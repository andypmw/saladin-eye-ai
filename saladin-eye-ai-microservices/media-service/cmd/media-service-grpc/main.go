package main

import (
	"net"
	"os"

	"github.com/andypmw/saladin-eye-ai/media-service/common/genproto"
	grpcHandler "github.com/andypmw/saladin-eye-ai/media-service/handler/grpc"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func init() {
	// Log setup
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}

func main() {
	log.Info().Msg("SaladinEye.AI - Media Service - GRPC Server")

	// Start gRPC server
	server := grpc.NewServer()
	mediaService := grpcHandler.New()
	genproto.RegisterMediaServiceServer(server, mediaService)

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
