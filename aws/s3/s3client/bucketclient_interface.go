package s3client

import "time"

type BucketClient interface {
	Exists(key string) (bool, error)
	GetSize(key string, callerPays bool) (int64, error)
	GetObject(key string, callerPays bool) ([]byte, error)
	Delete(key string) error
	Copy(src, dst string, validateEtag, callerPays bool) error
	IsSrcNewer(src, dst string, callerPays bool) (bool, error)
	GetPresignedURL(key string, duration time.Duration) (string, error)
	GetETag(key string) (string, error)
}
