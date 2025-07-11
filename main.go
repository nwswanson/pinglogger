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

	_ "modernc.org/sqlite"
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

func main() {
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
	db, err := sql.Open("sqlite", dbFile)
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


