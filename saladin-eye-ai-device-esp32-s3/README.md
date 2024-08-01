# SaladinEye.ai ESP32-S3 CCTV Device with AI/ML

SaladinEye.ai is an open-source project that aims to provide a robust and scalable CCTV system using the ESP32-S3 microcontroller. This project leverages AI and machine learning capabilities to capture and process images, focusing on efficient local storage and selective cloud uploading of important photos.

## Project Overview

- Platform: ESP32-S3 with 8MB PSRAM
- Language: C++
- Tooling: Platform-IO

Key Features:

1. Captures photos at configurable intervals.
1. Stores images on a local MicroSD card.
1. Performs local AI inference using TinyML or FOMO for computer vision tasks.
1. Uploads selected images to a centralized object storage based on AI filtering.
1. Communicate with SaladinEye.AI backend written on Node.js/Nest.js.

## Features

1. Efficient Image Capture: The system captures images at intervals defined by the user, optimizing storage and bandwidth usage.
1. Flexible Storage: Images can be stored on a MicroSD card, with support for hierarchical folder structures.
1. Local AI Processing: Utilizes TinyML or FOMO to process images locally, identifying and retaining only the most relevant ones.
1. Cloud Integration: The system can upload selected images to various cloud storage solutions, including Cloudflare R2, Amazon S3, Google Cloud Storage, and Digital Ocean Space.
1. Time Synchronization: Ensures accurate date and time information on captured images using NTP servers, with support for multiple time zones.
1. Scalable Management: SaladinEye.AI has a commmand center backend written on Node.js/Nest.js to seamlessly managing online and distributed CCTV based on ESP32-S3 microchip.

## Getting Started

### Prerequisites

- Visual Studio Code.
- Platform-IO installed on your development environment.
- An ESP32-S3 board with OV2640 camera module (recommendation: the one designed by Freenove).
- A MicroSD card for local storage.
- An RTC clock module.

### Installation

Clone the repository:

```
git clone https://github.com/andypmw/saladin-eye-ai.git
cd saladin-eye-ai
```

## Disclaimer

This project is intended for hobby and educational purposes only. It is not designed, developed, or tested to meet the standards and requirements of a production-grade system. There are no guarantees, warranties, or assurances provided regarding the performance, reliability, or security of this code. The author assumes no responsibility or liability for any damages, losses, or other issues that may arise from using this source code. Use at your own risk.
