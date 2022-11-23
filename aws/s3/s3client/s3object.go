package s3client

import (
	"errors"
	"fmt"
	"strings"
)

const (
	s3Schema = "s3://"
)

// S3ObjectPath contains a ful path to a bucket object.
type S3ObjectPath struct {
	Bucket string
	Key    string
}

// NewS3Object creates a path object from a path: [bucket]/[key].
func NewS3Object(path string) (S3ObjectPath, error) {
	fullPath := path
	if strings.HasPrefix(strings.ToLower(path), s3Schema) {
		fullPath = path[5:]
	}

	bucketIdx := strings.Index(fullPath, "/")
	if bucketIdx == -1 {
		return S3ObjectPath{}, errors.New("no bucket separator in S3 path")
	}

	if bucketIdx >= len(fullPath)-1 {
		return S3ObjectPath{}, errors.New("bucket with no key in S3 path")
	}

	return S3ObjectPath{
		Bucket: fullPath[:bucketIdx],
		Key:    fullPath[bucketIdx+1:],
	}, nil
}

// Path returns an S3 path constructed.
func (p *S3ObjectPath) Path() string {
	return fmt.Sprintf("%v/%v", p.Bucket, p.Key)
}

// FullPath returns a full S3 path constructed including the schema.
func (p *S3ObjectPath) FullPath() string {
	return fmt.Sprintf("%s%v/%v", s3Schema, p.Bucket, p.Key)
}
