package main

import (
	"bytes"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// AwsS3ClientAPI describes the AwsS3ClientAPI
type AwsS3ClientAPI struct {
	s3      *s3.S3
	session *session.Session
}

func NewAwsS3ClientAPI(keyID, secret string) (*AwsS3ClientAPI, error) {
	awsCfg := aws.NewConfig().WithRegion("eu-central-1") // it is not available in eu-central-1 yet
	// awsCfg.Credentials = credentials.NewCredentials(
	// 	&credentials.StaticProvider{
	// 		Value: credentials.Value{
	// 			AccessKeyID:     keyID,
	// 			SecretAccessKey: secret,
	// 		},
	// 	},
	// )
	session, err := session.NewSession(awsCfg)
	if err != nil {
		return nil, err
	}
	s3 := s3.New(session)
	return &AwsS3ClientAPI{
		s3:      s3,
		session: session,
	}, nil
}

func (api *AwsS3ClientAPI) GetObject(bucketName, key string) ([]byte, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	}
	resp, err := api.s3.GetObject(input)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func (api *AwsS3ClientAPI) PutObject(bucketName, key string, data []byte) error {
	input := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
	}
	_, err := api.s3.PutObject(input)
	return err
}

func (api *AwsS3ClientAPI) DeleteObject(bucketName, key string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	}
	_, err := api.s3.DeleteObject(input)
	return err
}

func (api *AwsS3ClientAPI) ListObjects(bucketName string) ([]*s3.Object, error) {
	input := &s3.ListObjectsInput{
		Bucket: aws.String(bucketName),
	}
	resp, err := api.s3.ListObjects(input)
	if err != nil {
		return nil, err
	}
	return resp.Contents, nil
}
