# Club Transfer Email Application

A Go application for processing club transfer data and sending notification emails to clubs.

## Running with Task

```sh
# Build and run PIF transfers
task send-email-pif

# Build and run DD transfers  
task send-email-dd

# Run tests
task test

# Run tests with coverage
task test:coverage
```

## Running with Docker

### Build the Docker image

```sh
task docker:build
# or
docker build -t club-transfer-app .
```

### Run with Docker

```sh
# Run with docker-compose
docker-compose up

# Run directly with docker
docker run --rm \
  -v $(pwd)/data:/app/data:ro \
  -v $(pwd)/config.yaml:/app/config.yaml:ro \
  -e AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID \
  -e AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY \
  -e AWS_SESSION_TOKEN=$AWS_SESSION_TOKEN \
  club-transfer-app \
  send-email -e dev -t PIF -i /app/data/pif_club_transfer.csv -s no-reply@the-hub.ai
```

### Run tests in Docker

```sh
task docker:test
```

## Configuration

The application can be configured using:

1. Command line flags
2. Environment variables
3. Configuration files

### Command Line Flags

```sh
# Basic usage
./email-app send-email --type=pif --input=data/pif_club_transfer.csv --env=dev --sender=no-reply@the-hub.ai

# Use a custom configuration file
./email-app --config=./my-config.yaml send-email --type=pif --input=data/pif_club_transfer.csv
```

### Environment Variables

Environment variables are prefixed with `CORAL_`:

```sh
# Set environment
export CORAL_ENVIRONMENT=dev

# Set sender email
export CORAL_EMAIL_SENDER=no-reply@the-hub.ai

# Set AWS region
export CORAL_AWS_REGION=ap-southeast-2
```

### Configuration Files

The application supports multiple configuration approaches:

1. **Single config file**: `config.yaml`
2. **Environment-specific files**: `config.dev.yaml`, `config.prod.yaml`
3. **Custom config file**: Use `--config` flag

### Loading Order

Configuration is loaded in the following order (later values override earlier ones):
1. Default values
2. Configuration files
3. Environment variables
4. Command line flags

### Example Configuration

```yaml
# config.yaml
environment: "dev"

email:
  sender: "no-reply@the-hub.ai"
  test: ""

aws:
  region: "ap-southeast-2"

worker:
  pool_size: 5
  delay_ms: 1000
```

## Testing

See [TESTING.md](TESTING.md) for comprehensive testing documentation.

```sh
# Run all tests
task test

# Run with coverage
task test:coverage

# Run with race detection
task test:race
```

## Development

```sh
# Install dependencies
task deps

# Format code
task fmt

# Run linter
task lint

# Clean build artifacts
task clean
```
