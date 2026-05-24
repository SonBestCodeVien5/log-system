# AGENTS.md — api-server

Thư mục này chứa Go API server dùng gin framework.
Đọc AGENTS.md ở root để hiểu context tổng thể.

## Package structure

```
api-server/
├── main.go              ← khởi tạo gin, ES client, alerting engine, routes
├── handlers/
│   ├── logs.go          ← handler cho /api/logs và /api/logs/count
│   └── alerts.go        ← handler cho WebSocket /ws/alerts và POST /api/alerts/config
├── alerting/
│   └── engine.go        ← AlertEngine struct, goroutine, sliding window, dedup
└── middleware/
    └── cors.go          ← CORS cho dashboard gọi API
```

## AlertEngine — cách hoạt động

```go
type AlertEngine struct {
    esClient    *elasticsearch.Client
    threshold   int           // đọc/ghi cần RWMutex
    mu          sync.RWMutex
    clients     map[*websocket.Conn]bool  // WebSocket connections
    windowSecs  int
    cooldown    int
    sent        map[string]time.Time  // deduplication tracking
}
```

Goroutine chạy mỗi ALERT_CHECK_INTERVAL_SECONDS:
1. Query ES đếm ERROR trong window vừa qua
2. Nếu count > threshold VÀ chưa gửi trong cooldown → broadcast WebSocket
3. Ghi vào sent map với timestamp

## ES Client khởi tạo

```go
cfg := elasticsearch.Config{
    Addresses: []string{fmt.Sprintf("http://elasticsearch:%s", os.Getenv("ES_PORT"))},
    Username:  "elastic",
    Password:  os.Getenv("ES_PASSWORD"),
}
```
