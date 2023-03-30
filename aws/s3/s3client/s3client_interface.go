package s3client

import "time"

// S3Client is a client for AWS S3 storage.
type S3Client interface {
	CreateBucket(bucket string) (string, error)
	Exists(obj S3Path) (bool, error)
	GetSize(obj S3Path, callerPays bool) (int64, error)
	GetObject(objPath S3Path, callerPays bool) ([]byte, error)
	Delete(obj S3Path) error
	Copy(src, dst S3Path, validateEtag, callerPays bool) error
	IsSrcNewer(src, dst S3Path, callerPays bool) (bool, error)
	GetPresignedURL(obj S3Path, duration time.Duration) (string, error)
	GetETag(obj S3Path) (string, error)
}
