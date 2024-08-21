package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	mqttHandler "github.com/andypmw/saladin-eye-ai/media-service/handler/mqtt"
)

func init() {
	// Log setup
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}

func main() {
	log.Info().Msg("SaladinEye.AI - Media Service - MQTT Handler")

	mqttHandler := mqttHandler.New()
	mqttHandler.Start()
}
