// =================================================================
// Alerting Engine — Sliding Window + Deduplication + Dynamic Threshold
// =================================================================
package alerting

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
	"fmt"
	
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gorilla/websocket"
)

// ---------------------------------------------------------------
// AlertMessage — gửi qua WebSocket về dashboard
// ---------------------------------------------------------------
type AlertMessage struct {
	Type      string    `json:"type"`
	Count     int64     `json:"count"`
	Threshold int       `json:"threshold"`
	Window    string    `json:"window"`
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
}

// ---------------------------------------------------------------
// AlertConfig — nhận từ dashboard qua POST /api/alerts/config
// ---------------------------------------------------------------
type AlertConfig struct {
	Threshold      int `json:"threshold"`
	WindowSeconds  int `json:"window_seconds"`
	CooldownSeconds int `json:"cooldown_seconds"`
}

// ---------------------------------------------------------------
// AlertEngine
// ---------------------------------------------------------------
type AlertEngine struct {
	es      *elasticsearch.Client

	// Config — đọc/ghi cần mu
	mu              sync.Mutex
	threshold       int
	windowSeconds   int
	cooldownSeconds int
	checkInterval   time.Duration

	// Deduplication — track thời gian gửi alert gần nhất
	// Dùng cùng mu để đảm bảo check+write là atomic
	sent map[string]time.Time

	// WebSocket clients
	clientsMu sync.RWMutex
	clients   map[*websocket.Conn]bool
}

// ---------------------------------------------------------------
// NewEngine — khởi tạo từ biến môi trường
// ---------------------------------------------------------------
func NewEngine(es *elasticsearch.Client) *AlertEngine {
	return &AlertEngine{
		es:              es,
		threshold:       getEnvInt("ALERT_THRESHOLD", 10),
		windowSeconds:   getEnvInt("ALERT_WINDOW_SECONDS", 300),
		cooldownSeconds: getEnvInt("ALERT_COOLDOWN_SECONDS", 60),
		checkInterval:   time.Duration(getEnvInt("ALERT_CHECK_INTERVAL_SECONDS", 10)) * time.Second,
		sent:            make(map[string]time.Time),
		clients:         make(map[*websocket.Conn]bool),
	}
}

// ---------------------------------------------------------------
// Run — goroutine chính, chạy Sliding Window loop
// ---------------------------------------------------------------
func (e *AlertEngine) Run() {
	ticker := time.NewTicker(e.checkInterval)
	defer ticker.Stop()

	log.Printf("[alerting] started — interval=%v threshold=%d window=%ds",
		e.checkInterval, e.threshold, e.windowSeconds)

	for range ticker.C {
		e.check()
	}
}

// ---------------------------------------------------------------
// check — đếm ERROR trong sliding window, trigger nếu vượt ngưỡng
// ---------------------------------------------------------------
func (e *AlertEngine) check() {
	// Đọc config an toàn
	e.mu.Lock()
	threshold     := e.threshold
	windowSeconds := e.windowSeconds
	e.mu.Unlock()

	// Query ES đếm ERROR trong window
	count, err := e.countErrors(windowSeconds)
	if err != nil {
		log.Printf("[alerting] ES query error: %v", err)
		return
	}

	if count <= int64(threshold) {
		return
	}

	// Vượt ngưỡng — check dedup
	alertKey := "error_spike"
	if !e.shouldAlert(alertKey) {
		return
	}

	// Broadcast alert
	msg := AlertMessage{
		Type:      "error_spike",
		Count:     count,
		Threshold: threshold,
		Window:    formatDuration(windowSeconds),
		Timestamp: time.Now().UTC(),
		Message: fmt.Sprintf(
			"%d errors in last %s (threshold: %d)",
			count, formatDuration(windowSeconds), threshold,
		),
	}

	e.broadcast(msg)
	log.Printf("[alerting] alert sent — count=%d threshold=%d", count, threshold)
}

// ---------------------------------------------------------------
// shouldAlert — check + write atomic để tránh double alert
// ---------------------------------------------------------------
func (e *AlertEngine) shouldAlert(key string) bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	cooldown := time.Duration(e.cooldownSeconds) * time.Second
	if lastSent, exists := e.sent[key]; exists {
		if time.Since(lastSent) < cooldown {
			return false
		}
	}
	e.sent[key] = time.Now()
	return true
}

