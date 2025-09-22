# SMTP Duplicator (gosmtp-dup)

An SMTP relay server that duplicates incoming emails to multiple destinations - one primary server and multiple backup servers.

## Overview

This Go application acts as an SMTP proxy that receives emails and forwards them to:
- **Primary destination**: Critical delivery - if this fails, the original email sending fails
- **Backup destinations**: Additional copies sent asynchronously for redundancy

The main purpose of this program is to enable **smooth migration from one mail server to another** (or multiple servers) by duplicating traffic in real-time. This allows you to gradually transition your mail infrastructure while ensuring no emails are lost during the migration process.

## Features

- **Dual delivery mode**: Synchronous delivery to primary server, asynchronous delivery to backups
- **High availability**: Email redundancy across multiple servers
- **Configurable**: YAML-based configuration for easy deployment
- **Logging**: Comprehensive logging with structured output using Zap
- **SMTP compliance**: Full SMTP protocol support with UTF-8 encoding

## Configuration

The application uses a `config.yaml` file that can be placed in:
- Current directory (`.`)
- `/etc/smtp-dup/`

### Example Configuration

```yaml
smtp:
  listen: "127.0.0.1:2525"
  domain: "localhost"

relay:
  destination_primary: "mailprimary.example.com:25"
  destination_backups:
    - "mail_backup1.example.com:25"
    - "mail_backup2.example.com:25"
  timeout_seconds: 10
```

### Configuration Options

- `smtp.listen`: Address and port for the SMTP server to listen on
- `smtp.domain`: Domain name that the SMTP server announces itself as (used in SMTP greeting and protocol identification)
- `relay.destination_primary`: Primary mail server (synchronous delivery)
- `relay.destination_backups`: List of backup mail servers (asynchronous delivery)
- `relay.timeout_seconds`: Timeout for relay operations

### Environment Variables

You can configure the application using environment variables with the `SMTP_DUP_` prefix:

```bash
# SMTP settings
export SMTP_DUP_SMTP_LISTEN="0.0.0.0:2525"
export SMTP_DUP_SMTP_DOMAIN="mail.example.com"

# Relay settings
export SMTP_DUP_RELAY_DESTINATION_PRIMARY="mailprimary.example.com:25"
export SMTP_DUP_RELAY_DESTINATION_BACKUPS="mail1.backup.com:25,mail2.backup.com:25"
export SMTP_DUP_RELAY_TIMEOUT_SECONDS="30"
```

Environment variables take precedence over config file values. The config file becomes optional when using environment variables.

## How It Works

1. **Receives SMTP connections** on the configured listen address
2. **Forwards to primary server** synchronously - if this fails, the client receives an error
3. **Forwards to backup servers** asynchronously in the background - failures are logged but don't affect the client
4. **Logs all operations** for monitoring and debugging

## Use Cases

- **Smooth mail server migration**: Gradually migrate from old to new mail servers with real-time duplication
- **Email backup and redundancy**: Ensure emails are delivered to multiple servers
- **Mail archiving**: Send copies to archive servers while delivering to production
- **Load distribution**: Distribute email load across multiple servers
- **Disaster recovery**: Maintain backup mail servers for business continuity

## Integration with Postfix

To use this duplicator with Postfix, you can configure it as a transport in your `master.cf` file:

### Step 1: Add to master.cf

Add this line to your `/etc/postfix/master.cf`:

```
# SMTP duplicator transport
smtp-dup    unix  -       -       n       -       -       smtp
    -o smtp_generic_maps=
    -o smtp_destination_concurrency_limit=2
    -o smtp_destination_rate_delay=1s
    -o smtp_connect_timeout=30s
    -o smtp_helo_timeout=30s
```

### Step 2: Configure transport maps

In your `/etc/postfix/main.cf`, add:

```
transport_maps = hash:/etc/postfix/transport
```

### Step 3: Create transport file

Create `/etc/postfix/transport`:

```
# Route specific domains through the duplicator
example.com     smtp-dup:[127.0.0.1]:2525
.example.com    smtp-dup:[127.0.0.1]:2525
```

### Step 4: Update transport map

```bash
postmap /etc/postfix/transport
systemctl reload postfix
```

This configuration will route emails for `example.com` and its subdomains through your SMTP duplicator running on port 2525.

## Building and Running

### Native Build

```bash
# Build the application
go build -o gosmtp-dup

# Run with configuration file in current directory
./gosmtp-dup
```

### Docker

#### Using Pre-built Image from GitHub Registry

```bash
# Pull the latest image
docker pull ghcr.io/cabonemailserver/gosmtp-dup:latest

# Browse all available images at:
# https://github.com/CaboneMailServer/GoSMTP-dup/pkgs/container/gosmtp-dup

# Run with Docker (using environment variables)
docker run -d \
  --name smtp-duplicator \
  -p 2525:2525 \
  -e SMTP_DUP_SMTP_LISTEN="0.0.0.0:2525" \
  -e SMTP_DUP_SMTP_DOMAIN="mail.example.com" \
  -e SMTP_DUP_RELAY_DESTINATION_PRIMARY="mailprimary.example.com:25" \
  -e SMTP_DUP_RELAY_DESTINATION_BACKUPS="mail1.backup.com:25,mail2.backup.com:25" \
  ghcr.io/cabonemailserver/gosmtp-dup:latest

# Or run with config file
docker run -d \
  --name smtp-duplicator \
  -p 2525:2525 \
  -v $(pwd)/config-example.yaml:/app/config-example.yaml:ro \
  ghcr.io/cabonemailserver/gosmtp-dup:latest
```

#### Building Locally

```bash
# Build Docker image
docker build -t gosmtp-dup .

# Run with Docker (mount config file)
docker run -d \
  --name smtp-duplicator \
  -p 2525:2525 \
  -v $(pwd)/config-example.yaml:/app/config-example.yaml:ro \
  gosmtp-dup

# Run with Docker Compose
docker-compose up -d
```

#### Docker Compose Example

Create a `docker-compose.yml` file:

```yaml
version: '3.8'

services:
  smtp-duplicator:
    image: ghcr.io/cabonemailserver/gosmtp-dup:latest
    # Or build locally: build: .
    ports:
      - "2525:2525"
    volumes:
      - ./config-example.yaml:/app/config-example.yaml:ro
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "nc", "-z", "localhost", "2525"]
      interval: 30s
      timeout: 10s
      retries: 3
```

