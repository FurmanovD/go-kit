// Package s3client is a wrapper around the AWS S3 connection.
package s3client

import (
	"errors"
	"fmt"
	"io"
	"sort"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	// DefaultMultipartChunkSize is the max size of a multipart chunk.
	DefaultMultipartChunkSize int64 = 1024 * 1024 * 100 // 100Mb chunk for multipart upload by default.

	awsErrNotFound               = "NotFound"             // AWS wide constant.
	awsSinglePartCopyLimit int64 = 1024 * 1024 * 1024 * 4 // limit a single-chunk upload by to 4Gb.

	payerRequester = "requester"
)

type s3Client struct {
	awsS3 *s3.S3
}

// NewClientFromS3 sets the AWS S3 connection and returns an S3Client interface.
func NewClientFromS3(awsS3client *s3.S3) S3Client {
	return &s3Client{
		awsS3: awsS3client,
	}
}

// S3 returns an AWS S3 connection that was used when creating the object.
func (c *s3Client) S3() *s3.S3 {
	return c.awsS3
}

// CreateBucket ...
func (c *s3Client) CreateBucket(bucket string) (string, error) {
	bckt, err := c.awsS3.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		var awsErr awserr.RequestFailure
		if ok := errors.As(err, &awsErr); ok && awsErr.StatusCode() == 409 /*already exists*/ {
			return "/" + bucket, nil
		}

		return "", err
	}

	return *bckt.Location, nil
}

// Get returns an S3 object in a byte array view.
func (c *s3Client) GetObject(path S3Path, callerPays bool) ([]byte, error) {
	params := &s3.GetObjectInput{
		Bucket: aws.String(path.Bucket),
		Key:    aws.String(path.Key),
	}

	if callerPays {
		params.RequestPayer = aws.String(payerRequester)
	}

	resp, err := c.awsS3.GetObject(params)
	if err != nil {
		return nil, err
	}

	return io.ReadAll(resp.Body)
}

// GetSize returns a size in bytes of the object.
func (c *s3Client) GetSize(path S3Path, callerPays bool) (int64, error) {
	params := &s3.HeadObjectInput{
		Bucket: aws.String(path.Bucket),
		Key:    aws.String(path.Key),
	}

	if callerPays {
		params.RequestPayer = aws.String(payerRequester)
	}

	head, err := c.awsS3.HeadObject(params)
	if err != nil {
		return -1, err
	}

	return *head.ContentLength, nil
}

// Exists returns true if S3 object exists.
func (c *s3Client) Exists(path S3Path) (bool, error) {
	headParams := &s3.HeadObjectInput{
		Bucket: aws.String(path.Bucket),
		Key:    aws.String(path.Key),
	}

	_, err := c.awsS3.HeadObject(headParams)
	if err != nil {
		var awsErr awserr.Error
		if ok := errors.As(err, &awsErr); ok {
			if awsErr.Code() == awsErrNotFound {
				return false, nil
			}
		}

		return false, err
	}

	return true, nil
}

// Delete deletes an S3 object.
func (c *s3Client) Delete(path S3Path) error {
	params := &s3.DeleteObjectInput{
		Bucket: aws.String(path.Bucket),
		Key:    aws.String(path.Key),
	}

	_, err := c.awsS3.DeleteObject(params)
	if err != nil {
		var awsErr awserr.Error
		if ok := errors.As(err, &awsErr); ok {
			if awsErr.Code() == awsErrNotFound {
				return nil
			}
		}

		return err
	}

	return nil
}

// Copy copies source to destination and checks if required the result integrity
// by comparing an ETag of source and destination.
func (c *s3Client) Copy(src, dst S3Path, validateEtag, callerPays bool) error {
	return c.copyObject(src, dst, validateEtag, callerPays)
}

// copyObject copies source to destination.
func (c *s3Client) copyObject(src, dst S3Path, validateEtag bool, callerPays bool) error {
	srcSize, err := c.GetSize(src, callerPays)
	if err != nil {
		return err
	}

	eTag := ""
	if validateEtag {
		eTag, err = c.GetETag(src)
		if err != nil {
			return err
		}
	}

	if srcSize > awsSinglePartCopyLimit {
		if validateEtag {
			return fmt.Errorf("file size %v requires to use a miltipart copy operation that cannot be verified using ETags", srcSize)
		}

		return c.copyMultipartInt(src, dst, srcSize, DefaultMultipartChunkSize, callerPays)
	}

	copyParams := &s3.CopyObjectInput{
		Bucket:     aws.String(dst.Bucket),
		Key:        aws.String(dst.Key),
		CopySource: aws.String(src.Path()),
	}
	if callerPays {
		copyParams.RequestPayer = aws.String(payerRequester)
	}

	copyResult, err := c.awsS3.CopyObject(copyParams)
	if err != nil {
		return err
	}

	return validateCopyResult(validateEtag, eTag, copyResult)
}

func validateCopyResult(validateEtag bool, eTag string, copyResult *s3.CopyObjectOutput) error {
	if validateEtag {
		if copyResult == nil ||
			copyResult.CopyObjectResult == nil ||
			copyResult.CopyObjectResult.ETag == nil {
			return fmt.Errorf("copy result is %+v cannot be verified using ETags", copyResult)
		}

		if eTag != *copyResult.CopyObjectResult.ETag {
			return fmt.Errorf("copy operation result ETag %+v and the source one is %+v",
				*copyResult.CopyObjectResult.ETag,
				eTag,
			)
		}
	}

	return nil
}

