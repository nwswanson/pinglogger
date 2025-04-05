
# pingdaemon

**A minimal ping logger for local network diagnostics on macOS.**

## Overview

`pingdaemon` is a small CLI tool that pings a target IP address at a regular interval and logs results to a local SQLite database. It's intended for users who want a simple way to collect basic network reachability data over time.

## Features

- Logs timestamp, success/failure, and RTT (round-trip time) for each ping
- Uses SQLite for persistent, queryable storage
- Works well on laptops (survives sleep/wake thanks to WAL mode)
- Should be low resource usage
- Easy to start/stop manually or via a launch agent

## Usage

1. **Build the tool**:

   ```bash
   go build -o pingdaemon
   ```

2. **Run it with sudo** (required for raw ICMP packets):

   ```bash
   sudo ./pingdaemon --ip 1.1.1.1 --interval 10s --db ~/pinglog.db
   ```

3. Leave it running, or wrap it in a `launchd` agent if you want it to start automatically.

## Why Use It?

It’s helpful when:
- You're trying to debug intermittent connection drops
- You want a lightweight log of when a machine could or couldn’t reach a known IP
- You're troubleshooting weird behavior and want to rule out basic connectivity issues

## Schema

```sql
CREATE TABLE pings (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	timestamp DATETIME,
	success BOOLEAN,
	rtt REAL -- in seconds
);
```

## Requirements

- macOS (uses raw ICMP sockets)
- Go 1.20+
- Root privileges to run
