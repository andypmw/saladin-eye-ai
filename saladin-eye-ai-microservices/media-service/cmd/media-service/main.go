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
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ctx        context.Context
	rdb        *redis.Client
	s3Client   *s3.S3
	bucketName string
)

func init() {
	// Setup Redis
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		log.Fatal().Msg("REDIS_ADDR environment variable not set")
	}

	// Initialize Redis client
	rdb = redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

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

	// Validations
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

	// Initialize an array to store filenames
	filenames := make([]string, 0)

	// Generate Redis key
	redisKey := fmt.Sprintf("media-service:list-files-by-date-hour:%s:%s:%d", deviceId, date, hour)

	// Check Redis cache
	cachedFiles, err := rdb.LRange(ctx, redisKey, 0, -1).Result()
	if err != nil && err != redis.Nil {
		log.Error().Msgf("failed to get cache from redis: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to get cache from redis: %v", err)
	}

	// Check cache hit first
	if len(cachedFiles) > 0 {
		log.Debug().Msg("cache hit")
		filenames = cachedFiles
	} else {
		// Cache miss, call the object-storage API.
		// Then set the cache in Redis.
		// Need to set TTL to the key:
		// - if the requested hour is the last hour, set TTL to 1 minute
		// - otherwise set TTL to 24 hours
		log.Debug().Msg("cache miss, call the object-storage API")

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

		// Fill the filenames array
		for _, item := range resp.Contents {
			fileName := strings.TrimPrefix(*item.Key, prefix+"/")
			if strings.HasSuffix(fileName, ".jpg") || strings.HasSuffix(fileName, ".JPG") {
				filenames = append(filenames, fileName)
			}
		}

		// Set the cache in Redis
		if len(filenames) > 0 {
			_, err := rdb.RPush(ctx, redisKey, filenames).Result()
			if err != nil {
				log.Error().Msgf("failed to set cache in redis: %v", err)
				return nil, status.Errorf(codes.Internal, "failed to set cache in redis: %v", err)
			}
		}

		// Need to set TTL to the key:
		// - if the requested hour is the last hour, set TTL to 1 minute
		// - otherwise set TTL to 24 hours

		// Step 1: Parse the date string "YYYY-MM-DD"
		parsedDate, err := time.Parse("2006-01-02", date)
		if err != nil {
			log.Error().Msgf("error parsing date: %v", err)
			return nil, status.Errorf(codes.Internal, "error parsing date: %v", err)
		}

		// Step 2: Add the hour to the parsed date
		constructedTime := parsedDate.Add(time.Duration(hour) * time.Hour)

		// Step 3: Get the current UTC time
		currentTime := time.Now().UTC()

		// Step 4: Calculate the difference
		timeDiff := currentTime.Sub(constructedTime)

		// Step 5: Check if the difference is more than 1 hour
		cacheExpireMinute := 1 * time.Minute
		if timeDiff > time.Hour {
			cacheExpireMinute = 24 * 60 * time.Minute
		}

		log.Debug().Msgf("set cache TTL to %v", cacheExpireMinute)

		// Set the TTL
		_, err = rdb.Expire(ctx, redisKey, cacheExpireMinute).Result()
		if err != nil {
			log.Error().Msgf("failed to set cache TTL in redis: %v", err)
			return nil, status.Errorf(codes.Internal, "failed to set cache TTL in redis: %v", err)
		}
	}

	// Build the files from filenames
	files := make([]*genproto.FileInfo, 0)
	for _, filename := range filenames {
		req, _ := s3Client.GetObjectRequest(&s3.GetObjectInput{
			Bucket: aws.String(bucketName),

			Key: aws.String(fmt.Sprintf("%s/%s/%02d/%s", deviceId, date, hour, filename)),
		})
		downloadURL, err := req.Presign(15 * time.Minute)
		if err != nil {
			log.Error().Msgf("failed to generate presigned URL: %v", err)
			return nil, status.Errorf(codes.Internal, "failed to generate presigned URL: %v", err)
		}
		files = append(files, &genproto.FileInfo{
			FileName:    filename,
			DownloadUrl: downloadURL,
		})
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
