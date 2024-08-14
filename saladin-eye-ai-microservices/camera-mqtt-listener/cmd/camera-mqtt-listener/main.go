package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-redis/redis/v8"
)

var (
	ctx = context.Background()
	rdb *redis.Client
)

func init() {
	// Initialize Redis client
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		log.Fatal("Redis address is not set in the environment variables")
	}

	rdb = redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
}

func main() {
	log.Println("SaladinEye.AI - Camera MQTT Listener")
	log.Println("====================================")

	// Read MQTT credentials from environment variables
	mqttBroker := os.Getenv("MQTT_BROKER")
	mqttUsername := os.Getenv("MQTT_USERNAME")
	mqttPassword := os.Getenv("MQTT_PASSWORD")

	if mqttUsername == "" || mqttPassword == "" || mqttBroker == "" {
		log.Fatal("MQTT credentials or broker address are not set in the environment variables")
	}

	// MQTT broker connection options
	opts := mqtt.NewClientOptions().AddBroker(mqttBroker)
	opts.SetClientID("saladin-eye-camera-mqtt-listener")
	opts.SetUsername(mqttUsername) // Add your username here
	opts.SetPassword(mqttPassword) // Add your password here
	opts.SetDefaultPublishHandler(messageHandler)

	// Create and start an MQTT client
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	// Subscribe to the topic
	if token := client.Subscribe("saladin-eye/device/+/status", 0, nil); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	// Keep the program running
	select {}
}

func messageHandler(client mqtt.Client, msg mqtt.Message) {
	// Extract deviceId from the topic
	topicParts := strings.Split(msg.Topic(), "/")
	if len(topicParts) != 4 {
		return
	}
	deviceId := topicParts[2]

	log.Println("Set online presence for device:", deviceId)

	// Set the Redis key with TTL
	key := fmt.Sprintf("saladin-eye:camera-service:device-online-presence:%s", deviceId)
	err := rdb.Set(ctx, key, "1", 3*time.Minute).Err()
	if err != nil {
		log.Println("Failed to set Redis key:", err)
	}
}
