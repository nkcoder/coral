# Club Transfer Email Application

A Go application for processing club transfer data and sending notification emails to clubs.

## Configuration

The application can be configured using:

1. Command line flags
2. Environment variables
3. Configuration files

[Viper](https://github.com/spf13/viper) uses the following precedence order. Each item takes precedence over the item below it:
- explicit call to Set
- flag
- env
- config
- key/value store
- default

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

### Configuration File

A sample configuration file is provided in `config.sample.yaml`. Copy it to `config.yaml` and modify as needed:

```sh
cp config.sample.yaml config.yaml
```

Configuration is loaded from these locations in order:

1. `./config.yaml`
2. `./config/config.yaml`
3. `$HOME/.coral/config.yaml`
4. Custom path specified with `--config` flag
