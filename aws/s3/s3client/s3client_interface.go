package s3client

import "time"

type S3Client interface {
	CreateBucket(bucket string) (string, error)
	Exists(obj S3ObjectPath) (bool, error)
	GetSize(obj S3ObjectPath, callerPays bool) (int64, error)
	GetObject(objPath S3ObjectPath, callerPays bool) ([]byte, error)
	Delete(obj S3ObjectPath) error
	Copy(src, dst S3ObjectPath, validateEtag, callerPays bool) error
	IsSrcNewer(src, dst S3ObjectPath, callerPays bool) (bool, error)
	GetPresignedURL(obj S3ObjectPath, duration time.Duration) (string, error)
	GetETag(obj S3ObjectPath) (string, error)
}
