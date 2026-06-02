# Tài liệu tri thức — Dự án Log System
> Tổng hợp toàn bộ phân tích, quyết định kỹ thuật và kiến thức nền
> trước khi bắt đầu triển khai. Dùng làm tài liệu tham chiếu khi bảo vệ.

---

## 1. Bối cảnh và yêu cầu gốc

### Yêu cầu chức năng
- **Log Collector** — Gom log từ nhiều service gửi về Elasticsearch tập trung
- **Log Viewer Dashboard** — Hiển thị danh sách log, lọc theo mức độ (INFO/WARN/ERROR) hoặc tên ứng dụng
- **Alerting** — Cảnh báo khi hệ thống có nhiều ERROR log

### Yêu cầu kỹ thuật
- Hệ thống hoạt động đúng với tài liệu giải pháp
- Đảm bảo hiệu năng
- Sẵn sàng triển khai trên môi trường production

### Công nghệ yêu cầu gốc
Elasticsearch, Logstash, Spring Java

### Phạm vi
1 tháng, 1 người, bảo vệ trước hội đồng

---

## 2. Phân tích khả thi

### Người thực hiện
- Có nền Node.js và Go (không sâu)
- Quen thiết kế và nghiệp vụ hơn là code thuần
- Từng dùng vibe code (AI-assisted coding)
- Chưa có kinh nghiệm production deployment
- Background: CRUD web chạy localhost

### Đánh giá tổng thể
**Hoàn toàn khả thi** trong 1 tháng với điều kiện:
- Hiểu rõ từng quyết định kỹ thuật (không chỉ làm cho chạy)
- Viết tài liệu song song với code
- Định nghĩa "hiệu năng" thành con số đo được
- Dùng AI để *giải thích* code, không chỉ *viết* code

### Cái thực sự mới cần học
| Kiến thức | Thời gian ước tính | Mức độ khó |
|---|---|---|
| Docker Compose + networking | 3 ngày | Trung bình |
| Elasticsearch basics | 4 ngày | Trung bình |
| Logstash Grok pattern | 2 ngày | Thấp (config) |
| Filebeat config | 1 ngày | Thấp |
| Go + gin + goroutine | 3 ngày | Thấp (đã có nền) |
| WebSocket cơ bản | 1 ngày | Thấp |

### Cái KHÔNG cần học từ đầu
- REST API, JSON handling → đã có từ CRUD
- Kết nối backend với database → ES giống DB, chỉ khác cú pháp query
- Dựng giao diện filter danh sách → đã làm rồi

---

## 3. Kiến trúc hệ thống

### Luồng dữ liệu tổng thể
```
[Demo Services]
  Node.js / Go  →  ghi file  →  /logs/**/*.log
                                      |
                              [Filebeat - Collector]
                               tail + buffer + ship
                                      |
                                   :5044
                              [Logstash - Processor]
                               Grok parse → JSON
                                      |
                                   :9200
                              [Elasticsearch - Storage]
                               index: logs-YYYY.MM.DD
                                      |
                          +-----------+-----------+
                          |                       |
                  [Go API Server]        [Alerting Engine]
                   gin :8080              goroutine nền
                   REST + WS              Sliding Window
                          |               Deduplication
                          |               Dyn Threshold
                          |                       |
                          +----------+------------+
                                     |
                               [Dashboard]
                               HTML + JS thuần
                               filter + alert UI
```

### 5 tầng và trách nhiệm

**Tầng 1 — Sources**
Các service sinh ra log. Với dự án này dùng 2 demo service (Node.js + Go) sinh log giả lập INFO/WARN/ERROR theo tỉ lệ ngẫu nhiên.

**Tầng 2 — Filebeat (Collector)**
Agent nhẹ chạy cùng service, tail log file, buffer khi downstream down, ship tới Logstash. Có *registry* — nhớ đã đọc đến dòng nào, không mất log khi restart.

**Tầng 3 — Logstash (Processor)**
Nhận raw text từ Filebeat, dùng Grok pattern parse thành JSON có cấu trúc. Không cần viết code — chỉ viết file config `.conf`.

**Tầng 4 — Elasticsearch (Storage)**
Lưu trữ log dạng JSON document, index theo ngày. Full-text search cực nhanh nhờ inverted index. ILM policy tự xóa index cũ.

**Tầng 5 — Go API + Dashboard**
REST API query ES, WebSocket push alert. Dashboard HTML thuần gọi API và nhận WebSocket.

---

## 4. Quyết định công nghệ

### Go thay vì Spring Java

**Lý do kỹ thuật:**

