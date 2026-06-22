# Hướng dẫn cài đặt và triển khai

## Yêu cầu môi trường

| Thành phần | Yêu cầu tối thiểu |
|---|---|
| Docker Engine | >= 24.0 |
| Docker Compose | v2 (`docker compose`, không phải `docker-compose`) |
| RAM | 8GB trống |
| Disk | 10GB trống |
| OS | Linux, macOS, Windows + WSL2 |

---

## Cài đặt

### Bước 1 — Clone repo

```bash
git clone git@github.com:SonBestCodeVien5/log-system.git
cd log-system
```

### Bước 2 — Cấu hình

```bash
# Tùy chọn: chỉ cần tạo .env nếu muốn override default trong docker-compose.yml
cp .env.example .env
```

File `.env` là **tùy chọn** trong môi trường dev hiện tại. Nếu không tạo `.env`,
Docker Compose vẫn chạy bằng các default trong `docker-compose.yml`, ví dụ
`${ES_PASSWORD:-changeme123}` và `${API_PORT:-8080}`.

Chỉ cần tạo/chỉnh sửa `.env` nếu muốn override default:

```bash
ES_PASSWORD=changeme123        # đổi password
ALERT_THRESHOLD=10             # số ERROR trong window để trigger alert
ALERT_CHECK_INTERVAL_SECONDS=10  # chu kỳ check alerting
```

### Bước 3 — WSL: set vm.max_map_count

Bắt buộc trên WSL trước khi chạy Elasticsearch:

```bash
sudo sysctl -w vm.max_map_count=262144
```

Để persistent sau reboot:

```bash
echo "[boot]" | sudo tee -a /etc/wsl.conf
echo "command = sysctl -w vm.max_map_count=262144" | sudo tee -a /etc/wsl.conf
```

### Bước 4 — Khởi động

```bash
docker compose up -d
```

Lần đầu pull image ES + Logstash khoảng 5–10 phút tùy mạng.

### Bước 5 — Verify

```bash
# ES có sẵn sàng không?
curl -u elastic:${ES_PASSWORD:-changeme123} http://localhost:9200/_cluster/health

# Log đã vào ES chưa?
curl -u elastic:${ES_PASSWORD:-changeme123} "http://localhost:9200/logs-*/_count"

# API có chạy không?
curl http://localhost:8080/api/health
```

### Bước 6 — Mở dashboard

```
http://localhost:8080
```

---

## Kế hoạch triển khai 4 tuần

> Trạng thái hiện tại: Infrastructure, Pipeline, Demo services, Go API,
> Alerting Engine, Dashboard và incident replay đã hoàn thành và có bằng chứng
> E2E. Giai đoạn hiện tại đóng băng scope ứng dụng, tập trung clean-clone,
> screenshot, báo cáo, slide và rehearsal.
> Roadmap chi tiết xem `docs/project-roadmap.md`; roadmap học/verify/bảo vệ 1 tháng cuối xem
> [`docs/one-month-defense-roadmap.md`](one-month-defense-roadmap.md).

### Tuần 1 — Infrastructure

**Mục tiêu:** Pipeline chạy end-to-end, log vào được ES.

**Trạng thái:** Hoàn thành và đã verify E2E; output lưu tại
[`testing-evidence.md`](testing-evidence.md).

**Tasks:**
- Hoàn thiện `docker-compose.yml`
- Viết `logstash/pipeline/logstash.conf` — JSON codec + Grok enrich
- Viết `filebeat/filebeat.yml` — tail `/logs/**/*.log`
- Viết demo services sinh log JSON Lines

**Deliverable đo được:**
```bash
# Chạy lệnh này, kết quả count > 0 là thành công
curl -u elastic:${ES_PASSWORD:-changeme123} "http://localhost:9200/logs-*/_count"
# {"count": 127, ...}
```

---

### Tuần 2 — Go API Server

**Mục tiêu:** Có thể query và filter log qua REST API.

**Trạng thái:** Hoàn thành và đã verify runtime với filter level, service, count
và response shape chuẩn.

**Tasks:**
- `main.go` — khởi tạo gin, ES client, routes
- `handlers/logs.go` — GET /api/logs với filter
- `handlers/logs.go` — GET /api/logs/count
- `middleware/cors.go` — cho phép dashboard gọi API

**Deliverable đo được:**
```bash
# Filter 20 ERROR gần nhất
curl "http://localhost:8080/api/logs?level=ERROR&size=20"
# Trả về JSON đúng format {"data":[...],"total":N}

# Đo response time thực tế (ghi vào tài liệu)
time curl "http://localhost:8080/api/logs?size=100" -o /dev/null
```

---

### Tuần 3 — Alerting + Dashboard

**Mục tiêu:** Alert banner xuất hiện khi hệ thống có spike ERROR.

**Trạng thái:** Hoàn thành và đã verify. Dashboard kết nối WebSocket, incident
replay được ingest và API server ghi nhận alert sent. Screenshot banner vẫn cần
chụp để dùng trong slide.

**Tasks:**
- `alerting/engine.go` — Sliding Window + Deduplication
- `handlers/alerts.go` — WebSocket /ws/alerts
- `handlers/alerts.go` — POST /api/alerts/config (Dynamic Threshold)
- `dashboard/index.html + app.js` — bảng log, filter, alert banner

