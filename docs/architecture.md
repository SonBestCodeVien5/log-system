# Kiến trúc hệ thống

## Tổng quan

Hệ thống log tập trung (Centralized Logging Platform) gom log từ nhiều service,
parse và lưu vào Elasticsearch, cung cấp dashboard tìm kiếm và alerting real-time.

Roadmap học, verify và chuẩn bị bảo vệ trong 1 tháng cuối nằm ở
[`docs/one-month-defense-roadmap.md`](one-month-defense-roadmap.md).

---

## Luồng dữ liệu

```
[Services]
  demo-node (Node.js)  -->  /logs/demo-node/app.log  (JSON Lines)
  demo-go   (Go)       -->  /logs/demo-go/app.log    (JSON Lines)
       |
       v
[Filebeat]
  - Tail log files
  - Buffer khi Logstash down (registry)
  - Ship tới Logstash :5044
       |
       v
[Logstash]
  Input:  beats (port 5044)
  Filter: codec json parse → Grok enrich thêm field
  Output: Elasticsearch index logs-YYYY.MM.DD
       |
       v
[Elasticsearch :9200]
  Lưu trữ JSON document
  Index theo ngày
       |
       +-----------> [Go API Server :8080]
       |
       +-----------> [Alerting Engine goroutine]
                          |
                          v WebSocket
                     [Dashboard]
```

---

## Format Log

### JSON Lines — format chính

Demo services ghi mỗi dòng là 1 JSON object hoàn chỉnh:

```json
{"timestamp":"2024-01-15T10:23:11Z","level":"ERROR","service":"demo-node","message":"Payment gateway timeout","metadata":{"order_id":"789"}}
```

**Tại sao JSON Lines thay vì text thô:**
- Logstash chỉ cần `codec => json`, không phụ thuộc Grok để parse format cơ bản
- Ít lỗi parse hơn — JSON sai format dễ debug hơn Grok pattern sai
- Phản ánh thực tế production — các logger hiện đại (Winston, Zap, zerolog) đều output JSON

### Grok — enrich phụ

Sau khi parse JSON, Logstash promote `message` thành `log_message`, rồi dùng Grok
để enrich thêm field từ `log_message`:

```
filter {
  # parse JSON trước
  # Grok enrich thêm nếu log_message có pattern đặc biệt
  grok {
    match => { "log_message" => "(?:%{WORD:error_code}:)?%{GREEDYDATA:error_detail}" }
    tag_on_failure => []  # bỏ qua silently nếu không match
  }
}
```

Grok lỗi không làm mất log — `tag_on_failure => []` đảm bảo log vẫn vào ES.

---

## Các tầng hệ thống

### Trạng thái triển khai hiện tại

| Thành phần | Trạng thái | Bằng chứng trong repo | Cần làm tiếp |
|---|---|---|---|
| Docker Compose stack | Hoàn thành | `docker-compose.yml` định nghĩa ES, Logstash, Filebeat, demo services, API | Verify runtime bằng `docker compose ps` |
| Demo services | Hoàn thành | `services/demo-node/`, `services/demo-go/` ghi JSON Lines | Quan sát file log trong `./logs` |
| Filebeat + Logstash pipeline | Hoàn thành | `filebeat/filebeat.yml`, `logstash/pipeline/logstash.conf` | Verify `logs-*/_count > 0` |
| Elasticsearch storage | Hoàn thành | ES 8.13.0 trong Compose, index `logs-*` | Ghi số liệu count/index sau test |
| Go API server | Hoàn thành | `/api/health`, `/api/logs`, `/api/logs/count`, `/api/alerts/config` | Chạy filter API ở Bước 10.2 |
| Alerting engine | Hoàn thành | Sliding window, dedup, dynamic threshold trong `api-server/alerting/engine.go` | Trigger threshold thấp ở Bước 10.3 |
| Dashboard | Hoàn thành, cần verify runtime | `dashboard/index.html`, `dashboard/app.js`, `dashboard/style.css` | Rebuild API container và mở `http://localhost:8080` |
| Docs + bảo vệ | Đang hoàn thiện | `README.md`, `docs/project-roadmap.md` | Điền số liệu thực tế và demo script |

### Tầng 1 — Demo Services

| Service | Công nghệ | Output |
|---|---|---|
| demo-node | Node.js 20 | `/logs/demo-node/app.log` |
| demo-go | Go 1.22 | `/logs/demo-go/app.log` |

Tỉ lệ log: INFO 60% — WARN 25% — ERROR 15%

