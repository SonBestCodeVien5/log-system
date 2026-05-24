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
  → ghi log ra file /logs/**/*.log
  → Filebeat tail file, ship tới Logstash :5044
  → Logstash Grok parse → Elasticsearch :9200
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

## Cấu trúc thư mục

```
log-system/
├── AGENTS.md                   ← file này
├── docker-compose.yml          ← khởi động toàn bộ hệ thống
├── .env                        ← config (không commit)
├── .env.example                ← template config
├── filebeat/
│   └── filebeat.yml            ← tail log files, ship to Logstash
├── logstash/
│   ├── pipeline/logstash.conf  ← Grok parse, enrich, output to ES
│   └── config/logstash.yml
├── elasticsearch/
│   └── config/elasticsearch.yml
├── services/
│   ├── demo-node/              ← Node.js service sinh log liên tục
│   └── demo-go/                ← Go service sinh log liên tục
├── api-server/                 ← Go + gin REST API + WebSocket
│   ├── main.go
│   ├── handlers/
│   │   ├── logs.go             ← GET /api/logs
│   │   └── alerts.go           ← WebSocket /ws/alerts
│   ├── alerting/
│   │   └── engine.go           ← Sliding Window + Deduplication + Dynamic Threshold
│   └── middleware/
│       └── cors.go
├── dashboard/
│   ├── index.html
│   ├── app.js
│   └── style.css
└── logs/                       ← volume mount, Filebeat đọc từ đây
    ├── demo-node/
    └── demo-go/
```

---

## Cấu hình qua biến môi trường

Mọi config đều đọc từ `.env` — **không bao giờ hardcode** giá trị vào code.

```
ES_VERSION=8.13.0
ES_PORT=9200
ES_PASSWORD=changeme123
LOGSTASH_PORT=5044
API_PORT=8080
ALERT_THRESHOLD=10
ALERT_WINDOW_SECONDS=300
ALERT_COOLDOWN_SECONDS=60
ALERT_CHECK_INTERVAL_SECONDS=10
```

---

## Format Log chuẩn

Tất cả demo service phải ghi log theo format này để Grok parse được:

```
[2024-01-15T10:23:11Z] [ERROR] [demo-node] Payment gateway timeout after 30s
[2024-01-15T10:23:12Z] [INFO]  [demo-go]   User login successful uid=12345
[2024-01-15T10:23:13Z] [WARN]  [demo-node] Retry attempt 2/3 for order=789
```

Pattern: `[timestamp] [LEVEL] [service-name] message`

---

## Go API — Conventions

**Module path:** `github.com/SonBestCodeVien5/log-system/api-server`

**Error handling:** luôn return error, không panic trừ main():
```go
// ĐÚNG
if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
}

// SAI — không dùng
if err != nil {
    panic(err)
}
```

**Elasticsearch query:** dùng `go-elasticsearch` v8, query bằng JSON string:
```go
query := `{
  "query": {
    "bool": {
      "must": [
        { "term": { "level.keyword": "ERROR" } }
      ],
      "filter": {
        "range": { "@timestamp": { "gte": "now-5m", "lte": "now" } }
      }
    }
  }
}`
```

**Response format chuẩn:**
```json
{
  "data": [...],
  "total": 100,
  "page": 1,
  "size": 20
}
```

**Alerting engine** chạy trong goroutine riêng, dùng `sync.RWMutex` cho shared state:
- Sliding Window: quét lùi `ALERT_WINDOW_SECONDS` giây mỗi `ALERT_CHECK_INTERVAL_SECONDS` giây
- Deduplication: không bắn alert trùng trong vòng `ALERT_COOLDOWN_SECONDS` giây
- Dynamic Threshold: nhận config mới qua WebSocket message từ dashboard

---

## Dashboard — Conventions

- **Không dùng framework** (không React, không Vue) — HTML + Vanilla JS thuần
- Gọi API bằng `fetch()`
- Nhận alert qua `WebSocket`
- CSS tự viết, không dùng Tailwind hay Bootstrap
- Tất cả trong 3 file: `index.html`, `app.js`, `style.css`

---

## Docker Compose — Conventions

- Mọi service đều có `restart: always`
- Mọi service đều có `healthcheck`
- ES và Logstash phải healthy trước khi Filebeat và API start (`depends_on: condition: service_healthy`)
- Network: tất cả trong cùng network `log-network`
- Volume logs mount từ `./logs` vào container

---

## Quy tắc KHÔNG được làm

- **Không hardcode** password, port, hostname vào code — dùng `os.Getenv()`
- **Không commit** file `.env`
- **Không dùng** `fmt.Println` trong Go API — dùng `log.Printf`
- **Không bỏ qua** error trong Go — luôn check và handle
- **Không thêm** dependency mới vào Go mà không cập nhật `go.mod`
- **Không đổi** format log của demo services — Grok pattern phụ thuộc vào format này

---

## Chạy và test

```bash
# Khởi động toàn bộ hệ thống
docker compose up -d

# Kiểm tra ES đã sẵn sàng
curl http://localhost:9200/_cluster/health

# Xem log của từng service
docker compose logs -f logstash
docker compose logs -f filebeat

# Chạy Go API local (development)
cd api-server
go run main.go

# Kiểm tra log đã vào ES chưa
curl "http://localhost:9200/logs-*/_count"
```

---

## Endpoints API

| Method | Path | Mô tả |
|---|---|---|
| GET | `/api/logs` | Lấy danh sách log, filter theo level/app/time |
| GET | `/api/logs/count` | Đếm log theo level trong khoảng thời gian |
| GET | `/api/health` | Health check |
| WebSocket | `/ws/alerts` | Real-time alert stream |
| POST | `/api/alerts/config` | Cập nhật ngưỡng alert động |

**Query params cho `/api/logs`:**
```
?level=ERROR
?app=demo-node
?from=2024-01-15T10:00:00Z
?to=2024-01-15T11:00:00Z
?q=payment+timeout        ← full-text search
?page=1&size=20
```
