---
layout: default
title: Installation
---

# Installation

## Binary Installation

### Linux

```bash
# Download the latest release
wget https://github.com/CaboneMailServer/GoSMTP-dup/releases/latest/download/gosmtp-dup-linux-amd64

# Make executable
chmod +x gosmtp-dup-linux-amd64

# Move to system path (optional)
sudo mv gosmtp-dup-linux-amd64 /usr/local/bin/gosmtp-dup
```

### Windows

1. Download `gosmtp-dup-windows-amd64.exe` from [releases](https://github.com/CaboneMailServer/GoSMTP-dup/releases/latest)
2. Place in your desired directory
3. Run from command prompt or PowerShell

### macOS

```bash
# Download for Intel Macs
wget https://github.com/CaboneMailServer/GoSMTP-dup/releases/latest/download/gosmtp-dup-darwin-amd64

# Download for Apple Silicon Macs
wget https://github.com/CaboneMailServer/GoSMTP-dup/releases/latest/download/gosmtp-dup-darwin-arm64

# Make executable
chmod +x gosmtp-dup-darwin-*

# Move to system path (optional)
sudo mv gosmtp-dup-darwin-* /usr/local/bin/gosmtp-dup
```

## Docker Installation

### Pull from Registry

```bash
# Pull latest image
docker pull ghcr.io/cabonemailserver/gosmtp-dup:latest

# Pull specific version
docker pull ghcr.io/cabonemailserver/gosmtp-dup:v1.0.0
```

### Build from Source

```bash
# Clone repository
git clone https://github.com/CaboneMailServer/GoSMTP-dup.git
cd GoSMTP-dup

# Build Docker image
docker build -t gosmtp-dup .
```

## Source Installation

### Prerequisites

- Go 1.24 or later

### Build from Source

```bash
# Clone repository
git clone https://github.com/CaboneMailServer/GoSMTP-dup.git
cd GoSMTP-dup

# Build
go build -o gosmtp-dup .

# Install (optional)
sudo cp gosmtp-dup /usr/local/bin/
```

## Systemd Service (Linux)

Create a systemd service file:

```bash
sudo tee /etc/systemd/system/gosmtp-dup.service > /dev/null <<EOF
[Unit]
Description=GoSMTP Duplicator
After=network.target

[Service]
Type=simple
User=nobody
Group=nogroup
ExecStart=/usr/local/bin/gosmtp-dup
Restart=always
RestartSec=5
WorkingDirectory=/etc/smtp-dup

[Install]
WantedBy=multi-user.target
EOF

# Reload systemd and start service
sudo systemctl daemon-reload
sudo systemctl enable gosmtp-dup
sudo systemctl start gosmtp-dup
```

## Verification

Test that the installation works:

```bash
# Check version (if available)
./gosmtp-dup --version

# Test with minimal config
echo "smtp:
  listen: \"127.0.0.1:2525\"
  domain: \"localhost\"
relay:
  destination_primary: \"localhost:25\"" > test-config.yaml

./gosmtp-dup
```

[← Back to Home](index.html) | [Next: Configuration →](configuration.html)