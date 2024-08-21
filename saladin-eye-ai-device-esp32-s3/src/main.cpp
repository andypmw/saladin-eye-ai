/**
 * Project: SaladinEye.AI
 * Author: Andy Primawan
 *
 * Description:
 *   An online and distributed CCTV with AI/ML capability,
 *   built using ESP32-S3 system-on-chip and OV2640 camera module.
 *
 *   This C++ code captures photos and stores them into a MicroSD card.
 *   It then performs computer vision processing on the captured photos.
 *
 *   The code communicates with the SaladinEye.AI backend,
 *   which is written using Node.js and Nest.js, for a variety of features.
 *
 */

// Glossaries
// - DS3231: a specific model of RTC module, refers to "Dallas Semiconductor" brand
// - I2C: Inter-Integrated Circuit, a communication protocol
// - NTP: Network Time Protocol
// - RTC: Real-Time Clock
// - SDA: Serial Data Line
// - SCL: Serial Clock Line

#include <iostream>
#include <fstream>
#include <Arduino.h>
#include <WiFi.h>
#include <WiFiUdp.h>
#include <NTPClient.h>
#include <Wire.h>
#include <RTClib.h>
#include <PubSubClient.h>
#include <HTTPClient.h>
#include <pb_encode.h>
#include <pb_decode.h>
#include "esp_log.h"
#include "esp_camera.h"
#include "FS.h"
#include "SD_MMC.h"
#include "genproto/media_service__get_photo_upload_url_request.pb.h"
#include "genproto/media_service__get_photo_upload_url_response.pb.h"

// GPIO pins for I2C communication with DS3231 RTC module
#define I2C_SDA 19
#define I2C_SCL 20

// GPIO pins for OV2640 camera module
#define PWDN_GPIO_NUM -1
#define RESET_GPIO_NUM -1
#define XCLK_GPIO_NUM 15
#define SIOD_GPIO_NUM 4
#define SIOC_GPIO_NUM 5
#define Y2_GPIO_NUM 11
#define Y3_GPIO_NUM 9
#define Y4_GPIO_NUM 8
#define Y5_GPIO_NUM 10
#define Y6_GPIO_NUM 12
#define Y7_GPIO_NUM 18
#define Y8_GPIO_NUM 17
#define Y9_GPIO_NUM 16
#define VSYNC_GPIO_NUM 6
#define HREF_GPIO_NUM 7
#define PCLK_GPIO_NUM 13

// GPIO pin for flash LED
#define FLASH_LED_GPIO_NUM 2

// GPIO pins for SD_MMC module
#define SD_MMC_CMD 38
#define SD_MMC_CLK 39
#define SD_MMC_D0 40

// Timer for capturing photos and publish request over MQTT
#define MSG_INTERVAL 10000

// Device ID, generate a unique ID for each device
const char* deviceId = "B7K9F2Q4L";

String mqttTopicMediaServiceString = "saladin-eye/server/media-service";
const char* mqttTopicMediaService = mqttTopicMediaServiceString.c_str();

String mqttTopicResponseWildcardString = "saladin-eye/device/" + String(deviceId) + "/response/#";
const char* mqttTopicResponseWildcard = mqttTopicResponseWildcardString.c_str();

String mqttTopicCommandString = "saladin-eye/device/command";
const char *mqttTopicCommand = mqttTopicCommandString.c_str();

// Define NTP properties
// TODO - make the UTC offset, NTP server, and NTP port configurable via SaladinEye.AI Nest.js command center
const long utcOffsetInSeconds = 7 * 3600;
const char *ntpServer = "pool.ntp.org";
const int ntpPort = 123;

WiFiUDP ntpUDP;
// TODO - make the update interval configurable via SaladinEye.AI Nest.js command center
NTPClient timeClient(ntpUDP, ntpServer, utcOffsetInSeconds, 60000);
// Define RTC
RTC_DS3231 rtc;

// MQTT
WiFiClient wifiClient;
const char* mqttBrokerAddress = "__REPLACE__";
String mqttClientIdString = "SaladinEye-ESP32S3-" + String(deviceId);
const char* mqttClientId = mqttClientIdString.c_str();
String mqttUsernameString = "SaladinEye-ESP32S3-" + String(deviceId);
const char* mqttUsername = mqttUsernameString.c_str();
const char* mqttPassword = "__REPLACE__";
PubSubClient mqttClient(wifiClient);
long lastMsg = 0;
char msg[50];
int value = 0;

// Functions declaration
void updateRtcFromNtp();
String getFormattedRtcTime();
bool initCamera();
bool initSdCard();
bool capturePhoto(String targetFullPath);
bool capturePhotoContinuously();
void mqttCallback(char *topic, byte *message, unsigned int length);
void mqttReconnect();
bool encode_string(pb_ostream_t* stream, const pb_field_t* field, void* const* arg);
bool decode_string(pb_istream_t *stream, const pb_field_t *field, void **arg);

