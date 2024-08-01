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
// - NTP: Network Time Protocol

#include <Arduino.h>
#include <WiFi.h>
#include <WiFiUdp.h>
#include <NTPClient.h>
#include "esp_log.h"

// Define NTP properties
// TODO - make the UTC offset, NTP server, and NTP port configurable via SaladinEye.AI Nest.js command center
const long utcOffsetInSeconds = 7 * 3600;
const char *ntpServer = "pool.ntp.org";
const int ntpPort = 123;

WiFiUDP ntpUDP;
// TODO - make the update interval configurable via SaladinEye.AI Nest.js command center
NTPClient timeClient(ntpUDP, ntpServer, utcOffsetInSeconds, 60000);

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
}

// Loop code, to run repeatedly
void loop()
{
  timeClient.update();

  // Print the formatted time
  log_d("Current time: %s", timeClient.getFormattedTime());

  // Print additional details if needed
  log_d("Epoch time: %d", timeClient.getEpochTime());

  // Wait for the next update
  delay(1000);
}
