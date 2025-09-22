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

### Step 3: Configure Transport Maps

In your `/etc/postfix/main.cf`, add or update:

```
transport_maps = hash:/etc/postfix/transport
```

### Step 4: Create Transport File

Create `/etc/postfix/transport` to define which emails go through the duplicator:

#### Option A: Specific Domains

```
# Route specific domains through the duplicator
example.com     smtp-dup:[127.0.0.1]:2525
.example.com    smtp-dup:[127.0.0.1]:2525
test.com        smtp-dup:[127.0.0.1]:2525
```

#### Option B: All Outbound Email

```
# Route all outbound email through duplicator
*               smtp-dup:[127.0.0.1]:2525
```

#### Option C: Specific Recipients

```
# Route specific recipients
user@example.com        smtp-dup:[127.0.0.1]:2525
admin@example.com       smtp-dup:[127.0.0.1]:2525
```

### Step 5: Update Transport Map and Reload

```bash
# Compile the transport map
postmap /etc/postfix/transport

# Reload Postfix configuration
systemctl reload postfix

# Verify configuration
postfix check
```

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
# Test SMTP connection to duplicator
echo -e "EHLO test\nQUIT" | nc 127.0.0.1 2525

# Check Postfix configuration syntax
postfix check

# View active transport maps
postconf -d | grep transport_maps

# Test email routing
echo "test@example.com" | postmap -q - /etc/postfix/transport
```

## Performance Considerations

- **Concurrency**: Limit `smtp_destination_concurrency_limit` to avoid overwhelming the duplicator
- **Timeouts**: Set appropriate timeouts in both Postfix and the duplicator
- **Queue Management**: Monitor Postfix queues for backup when duplicator is slow

[← Back: Docker Usage](docker.html) | [Next: Home →](index.html)