// Setup code, to run once
void setup()
{
  Serial.begin(115200);

  // Give the serial port time to initialize
  delay(1000);

  // Connect to Wi-Fi
  // TODO - use proper way for SSID and Password secret while development, for example using env file
  const char *ssid = "__REPLACE__";
  const char *password = "__REPLACE__";
  WiFi.begin(ssid, password);
  log_i("Connecting to WiFi...");
  while (WiFi.status() != WL_CONNECTED)
  {
    delay(500);
    log_i(".");
  }
  log_i(" connected.");

  // Initialize NTP client
  timeClient.begin();

  // Initialize I2C bus
  Wire.begin(I2C_SDA, I2C_SCL);

  // Initialize RTC
  if (!rtc.begin())
  {
    while (1)
    {
      log_e("Couldn't find RTC");
      delay(1000);
    }
  }

  // Get the current time from NTP and set to RTC
  updateRtcFromNtp();

  // Initialize camera
  if (!initCamera())
  {
    while (1)
    {
      log_e("Couldn't initialize camera");
      delay(1000);
    }
  }

  // Initialize the flash LED pin as an output
  pinMode(FLASH_LED_GPIO_NUM, OUTPUT);

  // Initialize SdCard
  if (!initSdCard())
  {
    while (1)
    {
      log_e("Couldn't initialize SdCard");
      delay(1000);
    }
  }

  mqttClient.setBufferSize(2048);
  mqttClient.setServer(mqttBrokerAddress, 1883);
  mqttClient.setCallback(mqttCallback);
}

// Loop code, to run repeatedly
void loop()
{
  if (!mqttClient.connected()) {
    mqttReconnect();
  }
  mqttClient.loop();

  unsigned long now = millis();
  if (now - lastMsg > MSG_INTERVAL) {
    lastMsg = now;
    bool captureResult = capturePhotoContinuously();
    if (!captureResult)
    {
      log_e("Failed to capture photo and sending GetPhotoUploadUrlRequest over MQTT");
    }
  }
}

void updateRtcFromNtp()
{
  // Get the current timestamp from NTP
  timeClient.update();
  unsigned long epochTime = timeClient.getEpochTime();

  // Set RTC with NTP time
  rtc.adjust(DateTime(epochTime));
  log_i("RTC updated with NTP time");
}

String getFormattedRtcTime()
{
  DateTime now = rtc.now();

  // Buffer to hold formatted date/time string
  char buffer[20];

  // Format the date and time
  snprintf(buffer, sizeof(buffer), "%04d-%02d-%02d %02d:%02d:%02d",
           now.year(), now.month(), now.day(),
           now.hour(), now.minute(), now.second());

  return String(buffer);
}

bool initCamera()
{
  camera_config_t config;
  config.ledc_channel = LEDC_CHANNEL_0;
  config.ledc_timer = LEDC_TIMER_0;
  config.pin_d0 = Y2_GPIO_NUM;
  config.pin_d1 = Y3_GPIO_NUM;
  config.pin_d2 = Y4_GPIO_NUM;
  config.pin_d3 = Y5_GPIO_NUM;
  config.pin_d4 = Y6_GPIO_NUM;
  config.pin_d5 = Y7_GPIO_NUM;
  config.pin_d6 = Y8_GPIO_NUM;
  config.pin_d7 = Y9_GPIO_NUM;
  config.pin_xclk = XCLK_GPIO_NUM;
  config.pin_pclk = PCLK_GPIO_NUM;
  config.pin_vsync = VSYNC_GPIO_NUM;
  config.pin_href = HREF_GPIO_NUM;
  config.pin_sccb_sda = SIOD_GPIO_NUM;
  config.pin_sccb_scl = SIOC_GPIO_NUM;
  config.pin_pwdn = PWDN_GPIO_NUM;
  config.pin_reset = RESET_GPIO_NUM;
  config.xclk_freq_hz = 20000000;
  config.pixel_format = PIXFORMAT_JPEG;

  // We have 8MB PSRAM, so we can handle larger frame size and higher frame-buffer
  // Set the frame size to the highest possible
  config.frame_size = FRAMESIZE_UXGA; // UXGA (1600x1200) resolution
  config.jpeg_quality = 10;           // Higher quality setting (lower number means higher quality)
  config.fb_count = 2;                // Reduce frame buffer count to free up memory
  config.grab_mode = CAMERA_GRAB_LATEST;
  config.fb_location = CAMERA_FB_IN_PSRAM;

  // Init the camera
  esp_err_t err = esp_camera_init(&config);
  if (err != ESP_OK)
  {
    log_e("Camera init failed with error 0x%x", err);
    return false;
  }

  return true;
}