### Tầng 2 — Filebeat

Agent nhẹ tail log file và ship tới Logstash. Có registry — không mất log khi restart.

**Tại sao không ship thẳng từ app:**
Filebeat footprint nhỏ (~50MB), chạy được trên mọi server.
Tách concern: app chỉ ghi file, Filebeat lo shipping.

### Tầng 3 — Logstash

Parse và enrich log trước khi lưu vào ES.

```
Input raw JSON:  {"level":"ERROR","message":"Payment failed","service":"demo-node"}
Output enriched: {
  "level":     "ERROR",
  "service":   "demo-node",
  "log_message": "Payment failed",
  "@timestamp": "2024-01-15T10:23:11Z",
  "error_code": "Payment",       ← thêm bởi Grok enrich
  "host":       "log-node-1"     ← thêm bởi Logstash
}
```

### Tầng 4 — Elasticsearch

Lưu trữ log dạng JSON document, index theo ngày.

**Tại sao không dùng MySQL/PostgreSQL:**

| Tiêu chí | Elasticsearch | MySQL |
|---|---|---|
| Index type | Inverted index | B-tree |
| Full-text search | Tự nhiên, nhanh | LIKE query, chậm |
| Time-range | `now-5m` built-in | Tính timestamp thủ công |

### Tầng 5 — Go API Server

REST API + WebSocket, trung gian giữa ES và Dashboard.

**Tại sao Go thay vì Spring Boot:**
- Docker, Kubernetes, Prometheus, Filebeat đều viết bằng Go — đúng domain
- Goroutine native phù hợp alerting concurrent
- RAM 20–50MB so với Spring 300–500MB

---

## Alerting Engine

### Sliding Window
Quét mỗi `ALERT_CHECK_INTERVAL_SECONDS` giây (mặc định 10s),
nhìn lùi `ALERT_WINDOW_SECONDS` giây (mặc định 300s = 5 phút).

```
t=0s:  đếm ERROR [t-300s, t] = 3  → bình thường
t=10s: đếm ERROR [t-300s, t] = 15 → vượt threshold → ALERT
t=20s: đếm ERROR [t-300s, t] = 14 → trong cooldown → bỏ qua
```

Độ trễ phát hiện tối đa = `ALERT_CHECK_INTERVAL_SECONDS`.
Mặc định 10s, có thể điều chỉnh qua `.env`.

### Alert Deduplication

Tránh Alert Fatigue — không gửi alert trùng trong cooldown period.

**Lưu ý implementation:** `shouldAlert` dùng single `Lock` (không tách RLock/Lock)
để đảm bảo check và ghi là atomic, tránh race condition double alert:

```go
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
```

### Dynamic Threshold

Dashboard gửi ngưỡng mới về API, goroutine cập nhật ngay không restart.
Implementation hiện tại dùng một `sync.Mutex` trong `AlertEngine` cho config
và dedup map; cách này giữ check/update đơn giản và đảm bảo phần dedup atomic.
Danh sách WebSocket clients dùng mutex riêng để tách khỏi config alerting.

---

## Hiệu năng — kỳ vọng và cơ chế đảm bảo

Thay vì cam kết con số tuyệt đối, hệ thống thiết kế theo các cơ chế sau.
Con số thực tế sẽ được đo và ghi nhận sau khi hoàn thành tích hợp.

| Mục tiêu | Cơ chế đảm bảo | Ghi chú |
|---|---|---|
| Query log nhanh | ES inverted index + index theo ngày | Đo trên môi trường dev sau tuần 2 |
| Dashboard không chậm theo data | Pagination 20 record/trang, lazy load | Không load toàn bộ log một lúc |
| Alerting phát hiện kịp thời | Check interval = 10s (configurable) | Độ trễ tối đa = interval value |
| Không mất log khi restart | Filebeat registry | Tiếp tục từ điểm dừng |

---

## Cấu trúc thư mục

```
log-system/
├── AGENTS.md
├── docker-compose.yml
├── .env / .env.example
├── filebeat/filebeat.yml
├── logstash/pipeline/logstash.conf
├── elasticsearch/config/elasticsearch.yml
├── services/demo-node/  demo-go/
├── api-server/
│   ├── main.go
│   ├── handlers/logs.go  alerts.go
│   ├── alerting/engine.go
│   └── middleware/cors.go
├── dashboard/index.html  app.js  style.css
├── logs/demo-node/  demo-go/
└── docs/
```