| Tiêu chí | Spring Boot | Go |
|---|---|---|
| Startup time | 10–30 giây | < 1 giây |
| RAM khi chạy | 300–500 MB | 20–50 MB |
| Docker image | ~500 MB | ~15 MB |
| Concurrency | @Scheduled annotation | Goroutine native |
| Setup complexity | Cao (DI, Bean, pom.xml) | Thấp (go mod) |
| Thời gian học thêm | 2–3 tuần | 3–5 ngày |
| Giải thích khi bảo vệ | Khó (có annotation magic) | Dễ (tường minh) |

**Lý do domain:**
Docker, Kubernetes, Prometheus, Filebeat đều viết bằng Go.
Hệ thống log là infrastructure tool — Go là lựa chọn đúng domain.

**2 câu trả lời khi thầy hỏi:**
> *"Em chọn Go vì đây là ngôn ngữ chuẩn của infrastructure tooling —
> Docker, Kubernetes, Prometheus đều viết bằng Go. Hơn nữa goroutine
> cho phép implement Sliding Window và Alert Deduplication tự nhiên hơn
> nhiều so với Spring Scheduler, và RAM nhỏ hơn 10 lần giúp hệ thống
> chạy ổn định khi ES và Logstash cũng đang chạy cùng lúc."*

### Stack cuối cùng

| Thành phần | Công nghệ | Version |
|---|---|---|
| Log collector | Filebeat | 8.13.0 |
| Log processor | Logstash | 8.13.0 |
| Storage & search | Elasticsearch | 8.13.0 |
| API server | Go + gin | Go 1.22 |
| Alerting | goroutine + gorilla/websocket | — |
| ES client | go-elasticsearch | v8 |
| Demo service A | Node.js | 20 LTS |
| Demo service B | Go | 1.22 |
| Dashboard | HTML + Vanilla JS | — |
| Infrastructure | Docker Compose | v2 |

---

## 5. Tính năng nâng cao — Điểm kỹ thuật bảo vệ

### 5.1 Sliding Window
**Vấn đề giải quyết:** Poll cố định mỗi 5 phút dễ bỏ sót spike lỗi ngắn.

**Cách làm:** Goroutine chạy mỗi 10 giây, mỗi lần quét lùi về 5 phút trước.

```
t=0s:  đếm ERROR [t-300s, t] = 3  → bình thường
t=10s: đếm ERROR [t-300s, t] = 15 → ALERT!
t=20s: đếm ERROR [t-300s, t] = 14 → dedup, không bắn lại
```

**Implement:** `time.NewTicker` + ES range query `now-5m`.

### 5.2 Alert Deduplication
**Vấn đề giải quyết:** Alert Fatigue — khi hệ thống sập xả 10.000 lỗi,
không gửi 10.000 alert.

**Cách làm:** Track thời gian gửi alert gần nhất theo key,
nếu trong cooldown thì bỏ qua.

```go
type AlertEngine struct {
    sent     map[string]time.Time  // dedup tracking
    cooldown time.Duration
    mu       sync.RWMutex
}

func (e *AlertEngine) shouldAlert(key string) bool {
    e.mu.RLock()
    lastSent, exists := e.sent[key]
    e.mu.RUnlock()

    if exists && time.Since(lastSent) < e.cooldown {
        return false
    }
    e.mu.Lock()
    e.sent[key] = time.Now()
    e.mu.Unlock()
    return true
}
```

**Khái niệm cần biết khi bảo vệ:** Alert Fatigue là thuật ngữ chuẩn trong SRE/DevOps.

### 5.3 Dynamic Threshold
**Vấn đề giải quyết:** Thay đổi ngưỡng cảnh báo mà không cần restart server.

**Cách làm:** Dashboard gửi config mới qua WebSocket hoặc REST,
goroutine alerting đọc giá trị mới ngay lập tức.

**Khái niệm cần biết:** `sync.RWMutex` — nhiều goroutine đọc đồng thời được,
chỉ 1 goroutine ghi tại một thời điểm. Nếu không dùng mutex → race condition.

```go
// Đọc (goroutine alerting)
mu.RLock()
current := threshold
mu.RUnlock()

// Ghi (goroutine WebSocket handler)
mu.Lock()
threshold = newValue
mu.Unlock()
```

### Thứ tự implement
1. Sliding Window — dễ nhất, làm ngay từ đầu (cách đúng)
2. Alert Deduplication — tuần 3, ~2 tiếng code
3. Dynamic Threshold — tuần 3 nếu còn thời gian, ~4 tiếng code

---

## 6. Các khái niệm kỹ thuật cần hiểu

### Elasticsearch

