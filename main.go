package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-ping/ping"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

var (
	ip       string
	interval time.Duration
	dbFile   string
)

type PingResult struct {
	Timestamp time.Time
	Success   bool
	RTT       time.Duration
}

var rootCmd = &cobra.Command{
	Use:   "pingdaemon",
	Short: "A simple ping daemon that logs to SQLite",
	Run: func(cmd *cobra.Command, args []string) {
		runDaemon()
	},
}

func isPrivileged() bool {
	return os.Geteuid() == 0
}

func main() {
	if !isPrivileged() {
		log.Fatal("This program must be run with elevated privileges (sudo) to send ICMP packets.")
	}
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&ip, "ip", "8.8.8.8", "IP address to ping")
	rootCmd.PersistentFlags().DurationVar(&interval, "interval", 5*time.Second, "Ping interval")
	rootCmd.PersistentFlags().StringVar(&dbFile, "db", "pings.db", "SQLite DB file")
}

func runDaemon() {
	db := openDB()
	defer db.Close()

	createSchema(db)

	results := make(chan PingResult, 100)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle signals for graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Writer goroutine
	go func() {
		for res := range results {
			_, err := db.Exec(
				"INSERT INTO pings (timestamp, success, rtt) VALUES (?, ?, ?)",
				res.Timestamp, res.Success, res.RTT.Seconds(),
			)
			if err != nil {
				log.Printf("DB insert failed: %v", err)
			} else {
				fmt.Printf("%s | Success: %t | RTT: %v\n", res.Timestamp.Format(time.RFC3339), res.Success, res.RTT)
			}
		}
	}()

	sem := make(chan struct{}, 1) // Limit concurrent pings
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

loop:
	for {
		select {
		case <-ticker.C:
			sem <- struct{}{}
			go func() {
				defer func() { <-sem }()
				success, rtt := pingHost(ip, interval)
				results <- PingResult{
					Timestamp: time.Now(),
					Success:   success,
					RTT:       rtt,
				}
			}()
		case <-sigs:
			break loop
		case <-ctx.Done():
			break loop
		}
	}

	// Cleanup
	cancel()
	close(results)
	time.Sleep(500 * time.Millisecond) // Give writer time to flush (optional)
}

func openDB() *sql.DB {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("PRAGMA journal_mode=WAL;")
	if err != nil {
		log.Fatalf("Failed to enable WAL mode: %v", err)
	}
	return db
}

func createSchema(db *sql.DB) {
	schema := `
	CREATE TABLE IF NOT EXISTS pings (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp DATETIME,
		success BOOLEAN,
		rtt REAL
	);`
	if _, err := db.Exec(schema); err != nil {
		log.Fatal(err)
	}
}

func pingHost(ip string, interval time.Duration) (bool, time.Duration) {
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
