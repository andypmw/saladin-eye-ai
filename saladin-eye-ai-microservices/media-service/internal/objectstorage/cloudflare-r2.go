package objectstorage

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type CloudflareR2 struct {
	ObjectStorageIface
	bucketName string
	s3Client   *s3.S3
}

func (objs *CloudflareR2) Init() error {
	// Get Cloudflare R2 credentials from environment variables
	accountID := os.Getenv("R2_ACCOUNT_ID")
	accessKeyID := os.Getenv("R2_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("R2_SECRET_ACCESS_KEY")
	objs.bucketName = os.Getenv("R2_BUCKET")

	// Check if all required environment variables are set
	if accountID == "" || accessKeyID == "" || secretAccessKey == "" || objs.bucketName == "" {
		return errors.New("missing Cloudflare R2 credentials in environment variables")
	}

	// Create a new AWS session
	awsSession, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
		Endpoint:    aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID)),
		Region:      aws.String("auto"),
	})
	if err != nil {
		return fmt.Errorf("failed to create cloudflare R2 session: %w", err)
	}

	// Create S3 service client
	objs.s3Client = s3.New(awsSession)

	return nil
}

func (objs *CloudflareR2) GeneratePresignedUploadUrl(path string, durationMinute int) (string, error) {
	s3req, _ := objs.s3Client.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(objs.bucketName),
		Key:    aws.String(path),
	})
	uploadUrlStr, err := s3req.Presign(time.Duration(durationMinute) * time.Minute)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return uploadUrlStr, nil
}

func (objs *CloudflareR2) GeneratePresignedDownloadUrl(path string, durationMinute int) (string, error) {
	s3req, _ := objs.s3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(objs.bucketName),
		Key:    aws.String(path),
	})
	downloadUrlStr, err := s3req.Presign(time.Duration(durationMinute) * time.Minute)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return downloadUrlStr, nil
}

func (objs *CloudflareR2) ListObjectsByPrefix(prefix string) ([]string, error) {
	filenames := make([]string, 0)

	resp, err := objs.s3Client.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(objs.bucketName),
		Prefix: aws.String(prefix),
	})
	if err != nil {
		return filenames, fmt.Errorf("failed to list objects: %w", err)
	}

	// Fill the filenames array
	for _, item := range resp.Contents {
		fileName := strings.TrimPrefix(*item.Key, prefix+"/")
		filenames = append(filenames, fileName)
	}

	return filenames, nil
}
