![Satellite](https://user-images.githubusercontent.com/1419087/133525370-79b6afe5-e54f-4eb2-b988-9b872322d89a.png)

## Badges

[![Build Status](https://github.com/petaki/satellite/workflows/tests/badge.svg)](https://github.com/petaki/satellite/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-brightgreen.svg)](LICENSE.md)

## Getting Started

Before you start, you need to install the prerequisites.

### Prerequisites

- Redis: `Version >= 5.0` for data reading

### Run with Docker

Image can be found at package page on [GitHub](https://github.com/petaki/satellite/pkgs/container/satellite).

```
docker run --rm \
-e APP_URL=http://127.0.0.1:4000 \
-e REDIS_URL=redis://192.168.0.200:6379/0 \
-p 4000:4000 \
ghcr.io/petaki/satellite
```

### Install from binary

Downloads can be found at releases page on [GitHub](https://github.com/petaki/satellite/releases).

---

### Install from source

#### Prerequisites for building

- GO: `Version >= 1.19`
- Node.js: `Version >= 14.0`
- Yarn or NPM

#### 1. Clone the repository:

```
git clone git@github.com:petaki/satellite.git
```

#### 2. Open the folder:

```
cd satellite
```

#### 3. Install the UI dependencies

```
yarn install
```

#### 4. Build the UI

```
yarn prod
```

#### 5. Build the Satellite:

```
go build
```

#### 6. Copy the example configuration:

```
cp .env.example .env
```

## Configuration

The configruation is stored in the `.env` file.

### Application Address

```
APP_ADDR=:4000
```

### Application URL

```
APP_URL=http://127.0.0.1:4000
```

---

### Redis URL

```
REDIS_URL=redis://127.0.0.1:6379/0
```

---

### Heartbeat Enabled

```
HEARTBEAT_ENABLED=false
```

### Heartbeat Wait (in minutes before first notification)

```
HEARTBEAT_WAIT=5
```

### Heartbeat Sleep (in seconds between notifications)

- `0` - Disabled

```
HEARTBEAT_SLEEP=300
```

### Heartbeat Webhook Method

```
HEARTBEAT_WEBHOOK_METHOD=POST
```

### Heartbeat Webhook URL

```
HEARTBEAT_WEBHOOK_URL=http://127.0.0.1:4000/heartbeat
```

### Heartbeat Webhook Header

```
HEARTBEAT_WEBHOOK_HEADER='{"Authorization": "Bearer TOKEN", "Accept": "application/json"}'
```

### Heartbeat Webhook Body

- `%p` - Probe
- `%t` - Start timestamp of current heartbeat period in `RFC3339` format
- `%x` - Start timestamp of current heartbeat period in `Unix` format
- `%l` - Satellite link (relative)

```
HEARTBEAT_WEBHOOK_BODY='{"probe": "%p", "timestamp_rfc3339": "%t", "timestamp_unix": %x, "link": "%l"}'
```

## Usage

Run the app using the following command:

```
./satellite web serve
```

## Data Collection

You can gather the necessary data with the [Probe](https://github.com/petaki/probe).

## Contributors

- [@dyipon](https://github.com/dyipon) for development ideas, bug reports and testing

## Reporting Issues

If you are facing a problem with this package or found any bug, please open an issue on [GitHub](https://github.com/petaki/satellite/issues).

## License

The MIT License (MIT). Please see [License File](LICENSE.md) for more information.