**Inverted Index vs B-tree (MySQL)**
- MySQL dùng B-tree: tìm record theo ID/value → nhanh với exact match
- ES dùng inverted index: từ term → danh sách document chứa nó → nhanh với full-text search
- Ví dụ: search "payment timeout" trong 10 triệu log → ES trả kết quả trong ms

**Query DSL — so sánh với NoSQL**
```js
// MongoDB
db.logs.find({ level: "ERROR", app: "demo-node" })

// Elasticsearch — tương tự, chỉ khác cú pháp
{
  "query": {
    "bool": {
      "must": [
        { "term": { "level.keyword": "ERROR" } },
        { "term": { "app.keyword": "demo-node" } }
      ]
    }
  }
}
```

**term vs match**
- `term`: khớp chính xác (dùng cho level, service name)
- `match`: full-text search (dùng cho message)

**Time-range query**
```json
"filter": {
  "range": {
    "@timestamp": { "gte": "now-5m", "lte": "now" }
  }
}
```
`now-5m`, `now-1h`, `now-7d` — ES hiểu tự nhiên, không cần tính thủ công.

### Docker

**Container networking**
Các container trong cùng docker-compose network giao tiếp qua *tên service*,
không phải localhost:
```yaml
# Filebeat kết nối Logstash
output.logstash:
  hosts: ["logstash:5044"]  # ← tên service, không phải localhost
```

**depends_on + healthcheck**
ES cần 30–60 giây để sẵn sàng. Logstash phải đợi ES healthy trước khi start:
```yaml
depends_on:
  elasticsearch:
    condition: service_healthy
```

### Go Concurrency

**Goroutine**
```go
go func() {
    for range ticker.C {
        checkAlerts()  // chạy song song với main goroutine
    }
}()
```
Goroutine nhẹ ~2KB RAM, OS thread ~1MB. Go có thể chạy hàng triệu goroutine.

**sync.RWMutex**
```go
var mu sync.RWMutex

// Nhiều goroutine đọc đồng thời — không block nhau
mu.RLock()
val := sharedValue
mu.RUnlock()

// Chỉ 1 goroutine ghi — block tất cả đọc/ghi khác
mu.Lock()
sharedValue = newVal
mu.Unlock()
```

**GC (Garbage Collector)**
Cơ chế tự động dọn bộ nhớ không còn dùng. Go GC ít pause hơn JVM GC.
JVM (Spring) có "stop-the-world pause" làm chương trình dừng ngắn.

**Binary compile**
Go compile ra 1 file binary chạy được ngay, không cần runtime (JVM, Node.js).
Docker image từ `scratch` + binary → image ~15MB thay vì ~500MB với JDK.

### Logstash Grok

