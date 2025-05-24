# Club Transfer Email Application

A Go application for processing club transfer data and sending notification emails to clubs.

## Configuration

You need to update the variables in `Taskfile.yml` to run for different environments and testing purpose:

```
vars:
  ENV: dev
  AWS_PROFILE: dev
  TEST_EMAIL: daniel.guo@vivalabs.com.au
  SENDER: no-reply@the-hub.ai
  PIF_INPUT: data/pif_club_transfer.csv
  DD_INPUT: data/dd_club_transfer.csv
```

## Running with Task locally

```sh
export AWS_PROFILE=xxx

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

### Run with Docker

```sh
task docker:send-email-pif
task docker:send-email-dd
```

### Run tests in Docker

```sh
task docker:test
```

