package filestore

import (
	"bytes"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// var endpoint = "http://localhost:4572"
var endpoint = os.Getenv("S3_ENDPOINT")

const (
	disableSSL      = true
	accessKeyID     = "x"
	secretAccessKey = "x"
	secretToken     = "x"
)

// FileManager Defines the structure used to manage the Files on Amazon S3
type FileManager struct {
	*s3.S3
}

// New Initializes the FileManager
func New() *FileManager {
	return &FileManager{}
}

// CreateBucket to store files
func (f *FileManager) CreateBucket(name string) error {
	sess, err := session.NewSession(&aws.Config{
		Endpoint:    aws.String(endpoint),
		DisableSSL:  aws.Bool(disableSSL),
		Credentials: credentials.NewStaticCredentials(accessKeyID, secretAccessKey, secretToken),
		Region:      aws.String("us-west-2")},
	)
	if err != nil {
		return nil
	}

	svc := s3.New(sess)
	_, err = svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(name),
	})
	if err != nil {
		return err
	}

	// Wait until bucket is created before finishing
	err = svc.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(name),
	})
	return err
}

func (f *FileManager) ListItems(bucket, prefix string) (filenames []string, err error) {
	sess, err := session.NewSession(&aws.Config{
		Endpoint:    aws.String(endpoint),
		DisableSSL:  aws.Bool(disableSSL),
		Credentials: credentials.NewStaticCredentials(accessKeyID, secretAccessKey, secretToken),
		Region:      aws.String("us-west-2")},
	)
	if err != nil {
		return nil, nil
	}

	svc := s3.New(sess)
	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(bucket),
		Prefix: aws.String(prefix)})
	if err != nil {
		return nil, err
	}

	filenames = make([]string, len(resp.Contents))
	for index, item := range resp.Contents {
		filenames[index] = *item.Key
	}
	return filenames, nil
}

// Upload file to bucket
func (f *FileManager) Upload(bucket, filename string, reader io.Reader) error {
	sess, err := session.NewSession(&aws.Config{
		Endpoint:    aws.String(endpoint),
		DisableSSL:  aws.Bool(disableSSL),
		Credentials: credentials.NewStaticCredentials(accessKeyID, secretAccessKey, secretToken),
		Region:      aws.String("us-west-2")},
	)
	if err != nil {
		return err
	}
	uploader := s3manager.NewUploader(sess)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
		Body:   reader,
	})
	return err
}

// Download File from bucket
func (f *FileManager) Download(bucket, filename string) (io.Reader, error) {
	sess, err := session.NewSession(&aws.Config{
		Endpoint:    aws.String(endpoint),
		DisableSSL:  aws.Bool(disableSSL),
		Credentials: credentials.NewStaticCredentials(accessKeyID, secretAccessKey, secretToken),
		Region:      aws.String("us-west-2")},
	)
	if err != nil {
		return nil, err
	}
	downloader := s3manager.NewDownloader(sess)
	buffer := aws.NewWriteAtBuffer([]byte{})
	_, err = downloader.Download(buffer,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(filename),
		})

	return bytes.NewReader(buffer.Bytes()), err
}
