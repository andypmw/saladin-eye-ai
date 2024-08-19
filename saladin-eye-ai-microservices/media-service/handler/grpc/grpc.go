package grpc

import (
	"context"
	"strings"

	"github.com/andypmw/saladin-eye-ai/media-service/common/genproto"
	"github.com/andypmw/saladin-eye-ai/media-service/service/photo"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MediaService struct {
	genproto.UnimplementedMediaServiceServer
	photoService photo.PhotoServiceIface
}

func New() *MediaService {
	photoService, err := photo.New()
	if err != nil {
		log.Fatal().Msgf("failed to create photo service: %v", err)
	}

	return &MediaService{
		photoService: photoService,
	}
}

func (handler MediaService) GetPhotoUploadUrl(ctx context.Context, req *genproto.GetPhotoUploadUrlRequest) (*genproto.GetPhotoUploadUrlResponse, error) {
	deviceId := strings.TrimSpace(req.DeviceId)

	uploadURL, err := handler.photoService.GenerateUploadPresignedUrl(ctx, deviceId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate presigned photo upload URL: %v", err)
	}

	return &genproto.GetPhotoUploadUrlResponse{
		DeviceId:  deviceId,
		UploadUrl: uploadURL,
	}, nil
}

func (handler MediaService) ListFilesByDateHour(ctx context.Context, req *genproto.ListFilesByDateHourRequest) (*genproto.ListFilesByDateHourResponse, error) {
	deviceId := strings.TrimSpace(req.DeviceId)
	date := strings.TrimSpace(req.Date)
	hour := req.Hour

	result, err := handler.photoService.ListObjectsByDateHour(ctx, deviceId, date, hour)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list files by date hour: %v", err)
	}

	files := make([]*genproto.FileInfo, 0)
	for _, obj := range result {
		files = append(files, &genproto.FileInfo{
			FileName:    obj.Name,
			DownloadUrl: obj.DownloadUrl,
		})
	}

	return &genproto.ListFilesByDateHourResponse{
		TotalFiles: int32(len(result)),
		Files:      files,
	}, nil
}
