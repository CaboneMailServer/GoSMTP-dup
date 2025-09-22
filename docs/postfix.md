---
layout: default
title: Postfix Integration
---

# Postfix Integration

Learn how to integrate GoSMTP-dup with Postfix for seamless email duplication.

## Overview

Postfix can route specific domains or all emails through the SMTP duplicator using transport maps. This allows you to duplicate emails without changing your existing mail flow.

## Configuration Steps

### Step 1: Configure the SMTP Duplicator

First, configure GoSMTP-dup to listen and relay emails:

```yaml
# config.yaml
smtp:
  listen: "127.0.0.1:2525"
  domain: "mail.example.com"
relay:
  destination_primary: "oldmail.example.com:25"    # Current mail server
  destination_backups:
    - "newmail.example.com:25"                      # New mail server (migration target)
    - "archive.example.com:25"                      # Archive server
```

### Step 2: Add Transport to master.cf

You have two options for running the duplicator:

#### Option A: Managed by Postfix Master (Recommended)

Let Postfix manage the duplicator as a permanent service:

```
# SMTP duplicator permanent service
127.0.0.1:2525  inet  n       -       y       -       1       /usr/local/bin/gosmtp-dup
    user=postfix directory=/etc/smtp-dup
```

**Note**: Runs permanently as an inet service managed by Postfix. The proxy speaks SMTP natively, so Postfix can communicate with it directly.

#### Option B: External Service

For external service, no changes needed in `/etc/postfix/master.cf`. Postfix will use the default SMTP transport to connect to the external service.

### Step 3: Configure Transport Maps

In your `/etc/postfix/main.cf`, add or update:

```
transport_maps = hash:/etc/postfix/transport
```

### Step 4: Create Transport File

Create `/etc/postfix/transport` to define which emails go through the duplicator:

**For both options (same configuration):**

```
# Route specific domains through the duplicator
example.com     [127.0.0.1]:2525
.example.com    [127.0.0.1]:2525
test.com        [127.0.0.1]:2525

# Route all outbound email through duplicator
*               [127.0.0.1]:2525

# Route specific recipients
user@example.com        [127.0.0.1]:2525
admin@example.com       [127.0.0.1]:2525
```

### Step 5: Update Transport Map and Reload

```bash
# Compile the transport map
postmap /etc/postfix/transport

# For Option B: Start external service
systemctl start gosmtp-dup
systemctl enable gosmtp-dup

# Reload Postfix configuration
systemctl reload postfix

# Verify configuration
postfix check
```

**Note**: For Option A (Postfix-managed), the duplicator process will be started automatically by Postfix when needed.

## Testing the Integration

### Test Email Flow

```bash
# Send a test email
echo "Test email" | mail -s "Test Subject" user@example.com

# Check Postfix logs
tail -f /var/log/mail.log

# Check duplicator logs
journalctl -u gosmtp-dup -f
```

### Verify Transport

```bash
# Check if transport is applied
postmap -q example.com /etc/postfix/transport

# Test transport lookup
echo "example.com" | postmap -q - /etc/postfix/transport
```

## Advanced Configurations

### Conditional Routing

Route only certain email types through the duplicator:

```
# In transport file
newsletter@example.com  smtp-dup:[127.0.0.1]:2525
alerts@example.com      smtp-dup:[127.0.0.1]:2525
```

### Multiple Duplicators

Use different duplicators for different purposes:

```
# Different duplicators for different domains
internal.com    smtp-dup:[127.0.0.1]:2525
external.com    smtp-dup2:[127.0.0.1]:2526
```

Then add another transport in `master.cf`:

```
smtp-dup2   unix  -       -       n       -       -       smtp
    -o smtp_generic_maps=
    -o smtp_destination_concurrency_limit=2
```

### Fallback Configuration

Configure fallback if duplicator is unavailable:

```
# In main.cf
fallback_transport = smtp
smtp-dup_destination_recipient_limit = 1
smtp-dup_destination_concurrency_limit = 1
```

## Migration Scenarios

### Scenario 1: Gradual Migration

```yaml
# Phase 1: Test with limited domains
# transport:
test.example.com    smtp-dup:[127.0.0.1]:2525

# duplicator config:
relay:
  destination_primary: "oldmail.example.com:25"
  destination_backups:
    - "newmail.example.com:25"
```

```yaml
# Phase 2: Migrate critical domains
# transport:
example.com         smtp-dup:[127.0.0.1]:2525
.example.com        smtp-dup:[127.0.0.1]:2525

# duplicator config stays the same
```

```yaml
# Phase 3: Switch primary destination
# transport stays the same

# duplicator config:
relay:
  destination_primary: "newmail.example.com:25"    # Switched!
  destination_backups:
    - "oldmail.example.com:25"                      # Now backup
```

### Scenario 2: A/B Testing

```yaml
# Split traffic between old and new servers
relay:
  destination_primary: "oldmail.example.com:25"
  destination_backups:
    - "newmail.example.com:25"
    - "monitoring.example.com:25"
```

## Monitoring and Logs

### Postfix Logs

```bash
# Monitor Postfix routing decisions
tail -f /var/log/mail.log | grep "transport=smtp-dup"

# Check for transport errors
grep "smtp-dup" /var/log/mail.log | grep "error"
```

### Duplicator Logs

```bash
# Monitor duplicator activity
journalctl -u gosmtp-dup -f

# Check for relay failures
journalctl -u gosmtp-dup | grep "error"
```

## Troubleshooting

### Common Issues

1. **Transport not applied**:
   ```bash
   # Verify transport map compilation
   postmap -q example.com /etc/postfix/transport

   # Check main.cf configuration
   postconf transport_maps
   ```

2. **Connection refused**:
   ```bash
   # Check if duplicator is listening
   netstat -tulpn | grep 2525

   # Test connection manually
   telnet 127.0.0.1 2525
   ```

3. **Emails not duplicated**:
   ```bash
   # Check duplicator configuration
   journalctl -u gosmtp-dup | grep "destination_primary"

   # Verify backup servers are reachable
   telnet backup.example.com 25
   ```

### Debugging Commands

```bash
# For Option A (Postfix permanent service):
# Check if service is running
netstat -tulpn | grep 2525
ps aux | grep gosmtp-dup

# For Option B (External Service):
# Test SMTP connection to duplicator
echo -e "EHLO test\nQUIT" | nc 127.0.0.1 2525

# General debugging:
# Check Postfix configuration syntax
postfix check

# View active transport maps
postconf -d | grep transport_maps

# Test email routing
echo "test@example.com" | postmap -q - /etc/postfix/transport
```

## Performance Considerations

### Option A (Postfix permanent service)
- **Always Running**: Single permanent process managed by Postfix
- **Better Performance**: No startup overhead per email
- **Single Process**: maxproc=1 is sufficient due to goroutine concurrency
- **Auto-restart**: Postfix will restart if process dies

### Option B (External Service)
- **Concurrency**: Limit `smtp_destination_concurrency_limit` to avoid overwhelming the duplicator
- **Timeouts**: Set appropriate timeouts in both Postfix and the duplicator
- **Queue Management**: Monitor Postfix queues for backup when duplicator is slow

[← Back: Docker Usage](docker.html) | [Next: Home →](index.html)