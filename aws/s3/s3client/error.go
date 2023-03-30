package s3client

import "errors"

var (
	ErrPathNoBucketSeparator = errors.New("no bucket separator in S3 path")
	ErrPathEmptyPath         = errors.New("no key in S3 path")
)
