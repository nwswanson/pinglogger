## Plan: Cross-Compiling `pinglogger` for macOS from Linux

### Objective
To successfully build the `pinglogger` application on a Linux environment such that the resulting executable can run natively on macOS.

### Understanding Go Cross-Compilation
Go's build system supports cross-compilation out-of-the-box. This means you can compile an executable for a different operating system and architecture than the one you are currently running on. This is achieved by setting two environment variables before running `go build`:

*   `GOOS`: Specifies the target operating system (e.g., `linux`, `windows`, `darwin`). This variable is crucial for Go's build tags, allowing different source files or code blocks to be included based on the target OS.
*   `GOARCH`: Specifies the target architecture (e.g., `amd64`, `arm64`).

For macOS targets, `GOOS` should be set to `darwin`. When `GOOS` is set, Go automatically includes files with corresponding build tags (e.g., `//go:build darwin` for macOS-specific code, `//go:build linux` for Linux-specific code), ensuring the correct platform-specific logic is compiled into the final executable.

### Proposed Steps

#### Step 1: Determine Target macOS Architecture
Before compiling, we need to know whether the target macOS machine is Intel-based or Apple Silicon-based.

*   **Action:** User to specify the `GOARCH` for their macOS target:
    *   `amd64` (for Intel Macs)
    *   `arm64` (for Apple Silicon Macs)

#### Step 2: Cross-Compile the Application
Once the target architecture is known, the `go build` command will be executed with the appropriate environment variables.

*   **Action:** The application has been successfully cross-compiled for `darwin/arm64`.
    ```bash
    GOOS=darwin GOARCH=arm64 go build -o pinglogger_macos
    ```
    *   **Note:** The output executable is named `pinglogger_macos` and is located in the project root directory.

#### Step 3: Transfer the Executable to macOS
After successful compilation, the `pinglogger_macos` executable needs to be transferred to the target macOS machine.

*   **Action:** Transfer the `pinglogger_macos` file from your Linux build environment to your macOS machine using a method like `scp`, `rsync`, or a shared drive.

#### Step 4: Test on macOS
Finally, the cross-compiled application needs to be tested on the macOS machine to verify its functionality.

*   **Action:** On the macOS machine, navigate to the directory where you transferred `pinglogger_macos` and make it executable, then run it:
    ```bash
    chmod +x pinglogger_macos
    ./pinglogger_macos
    ```
*   **Acceptance Criteria:**
    *   The `pinglogger_macos` daemon runs successfully on macOS without requiring `sudo`.
    *   Pings are sent, results are logged to the SQLite database (`pings.db`), and output to the console as expected.
    *   The application handles graceful shutdown via `SIGINT`/`SIGTERM` correctly.
