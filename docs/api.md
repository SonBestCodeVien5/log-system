# API Reference

Base URL: `http://localhost:8080`

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
| `from` | ISO8601 | Thời gian bắt đầu | `2024-01-15T10:00:00Z` |
| `to` | ISO8601 | Thời gian kết thúc | `2024-01-15T11:00:00Z` |
| `q` | string | Full-text search trong message | `payment+timeout` |
| `page` | int | Số trang, bắt đầu từ 1 | `1` |
| `size` | int | Số record mỗi trang, mặc định 20 | `20` |

**Response:**
```json
{
  "data": [
    {
      "timestamp": "2024-01-15T10:23:11Z",
      "level":     "ERROR",
      "service":   "demo-node",
      "message":   "Payment gateway timeout after 30s"
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

**Query parameters:** `from`, `to`, `app` (tương tự /api/logs)

**Response:**
```json
{
  "INFO":  1240,
  "WARN":  87,
  "ERROR": 23,
  "total": 1350
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

**Response:**
```json
{ "status": "updated", "threshold": 10 }
```

---

### WebSocket /ws/alerts
Nhận alert real-time khi ERROR rate vượt ngưỡng.

**Kết nối:**
```js
const ws = new WebSocket('ws://localhost:8080/ws/alerts')

ws.onmessage = (event) => {
  const alert = JSON.parse(event.data)
  console.log(alert)
}
```

**Message format:**
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
