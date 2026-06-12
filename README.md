# Log System - Centralized Logging Platform

Hệ thống log tập trung dùng để thu thập log từ nhiều service, parse/enrich qua
Logstash, lưu trữ và tìm kiếm bằng Elasticsearch, cung cấp REST API bằng Go và
hiển thị dashboard realtime bằng HTML/CSS/JavaScript thuần.

## Mục lục

- [Kiến trúc](#kien-truc)
- [Yêu cầu môi trường](#yeu-cau-moi-truong)
- [Cài đặt nhanh](#cai-dat-nhanh)
- [Cấu hình môi trường](#cau-hinh-moi-truong)
- [Chạy hệ thống bằng Docker Compose](#chay-he-thong-bang-docker-compose)
- [Kiểm tra sau khi chạy](#kiem-tra-sau-khi-chay)
- [Truy cập hệ thống](#truy-cap-he-thong)
- [API nhanh](#api-nhanh)
- [Vận hành thường ngày](#van-hanh-thuong-ngay)
- [Chạy API server ở chế độ dev](#chay-api-server-o-che-do-dev)
- [Triển khai trên máy mới hoặc server](#trien-khai-tren-may-moi-hoac-server)
- [Xử lý sự cố](#xu-ly-su-co)
- [Tài liệu liên quan](#tai-lieu-lien-quan)

<a id="kien-truc"></a>

## Kiến trúc

```text
Demo Services
  -> ghi JSON Lines vào ./logs/<service>/app.log
  -> Filebeat tail log file
  -> Logstash parse JSON + enrich message
  -> Elasticsearch lưu index logs-YYYY.MM.dd
  -> Go API Server query Elasticsearch
  -> Dashboard gọi REST API và nhận alert qua WebSocket
```

Thành phần chính:

| Thành phần | Vai trò | Port mặc định |
|---|---|---|
| `demo-node` | Service Node.js sinh log demo | Nội bộ Docker |
| `demo-go` | Service Go sinh log demo | Nội bộ Docker |
| `filebeat` | Tail file trong `./logs` và gửi tới Logstash | Nội bộ Docker |
| `logstash` | Parse JSON Lines, enrich field, đẩy vào Elasticsearch | `5044` |
| `elasticsearch` | Lưu trữ và tìm kiếm log | `9200` |
| `api-server` | REST API, WebSocket alert, serve dashboard | `8080` |
| `dashboard` | Giao diện HTML/JS tĩnh mount vào API container | Qua `api-server` |

<a id="yeu-cau-moi-truong"></a>

## Yêu cầu môi trường

| Thành phần | Khuyến nghị |
|---|---|
| Docker Engine | 24.x trở lên |
| Docker Compose | v2, dùng lệnh `docker compose` |
| RAM trống | Tối thiểu 8GB, khuyến nghị 12GB |
| Disk trống | Tối thiểu 10GB |
| Hệ điều hành | Linux, macOS, hoặc Windows + WSL2 |

Trên Windows nên chạy trong WSL2 để Elasticsearch và volume Docker ổn định hơn.

<a id="cai-dat-nhanh"></a>

## Cài đặt nhanh

```bash
git clone git@github.com:SonBestCodeVien5/log-system.git
cd log-system

cp .env.example .env
# Mở .env và đổi ES_PASSWORD trước khi chạy thật.

# Linux/WSL2 cần cấu hình này trước khi Elasticsearch start.
sudo sysctl -w vm.max_map_count=262144

docker compose up -d --build
```

Sau khi các container healthy, mở dashboard:

```text
http://localhost:8080
```

<a id="cau-hinh-moi-truong"></a>

## Cấu hình môi trường

File `.env.example` là mẫu cấu hình. Khi triển khai, tạo `.env` riêng:

```bash
cp .env.example .env
```

Các biến quan trọng:

| Biến | Mặc định trong compose | Ý nghĩa |
|---|---:|---|
| `ES_VERSION` | `8.13.0` | Version Elasticsearch, Logstash, Filebeat |
| `ES_PORT` | `9200` | Port Elasticsearch expose ra host |
| `ES_PASSWORD` | `changeme123` nếu không có `.env` | Password user `elastic` |
| `LOGSTASH_PORT` | `5044` | Port Beats input của Logstash |
| `API_PORT` | `8080` | Port REST API và dashboard |
| `ALERT_THRESHOLD` | `10` | Số ERROR tối thiểu để phát alert |
| `ALERT_WINDOW_SECONDS` | `300` | Cửa sổ thời gian đếm ERROR |
| `ALERT_COOLDOWN_SECONDS` | `60` | Thời gian chống gửi trùng alert |
| `ALERT_CHECK_INTERVAL_SECONDS` | `10` | Chu kỳ engine kiểm tra alert |

Lưu ý bảo mật:

- Không commit file `.env`.
- Đổi `ES_PASSWORD` trước khi triển khai thật.
- Nếu đã chạy Elasticsearch rồi mới đổi `ES_PASSWORD`, volume cũ vẫn giữ mật
  khẩu ban đầu. Với môi trường dev có thể reset bằng `docker compose down -v`
  rồi chạy lại.

<a id="chay-he-thong-bang-docker-compose"></a>

## Chạy hệ thống bằng Docker Compose

### 1. Chuẩn bị thư mục log

Repo đã có cấu trúc `logs/demo-node` và `logs/demo-go`. Nếu clone mới chưa có
do `.gitignore`, tạo lại:

```bash
mkdir -p logs/demo-node logs/demo-go
```

### 2. Cấu hình WSL/Linux cho Elasticsearch

Elasticsearch cần `vm.max_map_count` đủ lớn. Trên Linux hoặc WSL2, chạy:

```bash
sudo sysctl -w vm.max_map_count=262144
```

Trên WSL2, nếu muốn giữ sau khi restart WSL, thêm cấu hình vào `/etc/wsl.conf`
theo hướng dẫn trong `docs/deployment.md`.

### 3. Build và start

```bash
docker compose up -d --build
```

Lần chạy đầu có thể mất vài phút vì Docker cần pull image Elastic Stack.

### 4. Theo dõi trạng thái

```bash
docker compose ps
docker compose logs -f elasticsearch
docker compose logs -f logstash
docker compose logs -f filebeat
docker compose logs -f api-server
```

<a id="kiem-tra-sau-khi-chay"></a>

## Kiểm tra sau khi chạy

Nạp biến từ `.env` để dùng lại trong lệnh kiểm tra:

```bash
set -a
source .env
set +a
```

Kiểm tra Elasticsearch:

```bash
curl -u "elastic:${ES_PASSWORD}" \
  "http://localhost:${ES_PORT:-9200}/_cluster/health?pretty"
```

Kết quả hợp lệ có `status` là `green` hoặc `yellow`.

Kiểm tra index log:

```bash
curl -u "elastic:${ES_PASSWORD}" \
  "http://localhost:${ES_PORT:-9200}/logs-*/_count?pretty"
```

Nếu demo services đã chạy được một lúc, `count` nên lớn hơn `0`.

Kiểm tra API:

```bash
curl "http://localhost:${API_PORT:-8080}/api/health"
```

Kết quả kỳ vọng:

```json
{"elasticsearch":"connected","status":"ok"}
```

Kiểm tra query log qua API:

```bash
curl "http://localhost:${API_PORT:-8080}/api/logs?size=20"
curl "http://localhost:${API_PORT:-8080}/api/logs?level=ERROR&size=20"
curl "http://localhost:${API_PORT:-8080}/api/logs/count?from=now-1h"
```

<a id="truy-cap-he-thong"></a>

## Truy cập hệ thống

| Mục | URL |
|---|---|
| Dashboard | `http://localhost:8080` |
| API health | `http://localhost:8080/api/health` |
| API logs | `http://localhost:8080/api/logs` |
| API count | `http://localhost:8080/api/logs/count` |
| WebSocket alerts | `ws://localhost:8080/ws/alerts` |
| Elasticsearch | `http://localhost:9200` |

Elasticsearch yêu cầu Basic Auth:

```bash
curl -u "elastic:${ES_PASSWORD}" "http://localhost:9200"
```

<a id="api-nhanh"></a>

## API nhanh

Base URL mặc định:

```text
http://localhost:8080
```

| Method | Path | Mô tả |
|---|---|---|
| `GET` | `/api/health` | Kiểm tra API và kết nối Elasticsearch |
| `GET` | `/api/logs` | Lấy danh sách log, hỗ trợ filter và phân trang |
| `GET` | `/api/logs/count` | Đếm log theo level |
| `POST` | `/api/alerts/config` | Cập nhật threshold alert runtime |
| `GET` | `/ws/alerts` | WebSocket nhận alert realtime |

Ví dụ filter log:

```bash
curl "http://localhost:8080/api/logs?level=ERROR&app=demo-node&page=1&size=20"
curl "http://localhost:8080/api/logs?q=payment&from=now-1h&to=now"
```

Cập nhật cấu hình alert:

```bash
curl -X POST "http://localhost:8080/api/alerts/config" \
  -H "Content-Type: application/json" \
  -d '{"threshold":5,"window_seconds":300,"cooldown_seconds":60}'
```

Chi tiết contract API nằm trong `docs/api.md`.

<a id="van-hanh-thuong-ngay"></a>

## Vận hành thường ngày

Xem trạng thái:

```bash
docker compose ps
```

Xem log toàn hệ thống:

```bash
docker compose logs -f
```

Xem log một service:

```bash
docker compose logs -f api-server
docker compose logs -f demo-node
docker compose logs -f demo-go
```

Restart một service:

```bash
docker compose restart api-server
docker compose restart logstash
```

Dừng hệ thống nhưng giữ dữ liệu Elasticsearch và registry Filebeat:

```bash
docker compose down
```

Dừng và xóa toàn bộ dữ liệu local:

```bash
docker compose down -v
```

Lệnh trên sẽ xóa index Elasticsearch và vị trí đọc của Filebeat, chỉ dùng khi
muốn reset môi trường dev.

<a id="chay-api-server-o-che-do-dev"></a>

## Chạy API server ở chế độ dev

Khi cần sửa Go API nhanh, vẫn nên chạy Elasticsearch bằng Docker Compose:

```bash
docker compose up -d elasticsearch logstash filebeat demo-node demo-go
```

Sau đó chạy API local:

```bash
cd api-server
set -a
source ../.env
set +a
ES_HOST=localhost go run main.go
```

Chế độ này phù hợp để test REST API nhanh. Route dashboard `/` trong binary Go
cần thư mục static như khi chạy container, nên khi cần test dashboard đầy đủ hãy
chạy `api-server` bằng Docker Compose.

<a id="trien-khai-tren-may-moi-hoac-server"></a>

## Triển khai trên máy mới hoặc server

Checklist triển khai:

1. Cài Docker Engine và Docker Compose v2.
2. Clone repo vào server.
3. Tạo file `.env` từ `.env.example`.
4. Đổi `ES_PASSWORD` thành mật khẩu riêng.
5. Kiểm tra port `9200`, `5044`, `8080` chưa bị chiếm hoặc đổi port trong `.env`.
6. Chạy `sudo sysctl -w vm.max_map_count=262144` nếu server là Linux/WSL2.
7. Chạy `docker compose up -d --build`.
8. Verify Elasticsearch, API health và `logs-*/_count`.
9. Mở firewall/reverse proxy cho `API_PORT` nếu cần truy cập từ máy khác.

Khuyến nghị khi triển khai ngoài máy dev:

- Chỉ public dashboard/API qua reverse proxy nội bộ hoặc VPN.
- Không public Elasticsearch trực tiếp ra Internet.
- Backup Docker volume `es_data` nếu cần giữ log lâu dài.
- Theo dõi dung lượng disk vì Elasticsearch sẽ tăng theo lượng log.

<a id="xu-ly-su-co"></a>

## Xử lý sự cố

### Elasticsearch báo lỗi `vm.max_map_count`

Chạy:

```bash
sudo sysctl -w vm.max_map_count=262144
docker compose restart elasticsearch
```

### Curl Elasticsearch bị `401 Unauthorized`

Elasticsearch đang bật security. Dùng đúng user và password:

```bash
source .env
curl -u "elastic:${ES_PASSWORD}" "http://localhost:${ES_PORT:-9200}/_cluster/health?pretty"
```

Nếu đã đổi password sau khi volume Elasticsearch được tạo, reset môi trường dev:

```bash
docker compose down -v
docker compose up -d --build
```

### API health không `ok`

Kiểm tra API và Elasticsearch:

```bash
docker compose ps
docker compose logs -f api-server
docker compose logs -f elasticsearch
```

Đảm bảo `ES_PASSWORD` trong `.env` khớp với password Elasticsearch đang dùng.

### Log không vào Elasticsearch

Kiểm tra theo thứ tự pipeline:

```bash
ls -la logs/demo-node logs/demo-go
docker compose logs -f demo-node
docker compose logs -f demo-go
docker compose logs -f filebeat
docker compose logs -f logstash
source .env
curl -u "elastic:${ES_PASSWORD}" "http://localhost:${ES_PORT:-9200}/_cat/indices?v"
```

Nếu file log có dữ liệu nhưng Elasticsearch không có index `logs-*`, kiểm tra
Logstash và Filebeat trước. Nếu JSON parse lỗi, xem log của `logstash`.

### Dashboard không hiển thị dữ liệu

Kiểm tra API trực tiếp:

```bash
curl "http://localhost:8080/api/logs?size=5"
curl "http://localhost:8080/api/logs/count?from=now-1h"
```

Nếu API có dữ liệu nhưng dashboard trống, mở DevTools của browser để xem lỗi
JavaScript hoặc lỗi network tới `http://localhost:8080`.

### Port đã bị chiếm

Đổi port trong `.env`, ví dụ:

```env
ES_PORT=9201
API_PORT=18080
LOGSTASH_PORT=15044
```

Sau đó chạy lại:

```bash
docker compose up -d
```

<a id="tai-lieu-lien-quan"></a>

## Tài liệu liên quan

- [`docs/architecture.md`](docs/architecture.md) - kiến trúc và luồng dữ liệu
- [`docs/one-month-defense-roadmap.md`](docs/one-month-defense-roadmap.md) - roadmap 1 tháng cuối: verify, học sâu, incident demo, bảo vệ
- [`docs/api.md`](docs/api.md) - contract API và ví dụ request/response
- [`docs/deployment.md`](docs/deployment.md) - hướng dẫn cài đặt, triển khai, vận hành chi tiết
- [`docs/testing-evidence.md`](docs/testing-evidence.md) - lệnh kiểm thử và bằng chứng xác minh
- [`docs/decisions.md`](docs/decisions.md) - quyết định kỹ thuật và lý do lựa chọn
