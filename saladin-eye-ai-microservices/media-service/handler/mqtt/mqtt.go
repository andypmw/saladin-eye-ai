package mqtt

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/andypmw/saladin-eye-ai/media-service/common/constants"
	"github.com/andypmw/saladin-eye-ai/media-service/common/genproto"
	"github.com/andypmw/saladin-eye-ai/media-service/service/photo"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

type MqttHandlerIface interface {
	Start()
	messageHandler(client mqtt.Client, msg mqtt.Message)
	handleGetPhotoUploadUrl(deviceId, idempotencyKey string, requestByteArr []byte) (string, []byte, error)
}

type MqttHandler struct {
	clientId      string
	brokerAddress string
	username      string
	password      string
	client        mqtt.Client
	photoService  photo.PhotoServiceIface
}

func New() MqttHandlerIface {
	clientId := os.Getenv("MQTT_CLIENT_ID")
	broker := os.Getenv("MQTT_BROKER")
	username := os.Getenv("MQTT_USERNAME")
	password := os.Getenv("MQTT_PASSWORD")

	if broker == "" || username == "" || password == "" {
		log.Fatal().Msg("MQTT_BROKER, MQTT_USERNAME, and MQTT_PASSWORD environment variables must be set")
	}

	photoService, err := photo.New()
	if err != nil {
		log.Fatal().Msgf("failed to create photo service: %v", err)
	}

	return &MqttHandler{
		clientId:      clientId,
		brokerAddress: broker,
		username:      username,
		password:      password,
		photoService:  photoService,
	}
}

func (handler *MqttHandler) Start() {
	opts := mqtt.NewClientOptions().AddBroker(handler.brokerAddress)
	opts.SetClientID(handler.clientId)
	opts.SetUsername(handler.username)
	opts.SetPassword(handler.password)
	opts.SetCleanSession(true)

	handler.client = mqtt.NewClient(opts)
	if token := handler.client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal().Msgf("failed to connect to MQTT broker: %v", token.Error())
	}

	if token := handler.client.Subscribe(constants.MQTT_TOPIC_SUBSCRIBE, 0, handler.messageHandler); token.Wait() && token.Error() != nil {
		log.Fatal().Msgf("failed to subscribe to MQTT topic: %v", token.Error())
	}

	// Block main thread so that the application continues running
	select {}
}

// Sample topic name:
//
//	saladin-eye/server/media-service/request/[method-name]/[device-id]/[idempotency-key]
//
// We need to check the last part, to know what kind of request it is.
// It will be like URL path for REST API.
func (handler *MqttHandler) messageHandler(client mqtt.Client, msg mqtt.Message) {
	log.Info().Msgf("received message on topic: %s", msg.Topic())

	topicParts := strings.Split(msg.Topic(), "/")

	if len(topicParts) != 7 {
		log.Error().Msgf("Invalid topic format: %s", msg.Topic())
		return
	}

	methodName := topicParts[len(topicParts)-3]
	deviceId := topicParts[len(topicParts)-2]
	idempotencyKey := topicParts[len(topicParts)-1]

	log.Info().Msgf("method name %s deviceId %s idempotencyKey %s", methodName, deviceId, idempotencyKey)

	switch methodName {
	case "get-photo-upload-url":
		responseTopic, responseByteArr, err := handler.handleGetPhotoUploadUrl(deviceId, idempotencyKey, msg.Payload())
		if err != nil {
			log.Error().Msgf("failed to handle GetPhotoUploadUrl: %v", err)
			return
		}

		qos := byte(1)
		if token := client.Publish(responseTopic, qos, false, responseByteArr); token.Wait() && token.Error() != nil {
			log.Error().Msgf("failed to publish MQTT message: %v", token.Error())
			return
		}

		log.Info().Msgf("published response to MQTT topic: %s", responseTopic)
	default:
		log.Error().Msgf("unknown method name: %s", methodName)
		return
	}
}

func (handler *MqttHandler) handleGetPhotoUploadUrl(deviceId, idempotencyKey string, requestByteArr []byte) (string, []byte, error) {
	request := &genproto.GetPhotoUploadUrlRequest{}
	if err := proto.Unmarshal(requestByteArr, request); err != nil {
		log.Error().Msgf("failed to unmarshal protobuf GetPhotoUploadUrlRequest: %v", err)
		return "", nil, fmt.Errorf("failed to unmarshal protobuf GetPhotoUploadUrlRequest: %w", err)
	}

	uploadURL, err := handler.photoService.GenerateUploadPresignedUrl(context.Background(), deviceId, idempotencyKey)
	if err != nil {
		log.Error().Msgf("failed to generate upload presigned URL: %v", err)
		return "", nil, fmt.Errorf("failed to generate upload presigned URL: %w", err)
	}

	response := &genproto.GetPhotoUploadUrlResponse{
		DeviceId:          deviceId,
		UploadUrl:         uploadURL,
		OriginalPhotoPath: request.OriginalPhotoPath,
	}

	responseByteArr, err := proto.Marshal(response)
	if err != nil {
		log.Error().Msgf("failed to marshal GetPhotoUploadUrlResponse: %v", err)
		return "", nil, fmt.Errorf("failed to marshal GetPhotoUploadUrlResponse: %w", err)
	}

	responseTopic := fmt.Sprintf("saladin-eye/device/%s/response/media-service/%s", deviceId, "get-photo-upload-url")

	return responseTopic, responseByteArr, nil
}
