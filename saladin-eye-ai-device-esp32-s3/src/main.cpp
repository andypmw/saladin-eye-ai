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

#include <Arduino.h>
#include <WiFi.h>
#include <WiFiUdp.h>
#include <NTPClient.h>
#include <Wire.h>
#include <RTClib.h>
#include "esp_log.h"
#include "esp_camera.h"
#include "FS.h"
#include "SD_MMC.h"

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

// Functions declaration
void updateRtcFromNtp();
String getFormattedRtcTime();
bool initCamera();
bool initSdCard();
bool capturePhoto(String targetFullPath);
bool capturePhotoContinuously();

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
  Serial.print("Connecting to WiFi...");
  while (WiFi.status() != WL_CONNECTED)
  {
    delay(500);
    Serial.print(".");
  }
  Serial.println(" connected.");

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
}

// Loop code, to run repeatedly
void loop()
{
  // Print the formatted time
  log_d("Current RTC time: %s", getFormattedRtcTime().c_str());

  // Capture photo continously, and check the result value
  bool captureResult = capturePhotoContinuously();
  if (!captureResult)
  {
    log_e("Failed to capture photo");
  }

  // Wait for the next update
  delay(1000);
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

  return true;
}