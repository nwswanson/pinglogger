## Plan: Address macOS Compatibility for `pinglogger`

The goal is to ensure the `pinglogger` application can reliably send ICMP pings on macOS, specifically by resolving issues related to raw socket access and privileged operations.

### Task 1: Investigate `go-ping` macOS Compatibility

*   **Description:** Research the `github.com/go-ping/ping` library's documentation, GitHub issues, and source code to understand its macOS compatibility, particularly regarding privileged mode and raw socket usage. Look for recommended practices or known workarounds for macOS.
*   **Acceptance Criteria:**
    *   Identify if `go-ping` has a documented, reliable way to operate in privileged mode on macOS.
    *   Determine if there's a non-privileged mode that could be used, and if so, its limitations.
    *   Gather information on common pitfalls or alternative approaches suggested by the library's maintainers or community for macOS.
*   **Check-off:** [x]

### Task 2: Evaluate Alternative Ping Mechanisms (if Task 1 is inconclusive or negative)

*   **Description:** If `go-ping` proves problematic or unsuitable for macOS privileged operations, explore alternative methods for sending ICMP pings.
    *   **Option A: System `ping` utility:** Investigate using `os/exec` to call the native `ping` command-line utility available on macOS. This is generally reliable but less performant and "Go-native."
    *   **Option B: Alternative Go library:** Search for other Go libraries specifically designed for ICMP pinging that might have better macOS support or different approaches to raw socket handling.
*   **Acceptance Criteria:**
    *   For Option A: Confirm the `ping` command's output format is parseable and can provide the necessary success/RTT information.
    *   For Option B: Identify at least one viable alternative Go library with a clear indication of macOS compatibility.
*   **Check-off:** [x]

### Task 3: Implement Multi-Platform Ping Solution using Build Tags

*   **Description:** Refactor the `pingHost` function to use Go build tags for platform-specific implementations.
    *   Create `ping_linux.go` for Linux-specific ping logic (using `go-ping` with privileged mode).
    *   Create `ping_darwin.go` for macOS-specific ping logic (using `os/exec` to call the native `ping` command).
    *   The `main.go` file will contain a generic `pingHost` function that calls the platform-specific implementation.
*   **Acceptance Criteria:**
    *   `ping_linux.go` and `ping_darwin.go` files are created with appropriate build tags (`//go:build linux` and `//go:build darwin`).
    *   The `pingHost` function in `main.go` is refactored to call the platform-specific implementation.
    *   The `github.com/go-ping/ping` import is conditionally used only in `ping_linux.go`.
    *   The `os/exec`, `regexp`, and `strconv` imports are conditionally used only in `ping_darwin.go`.
    *   The `isPrivileged()` check is moved to `ping_linux.go` and applied only when building for Linux.
    *   The application compiles successfully for both Linux and macOS targets.
    *   Each platform-specific implementation correctly returns `success` (boolean) and `rtt` (`time.Duration`).
*   **Check-off:** [x]

### Task 4: Test on macOS (User Action Required)

*   **Description:** Provide instructions for the user to build and run the modified `pinglogger` application on a macOS environment to verify its functionality.
*   **Acceptance Criteria:**
    *   The user confirms that the `pinglogger` daemon runs successfully on macOS.
    *   Pings are sent, results are logged to the SQLite database, and output to the console as expected.
    *   The application handles graceful shutdown via `SIGINT`/`SIGTERM` correctly.
*   **Check-off:** [x]
