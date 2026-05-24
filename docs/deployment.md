# Hướng dẫn cài đặt và triển khai

## Yêu cầu môi trường

| Thành phần | Yêu cầu |
|---|---|
| Docker Engine | >= 24.0 |
| Docker Compose | v2 (lệnh `docker compose`, không phải `docker-compose`) |
| RAM | Tối thiểu 8GB trống, khuyến nghị 12GB |
| Disk | 10GB trống |
| OS | Linux, macOS, Windows + WSL2 |

## Cài đặt

### Bước 1 — Clone repo

```bash
git clone git@github.com:SonBestCodeVien5/log-system.git
cd log-system
```

### Bước 2 — Tạo file config

```bash
cp .env.example .env
```

Chỉnh sửa `.env` nếu cần — đặc biệt đổi `ES_PASSWORD`:

```bash
# .env
ES_VERSION=8.13.0
ES_PORT=9200
ES_PASSWORD=changeme123        # ← đổi password này

LOGSTASH_PORT=5044
API_PORT=8080

ALERT_THRESHOLD=10
ALERT_WINDOW_SECONDS=300
ALERT_COOLDOWN_SECONDS=60
ALERT_CHECK_INTERVAL_SECONDS=10
```

### Bước 3 — Khởi động

```bash
docker compose up -d
```

Lần đầu sẽ pull image ES + Logstash (~1.5GB), mất 5–10 phút tùy tốc độ mạng.

### Bước 4 — Kiểm tra hệ thống

```bash
# ES có sẵn sàng không?
curl http://localhost:9200/_cluster/health

# Kết quả mong đợi:
# {"status":"green"} hoặc {"status":"yellow"}

# Log đã vào ES chưa?
curl "http://localhost:9200/logs-*/_count"

# API server có chạy không?
curl http://localhost:8080/api/health
```

### Bước 5 — Mở dashboard

Truy cập `http://localhost:8080` trên trình duyệt.

## Các lệnh thường dùng

```bash
# Xem trạng thái tất cả service
docker compose ps

# Xem log của từng service
docker compose logs -f elasticsearch
docker compose logs -f logstash
docker compose logs -f filebeat
docker compose logs -f api-server

# Dừng toàn bộ
docker compose down

# Dừng và xóa data (cẩn thận — mất hết log trong ES)
docker compose down -v

# Restart 1 service cụ thể
docker compose restart logstash
```

## Xử lý sự cố thường gặp

### Log không vào Elasticsearch

Kiểm tra theo thứ tự:

```bash
# 1. ES có chạy không?
curl http://localhost:9200/_cluster/health

# 2. Logstash có nhận data từ Filebeat không?
docker compose logs logstash | grep "events"

# 3. Filebeat có đang tail đúng file không?
docker compose logs filebeat | grep "Harvester"

# 4. File log có tồn tại không?
ls -la ./logs/demo-node/
ls -la ./logs/demo-go/

# 5. Grok có parse lỗi không?
docker compose logs logstash | grep "_grokparsefailure"
```

### Elasticsearch không khởi động được

Thường do thiếu RAM hoặc vm.max_map_count quá thấp:

```bash
# Chạy lệnh này trên Linux/WSL
sudo sysctl -w vm.max_map_count=262144

# Để persistent sau reboot, thêm vào /etc/sysctl.conf
echo "vm.max_map_count=262144" | sudo tee -a /etc/sysctl.conf
```

### API server không kết nối được ES

Kiểm tra biến môi trường trong `.env` — `ES_PASSWORD` phải khớp với lúc khởi tạo ES.
Nếu đã thay đổi password sau khi ES đã chạy thì cần `docker compose down -v` và khởi động lại.

## Cấu trúc dữ liệu trong Elasticsearch

```bash
# Xem danh sách index
curl "http://localhost:9200/_cat/indices?v"

# Index đặt tên theo ngày: logs-YYYY.MM.DD
# Ví dụ: logs-2024.01.15

# Xem mapping của index
curl "http://localhost:9200/logs-*/_mapping"

# Query thử 5 log ERROR gần nhất
curl -X GET "http://localhost:9200/logs-*/_search" \
  -H "Content-Type: application/json" \
  -d '{
    "size": 5,
    "sort": [{"@timestamp": "desc"}],
    "query": {
      "term": { "level.keyword": "ERROR" }
    }
  }'
```

## Chi phí

| Thành phần | Chi phí |
|---|---|
| Elasticsearch Basic License | Miễn phí |
| Logstash + Filebeat | Miễn phí |
| Go + tất cả thư viện | Miễn phí |
| Docker | Miễn phí |
| Infrastructure (local) | Miễn phí |
| **Tổng** | **$0** |
