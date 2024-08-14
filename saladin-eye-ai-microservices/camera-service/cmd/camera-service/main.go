package main

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/andypmw/saladin-eye-ai/camera-service/common/genproto"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ctx = context.Background()
	rdb *redis.Client
)

func init() {
	// Initialize Redis client
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		log.Fatal().Msg("redis address is not set in the environment variables")
	}

	rdb = redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
}

type CameraService struct {
	genproto.UnimplementedCameraServiceServer
}

func (CameraService) GetCameraStatus(_ context.Context, req *genproto.GetCameraStatusRequest) (*genproto.GetCameraStatusResponse, error) {
	deviceId := req.DeviceId

	log.Debug().Msgf("GetCameraStatus for device_id %s", deviceId)

	// Get the camera online presence status from Redis
	key := fmt.Sprintf("saladin-eye:camera-service:device-online-presence:%s", deviceId)
	cameraStatus, err := rdb.Get(ctx, key).Result()
	if err != nil {
		if err != redis.Nil {
			return nil, status.Errorf(codes.Internal, "internal error: %v", err)
		}
	}

	response := genproto.GetCameraStatusResponse{
		DeviceId: deviceId,
		IsOnline: cameraStatus == "1",
	}

	return &response, nil
}

func main() {
	// Log setup
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Start gRPC server
	server := grpc.NewServer()
	var cameraService CameraService
	genproto.RegisterCameraServiceServer(server, cameraService)

	log.Info().Msg("SaladinEye.AI - gRPC Server - Camera Service")

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
