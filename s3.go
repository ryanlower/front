package main

import (
	"bytes"
	"errors"
	"io"
	"log"

	"github.com/awslabs/aws-sdk-go/aws"
	awsS3 "github.com/awslabs/aws-sdk-go/service/s3"
)

type s3 struct {
	conf *config
}

func (s3 *s3) newClient() *awsS3.S3 {
	return awsS3.New(&aws.Config{
		Credentials: aws.Creds(s3.conf.AWS.AccessKeyID, s3.conf.AWS.SecretAccessKey, ""),
		Region:      s3.conf.AWS.Region,
	})
}

func (s3 *s3) read(path string) (io.ReadCloser, error) {
	log.Print("s3 read path=", path)

	resp, _ := s3.newClient().GetObject(&awsS3.GetObjectInput{
		Bucket: aws.String(s3.conf.S3.Bucket),
		Key:    aws.String(path),
	})
	// TODO, handle s3 errors

	if resp.LastModified == nil {
		return nil, errors.New("Object not found")
	}

	return resp.Body, nil
}

func (s3 *s3) write(path string, i *img) error {
	log.Print("s3 write path=", path)

	b := new(bytes.Buffer)
	i.encode(b)

	_, err := s3.newClient().PutObject(&awsS3.PutObjectInput{
		Bucket:       aws.String(s3.conf.S3.Bucket),
		Key:          aws.String(path),
		ContentType:  aws.String(i.contentType()),
		Body:         bytes.NewReader(b.Bytes()),
		StorageClass: aws.String("REDUCED_REDUNDANCY"), // TODO, allow customization
		ACL:          aws.String("public-read"),        // TODO, allow customization
	})
	// TODO, handle s3 errors

	return err
}
