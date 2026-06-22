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
  - Memory queue khi downstream chậm
  - Registry lưu offset để resume sau restart
  - Ship tới Logstash :5044
       |
       v
[Logstash]
  Input:  beats (port 5044)
  Filter: JSON filter parse → Grok enrich thêm field
  Output: Elasticsearch index logs-YYYY.MM.dd
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

| Thành phần | Trạng thái | Bằng chứng |
|---|---|---|
| Docker Compose stack | Hoàn thành, đã verify | Sáu service healthy trong lần kiểm tra 2026-06-22 |
| Demo services | Hoàn thành, đã verify | Node.js và Go sinh JSON Lines liên tục; dashboard nhận log từ cả hai service |
| Filebeat + Logstash pipeline | Hoàn thành, đã verify | Dữ liệu được parse/enrich và xuất hiện qua API từ index `logs-*` |
| Elasticsearch storage | Hoàn thành, đã verify | API count và filter trả dữ liệu thực từ Elasticsearch |
| Go API server | Hoàn thành, đã verify | Health, logs, filter, count và alert config đã pass E2E |
| Alerting engine | Hoàn thành, đã verify | Incident replay tạo 20 ERROR; API log ghi alert sent |
| Dashboard | Hoàn thành, đã verify | Dashboard load, WebSocket Connected, stats/table/pagination có dữ liệu ngày 2026-06-22 |
| Docs + bảo vệ | Đang đóng gói | Còn clean-clone, screenshot alert banner và rehearsal; xem `report-notes.md` |

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
  "environment": "development",  ← thêm bởi Logstash
  "stack":       "log-system"    ← thêm bởi Logstash
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

Khi không bị cooldown chặn, độ trễ do polling tối đa xấp xỉ
`ALERT_CHECK_INTERVAL_SECONDS`. Mặc định là 10 giây và có thể điều chỉnh qua
biến môi trường. Cooldown có thể làm lần alert tiếp theo đến muộn hơn dù count
vẫn vượt ngưỡng; đây là hành vi chống gửi trùng có chủ đích.

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

## Hiệu năng — cơ chế và số đo dev

Các số dưới đây là phép đo `curl` đơn lẻ trên máy dev ngày 2026-06-17, không
phải benchmark tải nặng và không được dùng như SLA production.

| Mục tiêu | Cơ chế đảm bảo | Kết quả ghi nhận |
|---|---|---|
| Query log nhanh | ES inverted index + index theo ngày | `/api/logs?size=100`: HTTP 200 trong `0.033992s` |
| Dashboard không load toàn bộ data | Pagination 20 record/trang | `/`: HTTP 200 trong `0.012889s` |
| Health check nhẹ | ES cluster health qua Go API | `/api/health`: HTTP 200 trong `0.007250s` |
| Alerting có kiểm soát | Polling configurable + cooldown | Incident replay được ingest; alert quan sát sau ~33s do cooldown trước đó còn active |
| Tiếp tục đọc sau restart | Filebeat registry volume | Có cơ chế trong config; chưa có phép thử restart định lượng riêng |

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