bool initSdCard()
{
  SD_MMC.setPins(SD_MMC_CLK, SD_MMC_CMD, SD_MMC_D0);

  if (!SD_MMC.begin("/sdcard", true, true, SDMMC_FREQ_DEFAULT, 5))
  {
    log_e("SdCard mount failed");
    return false;
  }

  uint8_t cardType = SD_MMC.cardType();
  if (cardType == CARD_NONE)
  {
    log_e("No SD card attached");
    return false;
  }

  if (cardType == CARD_MMC)
  {
    log_i("MMC card");
  }
  else if (cardType == CARD_SD)
  {
    log_i("SDSC card");
  }
  else if (cardType == CARD_SDHC)
  {
    log_i("SDHC card");
  }
  else
  {
    log_i("UNKNOWN card");
  }

  uint64_t cardSize = SD_MMC.cardSize() / (1024 * 1024);
  log_i("SD_MMC Card Size: %lluMB\n", cardSize);
  log_i("Total space: %lluMB\r\n", SD_MMC.totalBytes() / (1024 * 1024));
  log_i("Used space: %lluMB\r\n", SD_MMC.usedBytes() / (1024 * 1024));

  return true;
}

// Function to capture and save a photo
bool capturePhoto(String targetFullPath)
{
  // Turn on flash LED
  digitalWrite(FLASH_LED_GPIO_NUM, HIGH);

  // Open file handler
  File photoFile = SD_MMC.open(targetFullPath, FILE_WRITE);

  // Create frame-buffer variable
  camera_fb_t *fb = NULL;

  // Take a photo and save it to SD card
  fb = esp_camera_fb_get();
  if (!fb)
  {
    digitalWrite(FLASH_LED_GPIO_NUM, LOW);
    log_e("Camera capture failed");
    return false;
  }

  // Write the JPG data to file
  photoFile.write(fb->buf, fb->len);

  // Close file handler
  photoFile.close();

  // Release the frame buffer
  esp_camera_fb_return(fb);

  // Turn off flash LED
  digitalWrite(FLASH_LED_GPIO_NUM, LOW);

  log_i("Photo captured and saved: %s", targetFullPath.c_str());

  return true;
}

bool capturePhotoContinuously()
{
  // Get current time from RTC
  DateTime now = rtc.now();

  // Prepare the YYYY-MM-DD folder name from RTC
  char folderName[18];
  sprintf(folderName, "/%04d-%02d-%02d", now.year(), now.month(), now.day());

  // Prepare the HH foldername from RTC
  char hourFolderName[4];
  sprintf(hourFolderName, "/%02d", now.hour());

  // Prepare the mm-ss.jpg file name from RTC
  char photoFilename[12];
  sprintf(photoFilename, "%02d-%02d.jpg", now.minute(), now.second());

  // Make sure the YYYY-MM-DD folder exists in the SD card
  if (!SD_MMC.exists(folderName))
  {
    SD_MMC.mkdir(folderName);
  }

  // Make sure the YYYY-MM-DD/HH folder exists in the SD card
  if (!SD_MMC.exists(String(folderName) + String(hourFolderName)))
  {
    SD_MMC.mkdir(String(folderName) + String(hourFolderName));

    // Also need to create the list-photo.txt file in the folder
    File listPhotoFile = SD_MMC.open(String(folderName) + String(hourFolderName) + "/list-photo.txt", FILE_WRITE);

    // Close the file handler
    listPhotoFile.close();
  }

  // Capture the photo with full path YYYY-MM-DD/HH/mm-ss.jpg
  bool captureResult = capturePhoto(String(folderName) + String(hourFolderName) + "/" + String(photoFilename));

  if (!captureResult)
  {
    return false;
  }

  // Append the photo filename to the list-photo.txt file
  File listPhotoFile = SD_MMC.open(String(folderName) + String(hourFolderName) + "/list-photo.txt", FILE_APPEND);
  listPhotoFile.println(String(photoFilename));
  listPhotoFile.close();
  log_d("Photo filename appended to list-photo.txt");

  // The JPG path
  String photoPathString = String(folderName) + String(hourFolderName) + "/" + String(photoFilename);
  const char* photoPath = photoPathString.c_str();

  // Send message to media-service over MQTT
  saladineye_GetPhotoUploadUrlRequest request = saladineye_GetPhotoUploadUrlRequest_init_zero;

  request.device_id.arg = (void*)deviceId;
  request.device_id.funcs.encode = encode_string;

  request.original_photo_path.arg = (void*)photoPath;
  request.original_photo_path.funcs.encode = encode_string;

  // Encode the message
  uint8_t buffer[128];
  pb_ostream_t stream = pb_ostream_from_buffer(buffer, sizeof(buffer));
  bool status = pb_encode(&stream, saladineye_GetPhotoUploadUrlRequest_fields, &request);
  size_t message_length = stream.bytes_written;

  if (!status) {
    Serial.println("Failed to encode protobuf message");
    return false;
  }

  // Publish the message to MQTT

  // Create random string with length 10
  char randomString[11];
  for (int i = 0; i < 10; i++) {
    randomString[i] = (char)random(65, 90);
  }

  String publishTopicString = "saladin-eye/server/media-service/request/get-photo-upload-url/" + String(deviceId) + "/" + String(randomString);
  const char* publishTopic = publishTopicString.c_str();

  mqttClient.publish(publishTopic, buffer, message_length);
  Serial.println("FromDevice message sent to media-service over MQTT!");
  Serial.println(publishTopic);

  return true;
}

