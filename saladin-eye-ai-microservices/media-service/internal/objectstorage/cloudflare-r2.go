package objectstorage

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sort"
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

func (objs *CloudflareR2) GeneratePresignedUploadUrl(ctx context.Context, path string, durationMinute int) (string, error) {
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

func (objs *CloudflareR2) GeneratePresignedDownloadUrl(ctx context.Context, path string, durationMinute int) (string, error) {
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

func (objs *CloudflareR2) ListObjectsByPrefix(ctx context.Context, prefix string) ([]string, error) {
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

func (objs *CloudflareR2) ListDate(ctx context.Context, deviceId string) ([]string, error) {
	// Create input parameters
	input := &s3.ListObjectsV2Input{
		Bucket:    aws.String(objs.bucketName),
		Prefix:    aws.String(deviceId + "/"),
		Delimiter: aws.String("/"),
	}

	// Call ListObjectsV2 to get the list of objects
	result, err := objs.s3Client.ListObjectsV2(input)
	if err != nil {
		return nil, fmt.Errorf("failed to list objects: %v", err)
	}

	// Create a map to store unique dates
	dates := make(map[string]bool)

	// Process the CommonPrefixes to extract dates
	for _, commonPrefix := range result.CommonPrefixes {
		parts := strings.Split(*commonPrefix.Prefix, "/")
		if len(parts) >= 2 {
			date := parts[1]
			dates[date] = true
		}
	}

	// Build the list of dates
	dateList := make([]string, 0)
	for date := range dates {
		dateList = append(dateList, date)
	}

	// Sort the list of dates
	sort.Strings(dateList)

	return dateList, nil
}

func (objs *CloudflareR2) ListHourByDate(ctx context.Context, deviceId, date string) ([]string, error) {
	// Create input parameters
	input := &s3.ListObjectsV2Input{
		Bucket:    aws.String(objs.bucketName),
		Prefix:    aws.String(deviceId + "/" + date + "/"),
		Delimiter: aws.String("/"),
	}

	// Call ListObjectsV2 to get the list of objects
	result, err := objs.s3Client.ListObjectsV2(input)
	if err != nil {
		return nil, fmt.Errorf("failed to list objects: %v", err)
	}

	// Create a map to store unique hours
	hours := make(map[string]bool)

	// Process the CommonPrefixes to extract hours
	for _, commonPrefix := range result.CommonPrefixes {
		parts := strings.Split(*commonPrefix.Prefix, "/")
		if len(parts) >= 3 {
			hour := parts[2]
			hours[hour] = true
		}
	}

	// Build the list of hours
	hourList := make([]string, 0)
	for hour := range hours {
		hourList = append(hourList, hour)
	}

	// Sort the list of hours
	sort.Strings(hourList)

	return hourList, nil
}
