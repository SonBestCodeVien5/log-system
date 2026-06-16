# AGENTS.md — Log System Project

## Tổng quan
Hệ thống log tập trung (Centralized Logging Platform) cho phép gom log từ nhiều service,
tìm kiếm nhanh, lọc theo mức độ và hiển thị cảnh báo real-time khi có lỗi.

**Chủ sở hữu:** SonBestCodeVien5
**Repo:** https://github.com/SonBestCodeVien5/log-system
**Trạng thái:** Đang phát triển — dự án tốt nghiệp

---

## Kiến trúc tổng thể

```
Demo Services (Node.js + Go)
  → ghi log JSON Lines ra file /logs/**/*.log
  → Filebeat tail file, ship tới Logstash :5044
  → Logstash codec json + Grok enrich → Elasticsearch :9200
  → Go API Server (gin) query ES
  → Dashboard (HTML/JS thuần) hiển thị log
  → Alerting Engine (goroutine) push WebSocket khi ERROR vượt ngưỡng
```

---

## Tech Stack & Version

| Thành phần | Công nghệ | Version |
|---|---|---|
| Log collector | Filebeat | 8.13.0 |
| Log processor | Logstash | 8.13.0 |
| Storage & search | Elasticsearch | 8.13.0 |
| API server | Go + gin | Go 1.22, gin v1.10 |
| Alerting | Go goroutine + gorilla/websocket | — |
| ES client | go-elasticsearch | v8 |
| Demo service A | Node.js | 20 LTS |
| Demo service B | Go | 1.22 |
| Dashboard | HTML + Vanilla JS | — |
| Infrastructure | Docker Compose v2 | — |

---

## Format Log

### JSON Lines — format chính (ưu tiên)
Demo services ghi log dạng JSON Lines, mỗi dòng là 1 JSON object:

```json
{"timestamp":"2024-01-15T10:23:11Z","level":"ERROR","service":"demo-node","message":"Payment gateway timeout","metadata":{"order_id":"789","retry":2}}
{"timestamp":"2024-01-15T10:23:12Z","level":"INFO","service":"demo-go","message":"User login successful","metadata":{"uid":"12345"}}
```

Logstash dùng `codec => json` để parse — không cần Grok cho format cơ bản.

### Grok — enrich phụ
Sau khi parse JSON, Logstash promote `message` thành `log_message`, rồi dùng Grok
để enrich thêm field từ `log_message`:

```
# Ví dụ: tách error_code từ log_message nếu có pattern
grok { match => { "log_message" => "(?:%{WORD:error_code}:)?%{GREEDYDATA:error_detail}" } }
```

Grok lỗi parse không làm mất log — dùng `tag_on_failure => []` để bỏ qua silently.

---

## Go API — Conventions

**Module path:** `github.com/SonBestCodeVien5/log-system/api-server`

**Error handling — ưu tiên:**
```go
// Ưu tiên: return error có context
if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
}
```

**Logging — ưu tiên trong production code:**
```go
// Ưu tiên dùng log.Printf trong production code
log.Printf("[handler] query failed: %v", err)

// fmt.Println chấp nhận trong debug tạm thời, dọn trước khi commit
fmt.Println("debug:", err)
```

**Elasticsearch query:**
```go
query := `{
  "query": {
    "bool": {
      "must": [{ "term": { "level.keyword": "ERROR" } }],
      "filter": { "range": { "@timestamp": { "gte": "now-5m", "lte": "now" } } }
    }
  }
}`
```

**Response format chuẩn:**
```json
{"data": [...], "total": 100, "page": 1, "size": 20}
```

---

## AlertEngine — Conventions

AlertEngine dùng single Lock (không phải RLock/Lock tách biệt) cho `shouldAlert`
để đảm bảo check và ghi là atomic, tránh double alert:

```go
// ĐÚNG — atomic check + write
func (e *AlertEngine) shouldAlert(key string) bool {
    e.mu.Lock()
    defer e.mu.Unlock()

    lastSent, exists := e.sent[key]
    if exists && time.Since(lastSent) < e.cooldown {
        return false
    }
    e.sent[key] = time.Now()
    return true
}

// SAI — race condition tiềm ẩn
func (e *AlertEngine) shouldAlert(key string) bool {
    e.mu.RLock()
    _, exists := e.sent[key]  // goroutine A và B đều đọc "chưa có"
    e.mu.RUnlock()
    // cả 2 goroutine cùng pass → double alert
    e.mu.Lock()
    e.sent[key] = time.Now()
    e.mu.Unlock()
    return true
}
```

---

## Dashboard — Conventions

- Không dùng framework — HTML + Vanilla JS thuần
- Gọi API bằng `fetch()`
- Nhận alert qua WebSocket
- Pagination 20 record/trang — không load toàn bộ log một lúc
- Tất cả trong 3 file: `index.html`, `app.js`, `style.css`

## Agent Context — Conventions

- Ưu tiên lưu phase context tại `.agents/context/features/<feature-slug>/`.
- Nếu `.agents/context` bị read-only trong phiên Codex, dùng `agent-context/features/<feature-slug>/` với cùng file phase (`01-discovery.md` đến `07-blocked.md`).
- Không lưu feature context tạm vào `docs/`; `docs/` chỉ dành cho tài liệu dự án/report.

---

## Rules cứng — không vi phạm

- Không commit file `.env`
- Không hardcode password, port vào code — dùng `os.Getenv()`
- Không bỏ qua error trong Go mà không xử lý hoặc log

## Conventions — ưu tiên tuân theo

- Dùng `log.Printf` thay `fmt.Println` trong production code
- Không thêm dependency mới vào Go mà không cập nhật `go.mod`
- Không đổi format log JSON của demo services khi Logstash đã config xong

---

## Chạy và test

```bash
# Khởi động
docker compose up -d

# Kiểm tra ES
curl http://localhost:9200/_cluster/health

# Kiểm tra log đã vào ES
curl "http://localhost:9200/logs-*/_count"

# API health check
curl http://localhost:8080/api/health

# Dev mode Go API
cd api-server && go run main.go
```

---

## Endpoints API

| Method | Path | Mô tả |
|---|---|---|
| GET | `/api/logs` | Danh sách log, filter level/app/time/search |
| GET | `/api/logs/count` | Đếm log theo level |
| GET | `/api/health` | Health check |
| WebSocket | `/ws/alerts` | Real-time alert stream |
| POST | `/api/alerts/config` | Cập nhật threshold động |
