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
cp .env.example .env
```

Chỉnh sửa `.env` nếu cần:

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
curl http://localhost:9200/_cluster/health

# Log đã vào ES chưa?
curl "http://localhost:9200/logs-*/_count"

# API có chạy không?
curl http://localhost:8080/api/health
```

### Bước 6 — Mở dashboard

```
http://localhost:8080
```

---

## Kế hoạch triển khai 4 tuần

> Trạng thái hiện tại: các phần Infrastructure, Pipeline, Demo services, Go API,
> Alerting Engine và Dashboard đã có trong repo. Giai đoạn tiếp theo không thêm
> scope lớn; tập trung chạy Bước 10 end-to-end test, điền số liệu thật vào docs,
> chuẩn bị bảo vệ và chỉ thêm kịch bản incident nhỏ nếu cần demo alert chủ động.
> Roadmap chi tiết xem `docs/project-roadmap.md`; roadmap học/verify/bảo vệ 1 tháng cuối xem
> [`docs/one-month-defense-roadmap.md`](one-month-defense-roadmap.md).

### Tuần 1 — Infrastructure

**Mục tiêu:** Pipeline chạy end-to-end, log vào được ES.

**Trạng thái:** Hoàn thành, cần ghi lại output verify mới nhất sau Bước 10.

**Tasks:**
- Hoàn thiện `docker-compose.yml`
- Viết `logstash/pipeline/logstash.conf` — JSON codec + Grok enrich
- Viết `filebeat/filebeat.yml` — tail `/logs/**/*.log`
- Viết demo services sinh log JSON Lines

**Deliverable đo được:**
```bash
# Chạy lệnh này, kết quả count > 0 là thành công
curl "http://localhost:9200/logs-*/_count"
# {"count": 127, ...}
```

---

### Tuần 2 — Go API Server

**Mục tiêu:** Có thể query và filter log qua REST API.

**Trạng thái:** Hoàn thành, cần verify runtime bằng các lệnh filter API trong
`docs/project-roadmap.md`.

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

**Trạng thái:** Hoàn thành về code, cần rebuild/start `api-server` và verify
dashboard + WebSocket alert trong Bước 10.1-10.3.

**Tasks:**
- `alerting/engine.go` — Sliding Window + Deduplication
- `handlers/alerts.go` — WebSocket /ws/alerts
- `handlers/alerts.go` — POST /api/alerts/config (Dynamic Threshold)
- `dashboard/index.html + app.js` — bảng log, filter, alert banner

**Deliverable đo được:**

Kịch bản test: tăng tốc sinh ERROR trong demo service → đo thời gian đến khi banner xuất hiện.

```
Kết quả kỳ vọng: banner xuất hiện trong khoảng ALERT_CHECK_INTERVAL_SECONDS giây
Kết quả thực tế: đo và ghi vào tài liệu sau khi hoàn thành
```

---

### Tuần 4 — Polish + Bảo vệ

**Mục tiêu:** Hệ thống chạy ổn định, tài liệu đầy đủ, demo được.

**Trạng thái:** Đang làm. Cần điền số liệu thật, chuẩn bị demo script 5 phút và
test clone sạch trước khi bảo vệ.

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
cp .env.example .env
sudo sysctl -w vm.max_map_count=262144
docker compose up -d
# Mở browser → http://localhost:8080 → thấy dashboard
```

---

## Số liệu hiệu năng

> **Lưu ý:** Các con số dưới đây sẽ được cập nhật sau khi đo thực tế
> trên môi trường phát triển (WSL2, 16GB RAM, SSD).
> Kết quả có thể thay đổi tùy cấu hình môi trường.

| Metric | Cơ chế đảm bảo | Kết quả đo thực tế |
|---|---|---|
| Query response time | ES inverted index + index theo ngày | Đo sau tuần 2 |
| Dashboard load time | Pagination 20 record/trang | Đo sau tuần 3 |
| Alert detection latency | Check interval = `ALERT_CHECK_INTERVAL_SECONDS` | Tối đa = interval value |
| Log durability | Filebeat registry | Không mất log khi restart |

### Lệnh đo cần chạy ở Bước 10

```bash
time curl -s "http://localhost:8080/api/logs?size=100" -o /dev/null
time curl -s "http://localhost:8080/api/health" -o /dev/null
```

Sau khi có kết quả, ghi số `real` vào bảng trên và lưu output chi tiết trong
`docs/testing-evidence.md`.

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
curl http://localhost:9200/_cluster/health

# 2. Logstash nhận data không?
docker compose logs logstash | grep -E "events|error"

# 3. Filebeat đang tail đúng không?
docker compose logs filebeat | grep -E "Harvester|error"

# 4. File log có tồn tại không?
ls -la ./logs/demo-node/ ./logs/demo-go/

# 5. JSON parse lỗi?
docker compose logs logstash | grep "json parse"

# 6. ES có index chưa?
curl "http://localhost:9200/_cat/indices?v"
```

### Elasticsearch không start

```bash
# Lỗi phổ biến nhất trên WSL
sudo sysctl -w vm.max_map_count=262144
docker compose restart elasticsearch
```

### API không kết nối được ES

Kiểm tra `ES_PASSWORD` trong `.env` khớp với lúc khởi tạo ES.
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