// IsSrcNewer returns true if source exist and newer thad destination, or when destination does not exist.
func (c *s3Client) IsSrcNewer(src, dst S3Path, callerPays bool) (bool, error) {
	headParams := &s3.HeadObjectInput{
		Bucket: aws.String(src.Bucket),
		Key:    aws.String(src.Key),
	}
	if callerPays {
		headParams.RequestPayer = aws.String(payerRequester)
	}

	headSrc, err := c.awsS3.HeadObject(headParams)
	if err != nil {
		return false, fmt.Errorf("can not query source head %v : %w", src, err)
	}

	// check destination.
	headParams = &s3.HeadObjectInput{
		Bucket: aws.String(dst.Bucket),
		Key:    aws.String(dst.Key),
	}

	headDst, err := c.awsS3.HeadObject(headParams)
	if err != nil {
		var awsErr awserr.Error
		if ok := errors.As(err, &awsErr); ok && awsErr.Code() == awsErrNotFound {
			return true, nil
		}

		return false, fmt.Errorf("error querying head %v : %w", dst, err)
	}

	return headSrc.LastModified.After(*headDst.LastModified), nil
}

// GetPresignedURL returns an S3 presigned URL for the given key.
func (c *s3Client) GetPresignedURL(path S3Path, duration time.Duration) (string, error) {
	req, _ := c.awsS3.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(path.Bucket),
		Key:    aws.String(path.Key),
	})

	return req.Presign(duration)
}

// GetPresignedURL returns an S3 presigned URL for the given key.
func (c *s3Client) GetETag(path S3Path) (string, error) {
	headInput := &s3.HeadObjectInput{
		Bucket: aws.String(path.Bucket),
		Key:    aws.String(path.Key),
	}

	result, err := c.awsS3.HeadObject(headInput)
	if err != nil {
		return "", fmt.Errorf("error querying head for ETag %v : %w", path, err)
	}

	if result.ETag == nil {
		return "", fmt.Errorf("head returned a nil ETag %v", path)
	}

	return *result.ETag, nil
}

// completedParts a utility type used to sort completed parts in a multipart upload.
type completedParts []*s3.CompletedPart

func (a completedParts) Len() int           { return len(a) }
func (a completedParts) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a completedParts) Less(i, j int) bool { return *a[i].PartNumber < *a[j].PartNumber }

// nolint:funlen
func (c *s3Client) copyMultipartInt(
	src, dst S3Path,
	srcSize int64,
	chunkSize int64,
	callerPays bool) error {
	multipartParams := &s3.CreateMultipartUploadInput{
		Bucket: aws.String(dst.Bucket),
		Key:    aws.String(dst.Key),
	}
	upload, err := c.awsS3.CreateMultipartUpload(multipartParams)
	if err != nil {
		return err
	}

	partsNumber := srcSize/chunkSize + 1
	resCh := make(chan *s3.CompletedPart, partsNumber)
	defer close(resCh)
	errCh := make(chan error, partsNumber)
	defer close(errCh)
	parts := 0

	var bytePosition int64
	var wg sync.WaitGroup
	for i := int64(1); bytePosition < srcSize; i++ {
		endRange := bytePosition + chunkSize - 1
		if endRange >= srcSize {
			endRange = srcSize - 1
		}

		partRange := fmt.Sprintf("bytes=%v-%v", bytePosition, endRange)

		partCopyParam := &s3.UploadPartCopyInput{
			Bucket:          aws.String(dst.Bucket),
			Key:             aws.String(dst.Key),
			CopySource:      aws.String(src.Path()),
			UploadId:        upload.UploadId,
			CopySourceRange: aws.String(partRange),
			PartNumber:      aws.Int64(i),
		}
		if callerPays {
			partCopyParam.RequestPayer = aws.String(payerRequester)
		}

		wg.Add(1)
		go func(wg *sync.WaitGroup, param *s3.UploadPartCopyInput) {
			res, err := c.awsS3.UploadPartCopy(param)
			if err != nil {
				errCh <- fmt.Errorf("failed to copy part: %w", err)
			} else {
				resCh <- &s3.CompletedPart{
					ETag:       res.CopyPartResult.ETag,
					PartNumber: param.PartNumber,
				}
			}
			wg.Done()
		}(&wg, partCopyParam)

		bytePosition += chunkSize
		parts++
	}

	// wait until all parts are uploaded.
	wg.Wait()

	errStr := ""
	if len(errCh) > 0 {
		errStr = "multipart upload error(s):"
		for len(errCh) > 0 {
			errStr += fmt.Sprintf(" [%+v]", <-errCh)
		}
	}
	if errStr != "" {
		return errors.New(errStr)
	}

	partsArr := make(completedParts, parts)
	for i := 0; i < parts; i++ {
		partsArr[i] = <-resCh
	}
	sort.Sort(partsArr)

	completedUpload := &s3.CompletedMultipartUpload{
		Parts: partsArr,
	}

	completeParam := &s3.CompleteMultipartUploadInput{
		Bucket:          aws.String(dst.Bucket),
		Key:             aws.String(dst.Key),
		UploadId:        upload.UploadId,
		MultipartUpload: completedUpload,
	}
	_, err = c.awsS3.CompleteMultipartUpload(completeParam)

	return err
}