void mqttCallback(char *topic, byte *mqttMessage, unsigned int length)
{
  // Convert topic to String for easier manipulation
  String topicString = String(topic);

  // Find the position of "/response/"
  int startIndex = topicString.indexOf("/response/");
  if (startIndex != -1) {
      // Extract the substring after "/response/"
      String methodName = topicString.substring(startIndex + 10); // 10 is the length of "/response/"
      Serial.println("Extracted method name: " + methodName);

      if (methodName == "media-service/get-photo-upload-url") {
        // Decode the message
        saladineye_GetPhotoUploadUrlResponse response = saladineye_GetPhotoUploadUrlResponse_init_zero;
        pb_istream_t istream = pb_istream_from_buffer(mqttMessage, length);

        bool decode_status = pb_decode(&istream, saladineye_GetPhotoUploadUrlResponse_fields, &response);

        if (decode_status) {
          log_i("Protobuf decode success");
        } else {
            log_e("decoding protobuf saladineye_GetPhotoUploadUrlResponse failed");
            return;
        }

        // Open the JPG file
        File jpgFile = SD_MMC.open(response.original_photo_path);
        if (!jpgFile) {
          log_e("failed to open JPG file to be uploaded");
          return;
        }
        
        // Get the file size
        size_t fileSize = jpgFile.size();
        
        // Create an HTTP client
        HTTPClient http;
        
        // Begin the HTTP request
        http.begin(response.upload_url);
        http.addHeader("Content-Type", "image/jpeg");
        http.addHeader("Content-Length", String(fileSize));
        
        // Start the PUT request
        int httpResponseCode = http.sendRequest("PUT", &jpgFile, fileSize);
        
        if (httpResponseCode > 0) {
          String response = http.getString();
          log_i("OK HTTP Response: %s", response.c_str());
        } else {
          log_e("Error HTTP response code: %d", httpResponseCode);
        }
        
        // Close the file and HTTP connection
        jpgFile.close();
        http.end();
      } else {
        log_e("unknown MQTT method name");
      }
  } else {
      log_e("the /response/ not found in topic");
  }
}

void mqttReconnect()
{
  // Loop until we're reconnected
  while (!mqttClient.connected())
  {
    log_i("Attempting MQTT connection...");

    // Attempt to connect
    if (mqttClient.connect(mqttClientId, mqttUsername, mqttPassword))
    {
      log_i("MQTT connected");

      // Subscribe to the command topic
      if (mqttClient.subscribe(mqttTopicResponseWildcard, 1)) {
        log_i("Subscribed to topic successfully");
      } else {
        log_e("Subscription failed");
      }
    }
    else
    {
      log_i("failed, rc=");
      log_i("%d", mqttClient.state());
      log_i(" try again in 5 seconds");
      // Wait 5 seconds before retrying
      delay(5000);
    }
  }
}

bool encode_string(pb_ostream_t* stream, const pb_field_t* field, void* const* arg)
{
    const char* str = (const char*)(*arg);

    if (!pb_encode_tag_for_field(stream, field))
        return false;

    return pb_encode_string(stream, (uint8_t*)str, strlen(str));
}

bool decode_string(pb_istream_t *stream, const pb_field_t *field, void **arg) {
    Serial.println("decode_string function called");
    Serial.printf("Stream bytes left: %d\n", stream->bytes_left);

    static char str_buffer[1000];  // Adjust size as needed
    size_t length = stream->bytes_left;
    
    if (length >= sizeof(str_buffer)) {
        Serial.println("Decoding failed: Buffer too small");
        return false;
    }

    if (!pb_read(stream, (uint8_t *)str_buffer, length)) {
        Serial.println("Decoding failed: Reading from stream failed");
        return false;
    }
    
    str_buffer[length] = '\0';  // Null-terminate the string
    *(char**)arg = str_buffer;  // Set the pointer to our static buffer

    return true;
}