**Deliverable đo được:**

Kịch bản test: tăng tốc sinh ERROR trong demo service → đo thời gian đến khi banner xuất hiện.

```text
Kết quả 2026-06-17: 20 incident ERROR được ingest; alert sent quan sát sau ~33s.
Nguyên nhân lớn hơn check interval: alert trước đó vẫn còn trong cooldown 60s.
```

---

### Tuần 4 — Polish + Bảo vệ

**Mục tiêu:** Hệ thống chạy ổn định, tài liệu đầy đủ, demo được.

**Trạng thái:** Đang đóng gói. Số liệu thật và demo script đã có; còn test clone
sạch, screenshot, slide và rehearsal trước khi bảo vệ.

**Tasks:**
- Đo và ghi nhận số liệu hiệu năng thực tế vào tài liệu
- Hoàn thiện README — clone → 3 lệnh → chạy được
- Chuẩn bị demo script cho buổi bảo vệ
- Chuẩn bị câu trả lời cho 5 câu hỏi hay gặp

**Deliverable đo được:**
```bash
# Test từ repo mới clone — không được quá 5 phút
git clone git@github.com:SonBestCodeVien5/log-system.git fresh-test
cd fresh-test
# Tùy chọn nếu muốn override default:
# cp .env.example .env
sudo sysctl -w vm.max_map_count=262144
docker compose up -d
# Mở browser → http://localhost:8080 → thấy dashboard
```

---

## Số liệu hiệu năng

> **Lưu ý:** Các con số dưới đây được đo bằng request `curl` đơn lẻ trên môi
> trường dev ngày 2026-06-17. Đây không phải benchmark tải nặng hoặc SLA.

| Metric | Cơ chế đảm bảo | Kết quả đo thực tế |
|---|---|---|
| Query response time | ES inverted index + index theo ngày | 2026-06-17: `/api/logs?size=100` trả HTTP 200 trong `0.033992s`; `/api/health` trả HTTP 200 trong `0.007250s` |
| Dashboard load time | Pagination 20 record/trang | 2026-06-17: dashboard `/` trả HTTP 200 trong `0.012889s` |
| Alert detection latency | Check interval = `ALERT_CHECK_INTERVAL_SECONDS` | 2026-06-17: incident replay được ingest (`total: 20`); alert sent quan sát sau ~33s vì cooldown trước đó còn active |
| Log durability | Filebeat registry | Có cơ chế resume offset; chưa có phép thử restart định lượng riêng |

### Lệnh tái đo khi cần cập nhật báo cáo

```bash
time curl -s "http://localhost:8080/api/logs?size=100" -o /dev/null
time curl -s "http://localhost:8080/api/health" -o /dev/null
```

Khi tái đo trên máy khác, ghi rõ ngày, cấu hình môi trường và lưu output chi tiết
trong `docs/testing-evidence.md`.

## Việc triển khai còn mở trước bảo vệ

- Chạy flow từ một clean clone độc lập và ghi tổng thời gian đến khi dashboard có data.
- Chụp dashboard bình thường, trạng thái filter ERROR và alert banner.
- Rehearse incident replay sau khi cooldown cũ đã hết để demo không bị trễ bất ngờ.
- CORS và WebSocket origin hiện mở cho local demo; không expose stack trực tiếp ra Internet.

---

## Lệnh thường dùng

```bash
# Xem trạng thái
docker compose ps

# Xem log từng service
docker compose logs -f elasticsearch
docker compose logs -f logstash
docker compose logs -f filebeat

# Restart 1 service
docker compose restart logstash

# Dừng hệ thống
docker compose down

# Dừng và xóa data (mất hết log)
docker compose down -v
```

---

## Xử lý sự cố

### Log không vào Elasticsearch

```bash
# 1. ES có chạy không?
curl -u elastic:${ES_PASSWORD:-changeme123} http://localhost:9200/_cluster/health

# 2. Logstash nhận data không?
docker compose logs logstash | grep -E "events|error"

# 3. Filebeat đang tail đúng không?
docker compose logs filebeat | grep -E "Harvester|error"

# 4. File log có tồn tại không?
ls -la ./logs/demo-node/ ./logs/demo-go/

# 5. JSON parse lỗi?
docker compose logs logstash | grep "json parse"

# 6. ES có index chưa?
curl -u elastic:${ES_PASSWORD:-changeme123} "http://localhost:9200/_cat/indices?v"
```

### Elasticsearch không start

```bash
# Lỗi phổ biến nhất trên WSL
sudo sysctl -w vm.max_map_count=262144
docker compose restart elasticsearch
```

### API không kết nối được ES

Kiểm tra `ES_PASSWORD` trong `.env` hoặc default `changeme123` có khớp với lúc khởi tạo ES không.
Nếu đã đổi password sau khi ES đã chạy:

```bash
docker compose down -v  # xóa data cũ
docker compose up -d    # khởi động lại
```

---

## Chi phí

| Thành phần | Chi phí |
|---|---|
| Elasticsearch Basic License | Miễn phí |
| Logstash + Filebeat | Miễn phí |
| Go + tất cả thư viện | Miễn phí |
| Docker | Miễn phí |
| **Tổng** | **$0** |
