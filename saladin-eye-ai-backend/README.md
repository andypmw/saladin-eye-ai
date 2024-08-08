# SaladinEye.AI backend for online ESP32 Camera surveillance

This Nest.js application serves as the backend for the SaladinEye.AI online surveillance system, which utilizes ESP32 cameras. It manages device communication, data persistence, and real-time status tracking. The system employs PostgreSQL for long-term data storage and Redis for maintaining online status of ESP32 devices and caching. Communication with ESP32 devices is facilitated through an MQTT broker, with this backend acting as both publisher and subscriber to relevant topics. This architecture ensures efficient, scalable, and real-time operation of the SaladinEye.AI surveillance system.

## Installation

```bash
$ npm install
```

## Running the app

```bash
# development
$ npm run start

# watch mode
$ npm run start:dev

# production mode
$ npm run start:prod
```

## Test

```bash
# unit tests
$ npm run test

# e2e tests
$ npm run test:e2e

# test coverage
$ npm run test:cov
```

## Disclaimer

This project is intended for hobby and educational purposes only. It is not designed, developed, or tested to meet the standards and requirements of a production-grade system. There are no guarantees, warranties, or assurances provided regarding the performance, reliability, or security of this code. The author assumes no responsibility or liability for any damages, losses, or other issues that may arise from using this source code. Use at your own risk.
