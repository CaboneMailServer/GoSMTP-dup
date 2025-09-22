---
layout: default
title: Docker Usage
---

# Docker Usage

GoSMTP-dup provides multi-architecture Docker images for easy deployment.

## Quick Start

```bash
# Pull and run with environment variables
docker run -d \
  --name smtp-duplicator \
  -p 2525:2525 \
  -e SMTP_DUP_SMTP_LISTEN="0.0.0.0:2525" \
  -e SMTP_DUP_RELAY_DESTINATION_PRIMARY="primary.example.com:25" \
  ghcr.io/cabonemailserver/gosmtp-dup:latest
```

## Image Information

- **Registry**: [ghcr.io/cabonemailserver/gosmtp-dup](https://github.com/CaboneMailServer/GoSMTP-dup/pkgs/container/gosmtp-dup)
- **Architecture**: linux/amd64
- **Base**: Alpine Linux (minimal size)
- **User**: Non-root for security

## Available Tags

| Tag | Description |
|-----|-------------|
| `latest` | Latest stable release |
| `v1.x.x` | Specific version (e.g., `v1.0.0`) |

## Configuration Methods

### Method 1: Environment Variables (Recommended)

```bash
docker run -d \
  --name smtp-duplicator \
  -p 2525:2525 \
  -e SMTP_DUP_SMTP_LISTEN="0.0.0.0:2525" \
  -e SMTP_DUP_SMTP_DOMAIN="mail.example.com" \
  -e SMTP_DUP_RELAY_DESTINATION_PRIMARY="primary.example.com:25" \
  -e SMTP_DUP_RELAY_DESTINATION_BACKUPS="backup1.com:25,backup2.com:25" \
  ghcr.io/cabonemailserver/gosmtp-dup:latest
```

### Method 2: Configuration File

```bash
# Create config file
cat > config.yaml << EOF
smtp:
  listen: "0.0.0.0:2525"
  domain: "mail.example.com"
relay:
  destination_primary: "primary.example.com:25"
  destination_backups:
    - "backup1.example.com:25"
    - "backup2.example.com:25"
EOF

# Run with mounted config
docker run -d \
  --name smtp-duplicator \
  -p 2525:2525 \
  -v $(pwd)/config.yaml:/app/config.yaml:ro \
  ghcr.io/cabonemailserver/gosmtp-dup:latest
```

## Docker Compose

### Basic Setup

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
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "nc", "-z", "localhost", "2525"]
      interval: 30s
      timeout: 10s
      retries: 3
```

### With Configuration File

```yaml
version: '3.8'

services:
  smtp-duplicator:
    image: ghcr.io/cabonemailserver/gosmtp-dup:latest
    ports:
      - "2525:2525"
    volumes:
      - ./config.yaml:/app/config.yaml:ro
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "nc", "-z", "localhost", "2525"]
      interval: 30s
      timeout: 10s
      retries: 3
```

## Health Checks

The Docker image includes `netcat` for health checking:

```bash
# Manual health check
docker exec smtp-duplicator nc -z localhost 2525

# Docker health check (built-in)
docker run -d \
  --name smtp-duplicator \
  --health-cmd="nc -z localhost 2525" \
  --health-interval=30s \
  --health-timeout=10s \
  --health-retries=3 \
  -p 2525:2525 \
  -e SMTP_DUP_RELAY_DESTINATION_PRIMARY="primary.example.com:25" \
  ghcr.io/cabonemailserver/gosmtp-dup:latest
```

## Kubernetes Deployment

### Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: smtp-duplicator
  labels:
    app: smtp-duplicator
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
        livenessProbe:
          tcpSocket:
            port: 2525
          initialDelaySeconds: 10
          periodSeconds: 10
        readinessProbe:
          tcpSocket:
            port: 2525
          initialDelaySeconds: 5
          periodSeconds: 5
```

### Service

```yaml
apiVersion: v1
kind: Service
metadata:
  name: smtp-duplicator-service
spec:
  selector:
    app: smtp-duplicator
  ports:
  - protocol: TCP
    port: 2525
    targetPort: 2525
  type: ClusterIP
```

## Logging

View container logs:

```bash
# Follow logs
docker logs -f smtp-duplicator

# Last 100 lines
docker logs --tail 100 smtp-duplicator
```

## Troubleshooting

### Common Issues

1. **Port already in use**:
   ```bash
   # Check what's using port 2525
   netstat -tulpn | grep 2525

   # Use different port
   docker run -p 2526:2525 ...
   ```

2. **Permission denied**:
   ```bash
   # The container runs as non-root, ensure config files are readable
   chmod 644 config.yaml
   ```

3. **Configuration not loaded**:
   ```bash
   # Check if environment variables are set
   docker exec smtp-duplicator env | grep SMTP_DUP

   # Check mounted files
   docker exec smtp-duplicator ls -la /app/
   ```

### Debug Mode

```bash
# Run with interactive shell for debugging
docker run -it --rm \
  -p 2525:2525 \
  ghcr.io/cabonemailserver/gosmtp-dup:latest \
  /bin/sh
```

[← Back: Environment Variables](environment.html) | [Next: Postfix Integration →](postfix.html)