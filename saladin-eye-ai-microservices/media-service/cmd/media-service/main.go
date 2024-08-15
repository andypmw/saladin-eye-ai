package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/andypmw/saladin-eye-ai/media-service/common/genproto"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	s3Client   *s3.S3
	bucketName string
)

func init() {
	// Get R2 credentials from environment variables
	accountID := os.Getenv("R2_ACCOUNT_ID")
	accessKeyID := os.Getenv("R2_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("R2_SECRET_ACCESS_KEY")
	bucketName = os.Getenv("R2_BUCKET")

	// Check if all required environment variables are set
	if accountID == "" || accessKeyID == "" || secretAccessKey == "" || bucketName == "" {
		log.Fatal().Msg("missing Cloudflare R2 credentials in environment variables")
	}

	// Create a new AWS session
	awsSession, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
		Endpoint:    aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID)),
		Region:      aws.String("auto"),
	})
	if err != nil {
		log.Fatal().Msgf("failed to create cloudflare R2 session: %v", err)
	}

	// Create S3 service client
	s3Client = s3.New(awsSession)

	// Dummy call to list objects in the bucket with a limit of 1
	input := &s3.ListObjectsV2Input{
		Bucket:  aws.String(bucketName),
		MaxKeys: aws.Int64(1),
	}

	_, err = s3Client.ListObjectsV2(input)
	if err != nil {
		log.Fatal().Msgf("unable to list objects in bucket %s: %v", bucketName, err)
	}
}

type MediaService struct {
	genproto.UnimplementedMediaServiceServer
}

/**
 * The media files in the object storage will be grouped like this:
 *   [Device ID]/[YYYY-MM-DD]/[HH]/[mm]-[ss].jpg
 *
 * The date time will be in UTC.
 */
func (MediaService) GetPhotoUploadUrl(_ context.Context, req *genproto.GetPhotoUploadUrlRequest) (*genproto.GetPhotoUploadUrlResponse, error) {
	deviceId := strings.TrimSpace(req.DeviceId)

	if len(deviceId) != 9 {
		log.Error().Msgf("invalid device_id %s length %d", deviceId, len(deviceId))
		return nil, status.Errorf(codes.InvalidArgument, "invalid device_id length: %d", len(deviceId))
	}

	log.Debug().Msgf("GetPhotoUploadUrl for device_id %s", deviceId)

	// Generate the file name based on the current UTC time
	now := time.Now().UTC()
	fileName := fmt.Sprintf("%s/%s/%02d/%02d-%02d.jpg", deviceId, now.Format("2006-01-02"), now.Hour(), now.Minute(), now.Second())

	// Generate presigned URL
	s3req, _ := s3Client.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
	})
	uploadUrlStr, err := s3req.Presign(15 * time.Minute)
	if err != nil {
		log.Error().Msgf("failed to generate presigned URL: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to generate presigned URL: %v", err)
	}

	return &genproto.GetPhotoUploadUrlResponse{
		DeviceId:  deviceId,
		UploadUrl: uploadUrlStr,
	}, nil
}

func (MediaService) ListFilesByDateHour(ctx context.Context, req *genproto.ListFilesByDateHourRequest) (*genproto.ListFilesByDateHourResponse, error) {
	deviceId := strings.TrimSpace(req.DeviceId)
	date := strings.TrimSpace(req.Date)
	hour := req.Hour

	if len(deviceId) != 9 {
		log.Error().Msgf("invalid device_id %s length %d", deviceId, len(deviceId))
		return nil, status.Errorf(codes.InvalidArgument, "invalid device_id length: %d", len(deviceId))
	}

	if _, err := time.Parse("2006-01-02", date); err != nil {
		log.Error().Msgf("invalid date format %s", date)
		return nil, status.Errorf(codes.InvalidArgument, "invalid date format: %s", date)
	}

	if hour < 0 || hour > 23 {
		log.Error().Msgf("invalid hour %d", hour)
		return nil, status.Errorf(codes.InvalidArgument, "invalid hour: %d", hour)
	}

	log.Debug().Msgf("ListFilesByDateHour for device_id %s, date %s, hour %d", deviceId, date, hour)

	// List objects in the S3 bucket
	prefix := fmt.Sprintf("%s/%s/%02d", deviceId, date, hour)
	resp, err := s3Client.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
		Prefix: aws.String(prefix),
	})
	if err != nil {
		log.Error().Msgf("failed to list objects: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to list objects: %v", err)
	}

	var files []*genproto.FileInfo
	for _, item := range resp.Contents {
		fileName := strings.TrimPrefix(*item.Key, prefix+"/")
		if strings.HasSuffix(fileName, ".jpg") || strings.HasSuffix(fileName, ".JPG") {
			req, _ := s3Client.GetObjectRequest(&s3.GetObjectInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(*item.Key),
			})
			downloadURL, err := req.Presign(15 * time.Minute)
			if err != nil {
				log.Error().Msgf("failed to generate presigned URL: %v", err)
				return nil, status.Errorf(codes.Internal, "failed to generate presigned URL: %v", err)
			}
			files = append(files, &genproto.FileInfo{
				FileName:    fileName,
				DownloadUrl: downloadURL,
			})
		}
	}

	return &genproto.ListFilesByDateHourResponse{
		TotalFiles: int32(len(files)),
		Files:      files,
	}, nil
}

func main() {
	// Log setup
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Start gRPC server
	server := grpc.NewServer()
	var mediaService MediaService
	genproto.RegisterMediaServiceServer(server, mediaService)

	log.Info().Msg("SaladinEye.AI - gRPC Server - Media Service")

	port := os.Getenv("GRPC_PORT")
	if port == "" {
		log.Fatal().Msg("GRPC_PORT environment variable not set")
	}

	l, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal().Msgf("could not listen to %s: %v", port, err)
	}

	log.Fatal().Err(server.Serve(l))
}