// ---------------------------------------------------------------
// UpdateConfig — Dynamic Threshold, gọi từ HTTP handler
// ---------------------------------------------------------------
func (e *AlertEngine) UpdateConfig(cfg AlertConfig) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if cfg.Threshold > 0 {
		e.threshold = cfg.Threshold
	}
	if cfg.WindowSeconds > 0 {
		e.windowSeconds = cfg.WindowSeconds
	}
	if cfg.CooldownSeconds > 0 {
		e.cooldownSeconds = cfg.CooldownSeconds
	}

	log.Printf("[alerting] config updated — threshold=%d window=%ds cooldown=%ds",
		e.threshold, e.windowSeconds, e.cooldownSeconds)
}

// ---------------------------------------------------------------
// GetConfig — đọc config hiện tại
// ---------------------------------------------------------------
func (e *AlertEngine) GetConfig() AlertConfig {
	e.mu.Lock()
	defer e.mu.Unlock()
	return AlertConfig{
		Threshold:       e.threshold,
		WindowSeconds:   e.windowSeconds,
		CooldownSeconds: e.cooldownSeconds,
	}
}

// ---------------------------------------------------------------
// RegisterClient / UnregisterClient — quản lý WebSocket connections
// ---------------------------------------------------------------
func (e *AlertEngine) RegisterClient(conn *websocket.Conn) {
	e.clientsMu.Lock()
	e.clients[conn] = true
	total := len(e.clients)
	e.clientsMu.Unlock()
	log.Printf("[alerting] client connected, total=%d", total)
}

func (e *AlertEngine) UnregisterClient(conn *websocket.Conn) {
	e.clientsMu.Lock()
	delete(e.clients, conn)
	total := len(e.clients)
	e.clientsMu.Unlock()
	log.Printf("[alerting] client disconnected, total=%d", total)
}

// ---------------------------------------------------------------
// broadcast — gửi alert tới tất cả WebSocket clients
// Conn nào ghi lỗi sẽ bị đóng và xóa khỏi danh sách, tránh tích lũy
// kết nối chết (peer đã disconnect nhưng read goroutine chưa kịp gỡ).
// ---------------------------------------------------------------
func (e *AlertEngine) broadcast(msg AlertMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("[alerting] marshal error: %v", err)
		return
	}

	// Snapshot danh sách conn dưới RLock để không giữ lock trong khi
	// gọi WriteMessage (có thể block bởi network).
	e.clientsMu.RLock()
	conns := make([]*websocket.Conn, 0, len(e.clients))
	for c := range e.clients {
		conns = append(conns, c)
	}
	e.clientsMu.RUnlock()

	var dead []*websocket.Conn
	for _, conn := range conns {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Printf("[alerting] write error, dropping client: %v", err)
			dead = append(dead, conn)
		}
	}

	if len(dead) == 0 {
		return
	}

	e.clientsMu.Lock()
	for _, conn := range dead {
		delete(e.clients, conn)
		_ = conn.Close()
	}
	e.clientsMu.Unlock()
}

// countErrorsTimeout — timeout cho query đếm trong sliding window.
// Ngắn hơn check interval mặc định (10s) để không chèn lên vòng kế.
const countErrorsTimeout = 5 * time.Second

// ---------------------------------------------------------------
// countErrors — query ES đếm ERROR trong window seconds gần nhất
// ---------------------------------------------------------------
func (e *AlertEngine) countErrors(windowSeconds int) (int64, error) {
	window := fmt.Sprintf("now-%ds", windowSeconds)

	query := map[string]any{
		"query": map[string]any{
			"bool": map[string]any{
				"must": []map[string]any{
					{"term": map[string]any{"level.keyword": "ERROR"}},
				},
				"filter": []map[string]any{
					{"range": map[string]any{
						"@timestamp": map[string]string{
							"gte": window,
							"lte": "now",
						},
					}},
				},
			},
		},
	}

	bodyBytes, _ := json.Marshal(query)

	ctx, cancel := context.WithTimeout(context.Background(), countErrorsTimeout)
	defer cancel()

	res, err := e.es.Count(
		e.es.Count.WithContext(ctx),
		e.es.Count.WithIndex("logs-*"),
		e.es.Count.WithBody(bytes.NewReader(bodyBytes)),
	)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return 0, fmt.Errorf("es count error: %s", res.Status())
	}

	var result map[string]any
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return 0, err
	}

	if count, ok := result["count"].(float64); ok {
		return int64(count), nil
	}
	return 0, nil
}

// ---------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------
func formatDuration(seconds int) string {
	if seconds < 60 {
		return fmt.Sprintf("%ds", seconds)
	}
	return fmt.Sprintf("%dm", seconds/60)
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}
