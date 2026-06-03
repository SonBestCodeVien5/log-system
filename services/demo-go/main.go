// =================================================================
// Demo Service B — Go
// Sinh log JSON Lines liên tục với tỉ lệ INFO/WARN/ERROR ngẫu nhiên
// =================================================================
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

// ---------------------------------------------------------------
// LogEntry — cấu trúc 1 dòng log JSON Lines
// ---------------------------------------------------------------
type LogEntry struct {
	Timestamp string         `json:"timestamp"`
	Level     string         `json:"level"`
	Service   string         `json:"service"`
	Message   string         `json:"message"`
	Metadata  map[string]any `json:"metadata"`
}

// ---------------------------------------------------------------
// Message mẫu theo level
// ---------------------------------------------------------------
var messages = map[string][]string{
	"INFO": {
		"User login successful",
		"Order created successfully",
		"Payment processed",
		"Cache refreshed",
		"Health check passed",
		"Config reloaded",
		"Request completed in 38ms",
	},
	"WARN": {
		"Response time exceeded 500ms",
		"Retry attempt 2/3 for request",
		"Cache miss — falling back to DB",
		"Rate limit approaching threshold",
		"Memory usage at 80%",
	},
	"ERROR": {
		"Database connection refused",
		"Failed to process request: timeout",
		"Service unavailable: payment-api",
		"Panic recovered: nil pointer dereference",
		"Failed to write to storage",
	},
}

// ---------------------------------------------------------------
// Chọn level theo tỉ lệ: INFO 60%, WARN 25%, ERROR 15%
// ---------------------------------------------------------------
func pickLevel() string {
	r := rand.Float64()
	switch {
	case r < 0.60:
		return "INFO"
	case r < 0.85:
		return "WARN"
	default:
		return "ERROR"
	}
}

func pickMessage(level string) string {
	list := messages[level]
	return list[rand.Intn(len(list))]
}

// ---------------------------------------------------------------
// Main
// ---------------------------------------------------------------
func main() {
	serviceName := getEnv("SERVICE_NAME", "demo-go")
	logFile     := getEnv("LOG_FILE", "/var/log/app/app.log")
	intervalMs  := getEnvInt("LOG_INTERVAL_MS", 1500)

	// Đảm bảo thư mục tồn tại
	if err := os.MkdirAll(filepath.Dir(logFile), 0755); err != nil {
		log.Fatalf("failed to create log dir: %v", err)
	}

	// Mở file để append
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}
	defer f.Close()

	ticker   := time.NewTicker(time.Duration(intervalMs) * time.Millisecond)
	sequence := 0
	quit     := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	fmt.Printf("[%s] Starting — writing to %s every %dms\n", serviceName, logFile, intervalMs)

	for {
		select {
		case <-ticker.C:
			sequence++
			level   := pickLevel()
			message := pickMessage(level)

			entry := LogEntry{
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Level:     level,
				Service:   serviceName,
				Message:   message,
				Metadata: map[string]any{
					"pid":      os.Getpid(),
					"sequence": sequence,
				},
			}

			// Ghi JSON Lines
			data, err := json.Marshal(entry)
			if err != nil {
				log.Printf("marshal error: %v", err)
				continue
			}
			if _, err := f.Write(append(data, '\n')); err != nil {
				log.Printf("write error: %v", err)
				continue
			}

			// Stdout để xem trong docker logs
			fmt.Printf("[%s] %s\n", level, message)

		case <-quit:
			fmt.Printf("[%s] Shutting down...\n", serviceName)
			ticker.Stop()
			return
		}
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}
