// Package s3client is a wrapper around the AWS S3 connection.
package s3client

import (
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
)

type bucketClient struct {
	client S3Client
	bucket string
}

// NewBucketClient returns a new S3Bucket
func NewBucketClient(s3 *s3.S3, bucket string) BucketClient {
	return &bucketClient{
		client: NewClientFromS3(s3),
		bucket: bucket,
	}
}

// NewBucketClientWithClient creates a new bucket's client using S3Client interface.
func NewBucketClientWithClient(client S3Client, bucket string) BucketClient {
	return &bucketClient{
		client: client,
		bucket: bucket,
	}
}

// Connection returns an AWS S3 connection that was used when creating the object.
func (c *bucketClient) Client() S3Client {
	return c.client
}

// CreateBucket ...
func (c *bucketClient) CreateBucket(bucket string) (string, error) {
	return c.client.CreateBucket(bucket)
}

// Get returns an S3 object in a byte array view.
func (c *bucketClient) GetObject(key string, callerPays bool) ([]byte, error) {
	return c.client.GetObject(
		S3Path{
			Bucket: c.bucket,
			Key:    key,
		},
		callerPays,
	)
}

// GetSize returns a size in bytes of the object.
func (c *bucketClient) GetSize(key string, callerPays bool) (int64, error) {
	return c.client.GetSize(
		S3Path{
			Bucket: c.bucket,
			Key:    key,
		},
		callerPays,
	)
}

// Exists returns true if S3 object exists.
func (c *bucketClient) Exists(key string) (bool, error) {
	return c.client.Exists(
		S3Path{
			Bucket: c.bucket,
			Key:    key,
		},
	)
}

// Delete deletes an S3 object.
func (c *bucketClient) Delete(key string) error {
	return c.client.Delete(
		S3Path{
			Bucket: c.bucket,
			Key:    key,
		},
	)
}

// Copy copies source to destination and checks if required the result integrity
// by comparing an ETag of source and destination.
func (c *bucketClient) Copy(src, dst string, validateEtag, callerPays bool) error {
	return c.client.Copy(
		S3Path{
			Bucket: c.bucket,
			Key:    src,
		},
		S3Path{
			Bucket: c.bucket,
			Key:    dst,
		},
		validateEtag,
		callerPays,
	)
}

// IsSrcNewer returns true if source exist and newer thad destination, or when destination does not exist.
func (c *bucketClient) IsSrcNewer(src, dst string, callerPays bool) (bool, error) {
	return c.client.IsSrcNewer(
		S3Path{
			Bucket: c.bucket,
			Key:    src,
		},
		S3Path{
			Bucket: c.bucket,
			Key:    dst,
		},
		callerPays,
	)
}

// GetPresignedURL returns an S3 presigned URL for the given key.
func (c *bucketClient) GetPresignedURL(key string, duration time.Duration) (string, error) {
	return c.client.GetPresignedURL(
		S3Path{
			Bucket: c.bucket,
			Key:    key,
		},
		duration,
	)
}

// GetPresignedURL returns an S3 presigned URL for the given key.
func (c *bucketClient) GetETag(key string) (string, error) {
	return c.client.GetETag(
		S3Path{
			Bucket: c.bucket,
			Key:    key,
		},
	)
}
