![Satellite](https://github.com/user-attachments/assets/66196992-98c9-4376-992f-a0612a931cdc)

## Badges

[![Build Status](https://github.com/petaki/satellite/workflows/tests/badge.svg)](https://github.com/petaki/satellite/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-brightgreen.svg)](LICENSE.md)

## MCP Server

Satellite includes a built-in [MCP](https://modelcontextprotocol.io) (Model Context Protocol) server that lets AI agents query your monitoring data directly. When enabled, it exposes a Streamable HTTP endpoint with tools for listing probes, querying CPU/memory/load/disk metrics, reading logs, checking alerts, and managing probes.

Add it to your MCP client:

```bash
# Claude Code
claude mcp add --transport http satellite http://127.0.0.1:4000/mcp

# Codex
codex mcp add satellite --url http://127.0.0.1:4000/mcp
```

| Tool | Description |
|------|-------------|
| `list_probes` | List all probes with status, latest metrics, and alarm thresholds |
| `get_cpu` | CPU usage time series with top 3 processes |
| `get_memory` | Memory usage time series with top 3 processes |
| `get_load` | Load average time series (1m, 5m, 15m) |
| `get_disk` | Disk usage time series with path selector |
| `get_logs` | Log entries with log file path selector |
| `get_alerts` | Alarm thresholds and current metric values |
| `delete_probe` | Delete a probe and all its data |

Enable it by setting `MCP_ENABLED=true` in your `.env` file.

## Features

- **MCP Server** - AI agent access to all monitoring data via Model Context Protocol
- **CPU** - overall and per-process CPU usage charts
- **Memory** - overall and per-process memory usage charts
- **Load** - system load averages (1, 5, 15 min) charts
- **Disk** - per-partition disk usage charts with path selector
- **Log Viewer** - browse log tail entries collected by Probe
- **Dark Mode** - full dark mode support
- **Heartbeat** - webhook notifications when a probe goes offline
- **Series Filtering** - configurable time range buttons (5 min to 30 days)

## Getting Started

Follow the steps below to install and configure Satellite.

### Prerequisites

- Redis: `Version >= 5.0` for data reading

### Run with Docker

Image can be found at the package page on [GitHub](https://github.com/petaki/satellite/pkgs/container/satellite).

```bash
docker run --rm \
-e APP_URL=http://127.0.0.1:4000 \
-e REDIS_URL=redis://192.168.0.200:6379/0 \
-p 4000:4000 \
ghcr.io/petaki/satellite
```

### Install from Binary

Download the latest release for your platform from the [GitHub Releases](https://github.com/petaki/satellite/releases) page.

---

### Install from Source

#### Prerequisites

- Go: `Version >= 1.26`
- Node.js: `Version >= 22.0`
- Yarn or NPM

#### Steps

1. Clone the repository:

```bash
git clone git@github.com:petaki/satellite.git
```

2. Install UI dependencies and build:

```bash
cd satellite
yarn install
yarn build
```

3. Build the binary:

```bash
go build
```

4. Copy and edit the configuration:

```bash
cp .env.example .env
```

## Configuration

All configuration is done through environment variables in the `.env` file.

### General

#### Application Name

```
APP_NAME=
```

#### Application Address

```
APP_ADDR=:4000
```

#### Application URL

```
APP_URL=http://127.0.0.1:4000
```

#### Application Series Buttons

- Maximum `4` items.
- The first item is the `default`.
- The order does not matter from the second item.

```
APP_SERIES_BUTTONS=last_5_minutes,last_1_hour,last_24_hours,last_7_days
```

Available options:

| Option | Description |
|--------|-------------|
| `last_5_minutes` | Last 5 minutes |
| `last_15_minutes` | Last 15 minutes |
| `last_30_minutes` | Last 30 minutes |
| `last_1_hour` | Last 1 hour |
| `last_3_hours` | Last 3 hours |
| `last_6_hours` | Last 6 hours |
| `last_12_hours` | Last 12 hours |
| `last_24_hours` | Last 24 hours |
| `last_2_days` | Last 2 days |
| `last_7_days` | Last 7 days |
| `last_30_days` | Last 30 days |

---

### Redis

#### Redis URL

```
REDIS_URL=redis://127.0.0.1:6379/0
```

---

### MCP

Exposes monitoring data to AI agents via the Model Context Protocol. No authentication is included — control access at the network level.

#### MCP Enabled

```
MCP_ENABLED=false
```

---

### Heartbeat

Sends webhook notifications when a probe stops reporting. Requires Redis to be configured.

#### Heartbeat Enabled

```
HEARTBEAT_ENABLED=false
```

#### Heartbeat Wait (in minutes before first notification)

```
HEARTBEAT_WAIT=5
```

#### Heartbeat Sleep (in seconds between notifications)

Set to `0` to disable.

```
HEARTBEAT_SLEEP=300
```

---

### Heartbeat Webhook

#### Heartbeat Webhook Method

```
HEARTBEAT_WEBHOOK_METHOD=POST
```

#### Heartbeat Webhook URL

```
HEARTBEAT_WEBHOOK_URL=http://127.0.0.1:4000/heartbeat
```

#### Heartbeat Webhook Header

```
HEARTBEAT_WEBHOOK_HEADER='{"Authorization": "Bearer TOKEN", "Accept": "application/json"}'
```

#### Heartbeat Webhook Body

The body supports the following placeholders:

| Placeholder | Description |
|-------------|-------------|
| `%p` | Probe name |
| `%t` | Timestamp (`RFC3339`) |
| `%x` | Timestamp (`Unix`) |
| `%l` | Satellite link (relative) |

```
HEARTBEAT_WEBHOOK_BODY='{"probe": "%p", "timestamp_rfc3339": "%t", "timestamp_unix": %x, "link": "%l"}'
```

## Usage

```bash
./satellite web serve
```

## Data Collection

Collected data is provided by [Probe](https://github.com/petaki/probe).

## Contributors

- [@dyipon](https://github.com/dyipon) - development ideas, bug reports and testing

## Reporting Issues

If you are facing a problem with this package or found any bug, please open an issue on [GitHub](https://github.com/petaki/satellite/issues).

## License

The MIT License (MIT). Please see [License File](LICENSE.md) for more information.
