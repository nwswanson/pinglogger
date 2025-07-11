
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

### Building the Tool

`pingdaemon` can be built for your current operating system or cross-compiled for other supported platforms.

#### Build for your current OS

To build the executable for the operating system you are currently on:

```bash
go build -o pingdaemon
```

#### Cross-compile for macOS

If you are building on a Linux machine and targeting macOS, you need to cross-compile. Specify the target operating system (`GOOS=darwin`) and the target architecture (`GOARCH`).

*   **For Intel-based Macs:**
    ```bash
    GOOS=darwin GOARCH=amd64 go build -o pingdaemon_macos
    ```
*   **For Apple Silicon (M1/M2/etc.) Macs:**
    ```bash
    GOOS=darwin GOARCH=arm64 go build -o pingdaemon_macos
    ```
The output executable will be named `pingdaemon_macos`.

### Running the Tool

Once built, you can run `pingdaemon` with various command-line flags:

*   `--ip`: IP address to ping (default: `8.8.8.8`)
*   `--interval`: Ping interval (default: `5s`)
*   `--db`: SQLite DB file (default: `pings.db`)

**Example:**

*   **On Linux (requires sudo):**
    ```bash
    sudo ./pingdaemon --ip 1.1.1.1 --interval 10s --db ~/pinglog.db
    ```
*   **On macOS (no sudo required):**
    ```bash
    ./pingdaemon_macos --ip 1.1.1.1 --interval 10s --db ~/pinglog.db
    ```

Leave it running, or wrap it in a `launchd` agent (macOS) or `systemd` service (Linux) if you want it to start automatically.

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

- Go 1.24+

### Platform Compatibility

`pingdaemon` is designed to work on both Linux and macOS.

*   **Linux:** Uses raw ICMP sockets and requires elevated privileges (`sudo`) to run.
*   **macOS:** Uses the native `ping` command via `os/exec` and does *not* require elevated privileges.