Grok là regex có tên, dùng để parse log text thành JSON:
```
Input:  "[2024-01-15T10:23:11Z] [ERROR] [demo-node] Payment failed"
Pattern: \[%{TIMESTAMP_ISO8601:timestamp}\] \[%{LOGLEVEL:level}\s*\] \[%{DATA:service}\] %{GREEDYDATA:message}
Output: {
  "timestamp": "2024-01-15T10:23:11Z",
  "level":     "ERROR",
  "service":   "demo-node",
  "message":   "Payment failed"
}
```
Tool test: [grokdebugger.com](https://grokdebugger.com)

---

## 7. Yêu cầu đầu ra và cách đáp ứng

### "Hệ thống hoạt động đúng với tài liệu giải pháp"
→ Viết tài liệu **song song** với code, không phải sau cùng.
Tài liệu mô tả gì thì code phải làm đúng cái đó.

### "Đảm bảo hiệu năng"
Định nghĩa thành con số đo được:
- Query 10,000 log và trả kết quả < 200ms
- Dashboard load < 2 giây
- Alerting phát hiện ERROR trong vòng 10 giây

### "Production-ready"
Không có nghĩa là deploy lên AWS với load balancer.
Với scope sinh viên cần có:

| Tiêu chí | Cách làm |
|---|---|
| Config qua biến môi trường | `.env` + `os.Getenv()` |
| Tự restart khi crash | `restart: always` trong docker-compose |
| Health check | `healthcheck` trong docker-compose cho ES |
| Không lộ thông tin nhạy cảm | ES password trong `.env`, không commit |
| Log của chính hệ thống | Logstash và Go API tự ghi log hoạt động |

### Tài liệu cần có
1. **Kiến trúc** — sơ đồ, tại sao chọn từng công nghệ
2. **Cài đặt** — clone → 3 lệnh → chạy được
3. **Sử dụng** — screenshot dashboard, giải thích tính năng
4. **API Reference** — endpoint, request/response mẫu

---

## 8. Chuẩn bị bảo vệ

### 5 câu hỏi hay gặp

**"Tại sao dùng Elasticsearch mà không dùng MySQL?"**
ES dùng inverted index, search full-text trên hàng triệu log trong millisecond.
MySQL B-tree index không tối ưu cho pattern search trên text log.
Ngoài ra ES có `now-5m` built-in cho time-range query trong alerting.

**"Filebeat và Logstash khác nhau gì, sao cần cả 2?"**
Filebeat nhẹ (~50MB), chạy trên từng server chỉ để ship data.
Logstash nặng hơn (~500MB) nhưng có processing power để parse/transform.
Tách vai trò để không đặt Logstash nặng trên mọi server.

**"Nếu Logstash die thì log có mất không?"**
Không. Filebeat có registry nhớ đã đọc đến dòng nào.
Khi Logstash phục hồi, Filebeat tiếp tục gửi từ điểm dừng, không mất log.

**"Alerting của em có false positive không?"**
Có — đó là trade-off. Ngưỡng tĩnh luôn có false positive.
Có thể tune bằng cách tăng window time hoặc tăng threshold.
Production thường dùng sliding window nhiều bước và machine learning
để giảm false positive, nhưng nằm ngoài scope dự án này.

**"Em học được gì từ project này?"**
Trả lời thật: distributed system thinking (tại sao cần pipeline tách tầng),
container networking, search engine fundamentals, Go concurrency model.

### 3 điểm giá trị để trình bày

**Giá trị 1 — Observability tập trung**
Khi debug "lỗi ở đâu" không còn mất hàng giờ SSH từng server.
Search toàn bộ log của mọi service trong 1 chỗ, trong vài giây.

**Giá trị 2 — Alerting chủ động**
Phát hiện lỗi *trước người dùng* bằng cách đếm ERROR rate real-time,
thay vì chờ người dùng báo cáo sự cố.

**Giá trị 3 — Kiến trúc pipeline tách biệt**
4 tầng độc lập: Collection → Processing → Storage → Presentation.
Có thể thay Filebeat bằng Fluentd, thay dashboard bằng Grafana
mà không ảnh hưởng tầng khác. Đây là tư duy kiến trúc thực tế.

---

## 9. Nguồn học theo thứ tự

1. [elastic.co — What is ELK Stack](https://www.elastic.co/what-is/elk-stack) — 15 phút, bức tranh tổng thể
2. [betterstack.com — Log levels explained](https://betterstack.com/community/guides/logging/log-levels-explained/) — 10 phút, nền tảng conceptual
3. [docs.docker.com — Get started](https://docs.docker.com/compose/gettingstarted/) — 1 buổi thực hành
4. [elastic.co — ES getting started](https://www.elastic.co/guide/en/elasticsearch/reference/current/getting-started.html) — 1 buổi, gõ curl query thử
5. [go.dev — RESTful API with Gin](https://go.dev/doc/tutorial/web-service-gin) — 1 buổi, bạn đã có nền Go
6. [grokdebugger.com](https://grokdebugger.com) — test Grok pattern trực tiếp

---

## 10. Chi phí

**Hoàn toàn $0** khi chạy local.

| Thành phần | Chi phí |
|---|---|
| Elasticsearch Basic License | Miễn phí |
| Logstash + Filebeat | Miễn phí |
| Go + tất cả thư viện | Miễn phí |
| Docker Desktop (WSL) | Miễn phí |
| Infrastructure local | Miễn phí |
| **Tổng** | **$0** |

Trường hợp cần deploy lên server thật: Oracle Cloud Free Tier (1 VM),
Railway.app ($5 credit miễn phí), hoặc mượn VPS từ trường.

---

## 11. Timeline 4 tuần

```
Tuần 1 — Infrastructure
  Mục tiêu: docker compose up chạy ES + Logstash + Filebeat
  Verify: gửi được 1 dòng log vào ES, query ra bằng curl
  Files: docker-compose.yml, filebeat.yml, logstash.conf

Tuần 2 — Services & API
  Mục tiêu: demo services sinh log, Go API filter được log
  Files: demo-node/index.js, demo-go/main.go, api-server/

Tuần 3 — Alerting & Dashboard
  Mục tiêu: alert banner xuất hiện khi spam ERROR
  Features: Sliding Window, Deduplication, (Dynamic Threshold nếu kịp)
  Files: alerting/engine.go, dashboard/

Tuần 4 — Docs & Polish
  Mục tiêu: chạy ổn định, tài liệu đầy đủ, chuẩn bị demo
  Files: docs/, README.md hoàn chỉnh
```

---

*Tài liệu này tổng hợp từ quá trình phân tích trước khi bắt đầu code.
Cập nhật khi có quyết định kỹ thuật mới trong quá trình phát triển.*
