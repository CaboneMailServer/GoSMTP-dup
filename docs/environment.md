---
layout: default
title: Environment Variables
---

# Environment Variables

GoSMTP-dup can be configured entirely using environment variables with the `SMTP_DUP_` prefix.

## Environment Variable Format

Environment variables use the prefix `SMTP_DUP_` followed by the configuration path with dots replaced by underscores:

- `smtp.listen` → `SMTP_DUP_SMTP_LISTEN`
- `relay.destination_primary` → `SMTP_DUP_RELAY_DESTINATION_PRIMARY`

## Available Variables

| Environment Variable | Config Equivalent | Example Value |
|---------------------|-------------------|---------------|
| `SMTP_DUP_SMTP_LISTEN` | `smtp.listen` | `"0.0.0.0:2525"` |
| `SMTP_DUP_SMTP_DOMAIN` | `smtp.domain` | `"mail.example.com"` |
| `SMTP_DUP_RELAY_DESTINATION_PRIMARY` | `relay.destination_primary` | `"primary.example.com:25"` |
| `SMTP_DUP_RELAY_DESTINATION_BACKUPS` | `relay.destination_backups` | `"backup1.com:25,backup2.com:25"` |
| `SMTP_DUP_RELAY_TIMEOUT_SECONDS` | `relay.timeout_seconds` | `"30"` |

## Usage Examples

### Shell Export

```bash
# Basic configuration
export SMTP_DUP_SMTP_LISTEN="0.0.0.0:2525"
export SMTP_DUP_SMTP_DOMAIN="mail.example.com"
export SMTP_DUP_RELAY_DESTINATION_PRIMARY="primary.example.com:25"
export SMTP_DUP_RELAY_DESTINATION_BACKUPS="backup1.example.com:25,backup2.example.com:25"
export SMTP_DUP_RELAY_TIMEOUT_SECONDS="30"

# Run the application
./gosmtp-dup
```

### Docker Environment

```bash
docker run -d \
  --name smtp-duplicator \
  -p 2525:2525 \
  -e SMTP_DUP_SMTP_LISTEN="0.0.0.0:2525" \
  -e SMTP_DUP_SMTP_DOMAIN="mail.example.com" \
  -e SMTP_DUP_RELAY_DESTINATION_PRIMARY="primary.example.com:25" \
  -e SMTP_DUP_RELAY_DESTINATION_BACKUPS="backup1.example.com:25,backup2.example.com:25" \
  -e SMTP_DUP_RELAY_TIMEOUT_SECONDS="30" \
  ghcr.io/cabonemailserver/gosmtp-dup:latest
```

### Docker Compose

```yaml
version: '3.8'
services:
  smtp-duplicator:
    image: ghcr.io/cabonemailserver/gosmtp-dup:latest
    ports:
      - "2525:2525"
    environment:
      SMTP_DUP_SMTP_LISTEN: "0.0.0.0:2525"
      SMTP_DUP_SMTP_DOMAIN: "mail.example.com"
      SMTP_DUP_RELAY_DESTINATION_PRIMARY: "primary.example.com:25"
      SMTP_DUP_RELAY_DESTINATION_BACKUPS: "backup1.example.com:25,backup2.example.com:25"
      SMTP_DUP_RELAY_TIMEOUT_SECONDS: "30"
    restart: unless-stopped
```

### Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: smtp-duplicator
spec:
  replicas: 1
  selector:
    matchLabels:
      app: smtp-duplicator
  template:
    metadata:
      labels:
        app: smtp-duplicator
    spec:
      containers:
      - name: smtp-duplicator
        image: ghcr.io/cabonemailserver/gosmtp-dup:latest
        ports:
        - containerPort: 2525
        env:
        - name: SMTP_DUP_SMTP_LISTEN
          value: "0.0.0.0:2525"
        - name: SMTP_DUP_SMTP_DOMAIN
          value: "mail.example.com"
        - name: SMTP_DUP_RELAY_DESTINATION_PRIMARY
          value: "primary.example.com:25"
        - name: SMTP_DUP_RELAY_DESTINATION_BACKUPS
          value: "backup1.example.com:25,backup2.example.com:25"
        - name: SMTP_DUP_RELAY_TIMEOUT_SECONDS
          value: "30"
```

## Configuration Priority

Environment variables take precedence over configuration files:

1. **Environment variables** (highest priority)
2. **Configuration file** (lower priority)
3. **Default values** (lowest priority)

## Default Values

| Setting | Default Value |
|---------|---------------|
| `smtp.listen` | `"127.0.0.1:2525"` |
| `smtp.domain` | `"localhost"` |
| `relay.timeout_seconds` | `10` |

## Array Values

For array configurations like `destination_backups`, use comma-separated values:

```bash
# Multiple backup servers
export SMTP_DUP_RELAY_DESTINATION_BACKUPS="server1.com:25,server2.com:25,server3.com:25"

# Single backup server
export SMTP_DUP_RELAY_DESTINATION_BACKUPS="backup.example.com:25"

# No backup servers (empty)
export SMTP_DUP_RELAY_DESTINATION_BACKUPS=""
```

## Validation

The application will:
- Use default values for missing optional settings
- Exit with error if required settings are missing
- Log configuration source (file vs environment variables)

[← Back: Configuration](configuration.html) | [Next: Docker Usage →](docker.html)