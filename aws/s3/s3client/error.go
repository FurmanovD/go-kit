package s3client

import "errors"

var (
	ErrPathNoBucketSeparator = errors.New("no bucket separator in S3 path")
	ErrPathNoKey             = errors.New("no key in S3 path")
)
