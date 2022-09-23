![Satellite](https://user-images.githubusercontent.com/1419087/133525370-79b6afe5-e54f-4eb2-b988-9b872322d89a.png)

## Badges

[![Build Status](https://github.com/petaki/satellite/workflows/tests/badge.svg)](https://github.com/petaki/satellite/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-brightgreen.svg)](LICENSE.md)

## Getting Started

Before you start, you need to install the prerequisites.

### Prerequisites

- Redis: `Version >= 5.0` for data reading
- GO: `Version >= 1.19` for building
- Node.js: `Version >= 14.0` for building
- Yarn or NPM: for building

### Install with Docker

Image can be found at package page on [GitHub](https://github.com/petaki/satellite/pkgs/container/satellite).

### Install from binary

Downloads can be found at releases page on [GitHub](https://github.com/petaki/satellite/releases).

### Install from source

1. Clone the repository:

```
git clone git@github.com:petaki/satellite.git
```

2. Open the folder:

```
cd satellite
```

3. Install the UI dependencies

```
yarn install
```

4. Build the UI

```
yarn prod
```

5. Build the Satellite:

```
go build
```

6. Copy the example configuration:

```
cp .env.example .env
```

## Configuration

The configruation is stored in the `.env` file.

#### Application Address:

```
APP_ADDR=:4000
```

#### Application URL:

```
APP_URL=http://127.0.0.1:4000
```

#### Redis URL:

```
REDIS_URL=redis://127.0.0.1:6379/0
```

## Usage

Run the app using the following command:

```
./satellite web serve
```

## Data Collection

You can gather the necessary data with the [Probe](https://github.com/petaki/probe).

## Reporting Issues

If you are facing a problem with this package or found any bug, please open an issue on [GitHub](https://github.com/petaki/satellite/issues).

## License

The MIT License (MIT). Please see [License File](LICENSE.md) for more information.
