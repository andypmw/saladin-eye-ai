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

// GPIO pins for I2C communication with DS3231 RTC module
#define I2C_SDA 19
#define I2C_SCL 20

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
}

// Loop code, to run repeatedly
void loop()
{
  // Print the formatted time
  log_d("Current RTC time: %s", getFormattedRtcTime().c_str());

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