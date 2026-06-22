# API Reference

Base URL: `http://localhost:8080`

Roadmap verify API trong 1 tháng cuối nằm ở
[`docs/one-month-defense-roadmap.md`](one-month-defense-roadmap.md).

## Endpoints

### GET /api/health
Health check.

**Response:**
```json
{ "status": "ok", "elasticsearch": "connected" }
```

---

### GET /api/logs
Lấy danh sách log, hỗ trợ filter và phân trang.

**Query parameters:**

| Param | Type | Mô tả | Ví dụ |
|---|---|---|---|
| `level` | string | Lọc theo level | `ERROR`, `WARN`, `INFO` |
| `app` | string | Lọc theo service | `demo-node` |
| `from` | date math hoặc ISO8601 | Thời gian bắt đầu, mặc định `now-1h` | `now-5m` |
| `to` | date math hoặc ISO8601 | Thời gian kết thúc, mặc định `now` | `now` |
| `q` | string | Full-text search trong message | `payment+timeout` |
| `page` | int | Số trang, bắt đầu từ 1; giá trị nhỏ hơn 1 được đưa về 1 | `1` |
| `size` | int | Số record mỗi trang, mặc định 20, tối đa 100 | `20` |

Kết quả được sort theo `@timestamp desc`. `level` và `app` dùng exact match;
`q` dùng full-text match trên `log_message`. `size` ngoài khoảng 1–100 được đưa
về mặc định 20.

**Response:**
```json
{
  "data": [
    {
      "@timestamp":  "2024-01-15T10:23:11Z",
      "level":       "ERROR",
      "service":     "demo-node",
      "log_message": "Payment gateway timeout after 30s",
      "metadata":    { "order_id": "789" }
    }
  ],
  "total": 100,
  "page":  1,
  "size":  20
}
```

**Ví dụ:**
```bash
# Lấy 20 ERROR gần nhất
curl "http://localhost:8080/api/logs?level=ERROR&size=20"

# Tìm log chứa "payment" của demo-node
curl "http://localhost:8080/api/logs?app=demo-node&q=payment"

# Lọc theo khoảng thời gian
curl "http://localhost:8080/api/logs?from=2024-01-15T10:00:00Z&to=2024-01-15T11:00:00Z"
```

---

### GET /api/logs/count
Đếm log theo level trong khoảng thời gian.

**Query parameters:** `from`, `to`, `app` (tương tự `/api/logs`). Mặc định đếm
trong một giờ gần nhất.

**Response:**
```json
{
  "INFO":  1240,
  "WARN":  87,
  "ERROR": 23,
  "total": 1350,
  "from": "now-1h",
  "to": "now"
}
```

---

### POST /api/alerts/config
Cập nhật cấu hình alerting động, không cần restart server.

**Request body:**
```json
{
  "threshold":       10,
  "window_seconds":  300,
  "cooldown_seconds": 60
}
```

Có thể gửi một hoặc nhiều field. Field bị bỏ qua sẽ giữ giá trị hiện tại:

```json
{ "threshold": 5 }
```

Các field nếu được gửi phải là số nguyên `>= 1`; body rỗng `{}` bị từ chối.

**Response:**
```json
{
  "status": "updated",
  "config": {
    "threshold": 10,
    "window_seconds": 300,
    "cooldown_seconds": 60
  }
}
```

---

### WebSocket /ws/alerts
Nhận alert real-time khi ERROR rate vượt ngưỡng.

Endpoint hiện cho phép mọi origin để phục vụ môi trường demo local. Khi triển
khai public phải thay bằng allowlist origin và cấu hình CORS tương ứng.

**Kết nối:**
```js
const ws = new WebSocket('ws://localhost:8080/ws/alerts')

ws.onmessage = (event) => {
  const alert = JSON.parse(event.data)
  console.log(alert)
}
```

**Message format:**

Khi vừa connect, server gửi config hiện tại:

```json
{
  "type": "config",
  "config": {
    "threshold": 10,
    "window_seconds": 300,
    "cooldown_seconds": 60
  }
}
```

Khi có spike ERROR, server gửi alert:

```json
{
  "type":      "error_spike",
  "count":     25,
  "threshold": 10,
  "window":    "5m",
  "timestamp": "2024-01-15T10:23:11Z",
  "message":   "25 errors in last 5 minutes (threshold: 10)"
}
```

## Error responses

Tất cả lỗi trả về cùng format:

```json
{ "error": "mô tả lỗi" }
```

| HTTP Status | Ý nghĩa |
|---|---|
| 200 | Thành công |
| 400 | Request không hợp lệ |
| 500 | Lỗi server hoặc Elasticsearch |
| 503 | Elasticsearch chưa kết nối được ở `/api/health` |

## Trạng thái xác minh MVP

Các contract chính đã được kiểm tra E2E ngày 2026-06-17: filter level, filter
service, count theo level, partial alert config và incident replay. Lần kiểm tra
dashboard ngày 2026-06-22 tiếp tục ghi nhận các request logs/count trả HTTP 200.
Chi tiết nằm trong [`testing-evidence.md`](testing-evidence.md).
