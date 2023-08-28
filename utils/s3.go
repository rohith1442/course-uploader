package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func CreateSession(accessKey, secretKey, region string) (*session.Session, error) {
	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(region),
			Credentials: credentials.NewStaticCredentials(
				accessKey,
				secretKey,
				"",
			),
		})

	if err != nil {
		panic(err)
	}

	return sess, nil
}

func UploadObject(sess *session.Session, bucket, video, objectKey, filePath, mimeType string) error {
	uploader := s3manager.NewUploader(sess)
	fileContent, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer fileContent.Close()

	// Upload to s3
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(filepath.Join(video, objectKey)),
		Body:        fileContent,
		ContentType: aws.String(mimeType),
	})

	if err != nil {
		fmt.Println("failed to upload object:", err)
		return err
	}

	fmt.Printf("Successfully uploaded to %q\n", bucket)
	return nil
}
