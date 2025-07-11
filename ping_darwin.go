//go:build darwin

package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"time"
)

func pingHost(ip string, interval time.Duration) (bool, time.Duration) {
	// Use the system's ping command
	ctx, cancel := context.WithTimeout(context.Background(), interval)
	defer cancel()

	cmd := exec.CommandContext(ctx, "ping", "-c", "1", "-t", fmt.Sprintf("%.0f", interval.Seconds()), ip)
	output, err := cmd.CombinedOutput()

	if err != nil {
		// Ping command failed (e.g., host unreachable, timeout)
		log.Printf("Ping command failed for %s: %v\nOutput: %s", ip, err, string(output))
		return false, 0
	}

	// Parse ping output for RTT
	// Example output line: "round-trip min/avg/max/stddev = 0.040/0.040/0.040/0.000 ms"
	rttRegex := regexp.MustCompile(`time=(\d+\.?\d*)\s*ms`)
	matches := rttRegex.FindStringSubmatch(string(output))

	if len(matches) > 1 {
		rttMs, _ := strconv.ParseFloat(matches[1], 64)
		return true, time.Duration(rttMs * float64(time.Millisecond))
	}

	log.Printf("Could not parse RTT from ping output for %s: %s", ip, string(output))
	return false, 0
}
