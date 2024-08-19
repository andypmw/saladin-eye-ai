package photo

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/andypmw/saladin-eye-ai/media-service/common/cache"
	"github.com/andypmw/saladin-eye-ai/media-service/common/constants"
	"github.com/andypmw/saladin-eye-ai/media-service/internal/objectstorage"
)

type PhotoServiceIface interface {
	GenerateUploadPresignedUrl(ctx context.Context, deviceId string) (string, error)
	ListObjectsByDateHour(ctx context.Context, deviceId string, date string, hour int32) ([]ObjectFile, error)
}

type PhotoServiceImpl struct {
	objStorage objectstorage.ObjectStorageIface
	rdb        redis.Cmdable
}

func New() (PhotoServiceIface, error) {
	objs, err := objectstorage.New(objectstorage.ProviderCloudflareR2)
	if err != nil {
		log.Fatal().Msgf("failed to create photo service: %v", err)
		return nil, fmt.Errorf("failed to create photo service: %w", err)
	}

	return &PhotoServiceImpl{
		objStorage: objs,
		rdb:        cache.New(),
	}, nil
}

/**
 * The media files in the object storage will be grouped like this:
 *   [Device ID]/[YYYY-MM-DD]/[HH]/[mm]-[ss].jpg
 *
 * The date time will be in UTC.
 */
func (ps *PhotoServiceImpl) GenerateUploadPresignedUrl(ctx context.Context, deviceId string) (string, error) {
	if len(deviceId) != 9 {
		msg := fmt.Sprintf("invalid device_id %s length %d", deviceId, len(deviceId))
		log.Error().Msg(msg)
		return "", errors.New(msg)
	}

	log.Debug().Msgf("GetPhotoUploadUrl for device_id %s", deviceId)

	// Generate the file name based on the current UTC time
	now := time.Now().UTC()
	fileName := fmt.Sprintf("%s/%s/%02d/%02d-%02d.jpg", deviceId, now.Format("2006-01-02"), now.Hour(), now.Minute(), now.Second())

	return ps.objStorage.GeneratePresignedUploadUrl(ctx, fileName, 15)
}

/**
 * The returned file names
 */
func (ps *PhotoServiceImpl) ListObjectsByDateHour(ctx context.Context, deviceId string, date string, hour int32) ([]ObjectFile, error) {
	// Initialize array to store filenames from cache or object storage API
	filenames := make([]string, 0)

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

	log.Debug().Msgf("ListObjectsByDateHour for device_id %s, date %s, hour %d", deviceId, date, hour)

	// Prefix in the object storage bucket
	prefix := fmt.Sprintf("%s/%s/%02d", deviceId, date, hour)

	// Create cache key
	redisKey := fmt.Sprintf("media-service:list-files-by-date-hour:%s:%s:%d", deviceId, date, hour)

	// Check cache
	cachedFiles, err := ps.rdb.LRange(ctx, redisKey, 0, -1).Result()
	if err != nil && err != redis.Nil {
		log.Error().Msgf("failed to get from cache: %v", err)
		return nil, fmt.Errorf("failed to get from cache: %w", err)
	}

	// Check cache hit first
	if len(cachedFiles) > 0 {
		log.Debug().Msg("cache hit")
		filenames = cachedFiles
	} else {
		// Cache miss, need to call the object-storage API.
		//
		// Then set the cache in Redis.
		//
		// Need to set TTL to the key:
		// - if the requested hour is the last hour, set TTL to 1 minute
		// - otherwise set TTL to 24 hours
		log.Debug().Msg("cache miss, call the object-storage API")

		// List objects in the S3 bucket
		resp, err := ps.objStorage.ListObjectsByPrefix(ctx, prefix)
		if err != nil {
			log.Error().Msgf("failed to list objects from object-storage API: %v", err)
			return nil, fmt.Errorf("failed to list objects from object-storage API: %w", err)
		}

		// Fill the filenames array
		for _, item := range resp {
			if strings.HasSuffix(item, ".jpg") || strings.HasSuffix(item, ".JPG") {
				filenames = append(filenames, item)
			}
		}

		// Set to cache
		if len(filenames) > 0 {
			_, err := ps.rdb.RPush(ctx, redisKey, filenames).Result()
			if err != nil {
				log.Error().Msgf("failed to set cache: %v", err)
				return nil, fmt.Errorf("failed to set cache: %w", err)
			}
		}

		// Need to set TTL to the key:
		// - if the requested hour is the last hour, set TTL to 1 minute
		// - otherwise set TTL to 24 hours

		// Step 1: Parse the date string "YYYY-MM-DD"
		parsedDate, err := time.Parse("2006-01-02", date)
		if err != nil {
			log.Error().Msgf("error parsing date from client: %v", err)
			return nil, fmt.Errorf("error parsing date from client: %w", err)
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
		_, err = ps.rdb.Expire(ctx, redisKey, cacheExpireMinute).Result()
		if err != nil {
			log.Error().Msgf("failed to set cache TTL in redis: %v", err)
			return nil, fmt.Errorf("failed to set cache TTL in redis: %w", err)
		}
	}

	// Build the result from filenames
	result := make([]ObjectFile, 0)
	for _, filename := range filenames {
		fullpath := fmt.Sprintf("%s/%s", prefix, filename)
		downloadURL, err := ps.objStorage.GeneratePresignedDownloadUrl(ctx, fullpath, constants.PHOTO_SERVICE_EXPIRATION_MINUTES)
		if err != nil {
			log.Error().Msgf("failed to generate presigned URL: %v", err)
			return nil, fmt.Errorf("failed to generate presigned URL: %w", err)
		}

		result = append(result, ObjectFile{
			Name:        filename,
			DownloadUrl: downloadURL,
		})
	}

	return result, nil
}
