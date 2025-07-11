//go:build linux

package main

import (
	"log"
	"os"
	"time"

	"github.com/go-ping/ping"
)

func isPrivileged() bool {
	return os.Geteuid() == 0
}

func pingHost(ip string, interval time.Duration) (bool, time.Duration) {
	if !isPrivileged() {
		log.Printf("This program must be run with elevated privileges (sudo) to send ICMP packets on Linux.")
		return false, 0
	}

	pinger, err := ping.NewPinger(ip)
	pinger.SetPrivileged(true)
	if err != nil {
		log.Printf("Ping setup failed: %v", err)
		return false, 0
	}
	pinger.Count = 1
	pinger.Timeout = interval * 80 / 100 // e.g., 80% of interval
	if err := pinger.Run(); err != nil {
		log.Printf("Ping run failed: %v", err)
		return false, 0
	}
	stats := pinger.Statistics()
	return stats.PacketsRecv > 0, stats.AvgRtt
}
