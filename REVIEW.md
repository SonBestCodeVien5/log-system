# Code Review — log-system (nhánh `claude/code-review-fixes`)

Review thực hiện trên `main` tại commit `497263e`. 7 commit, mỗi commit
tập trung vào 1 vấn đề.

## 🔴 Bug nghiêm trọng đã sửa

### 1. Resource leak trong `CountLogs` (`api-server/handlers/logs.go`)

`defer res.Body.Close()` đặt trong vòng `for` 3 level. `defer` chạy
khi function return, không phải mỗi iteration → 3 response body bị
giữ tới cuối handler. Khi `/api/logs/count` được gọi liên tục (dashboard
auto-refresh mỗi 10s), connection pool tới ES tích lũy.

**Fix**: tách thành hàm `countOneLevel` để defer chạy đúng phạm vi.
Tiện thêm `context.WithTimeout(5s)` cho cả `Search` lẫn `Count`, bỏ
hack `var _ = fmt.Sprintf / time.Now / io.Discard` cuối file.

### 2. api-server kết nối sai port khi đổi `ES_PORT` (`docker-compose.yml`)

`ES_PORT` trong `.env` được hiểu là **host-binding port** (phần trái
của `${ES_PORT}:9200`). Trước đây mình truyền y nguyên vào env của
api-server, nhưng api-server kết nối qua tên service `elasticsearch`
trong Docker network — port nội bộ **luôn là 9200**. Hậu quả: đặt
`ES_PORT=9201` trong `.env` (đúng như README hướng dẫn khi conflict
port trên host) → api-server fail kết nối ES.

**Fix**: hard-code `ES_PORT=9200` trong env của api-server. Phần
`ports` mapping vẫn dùng `${ES_PORT}` để giữ khả năng đổi host port.

### 3. Race + dead-conn buildup trong alerting (`api-server/alerting/engine.go`)

Ba thứ:

- `log.Printf("total=%d", len(e.clients))` đọc map **sau khi đã
  Unlock** → data race với register/unregister đồng thời.
- `broadcast()` giữ `RLock` trong khi gọi `WriteMessage` (blocking I/O
  dưới read lock), và nếu write fail thì không xóa conn → conn chết
  ở lại đến khi read goroutine bên handler nhận EOF. Trong khoảng đó
  mọi alert đều cố write conn chết → log spam, có thể block.
- `var fmt_placeholder = fmt.Sprintf` ở cuối file là dead code.

**Fix**: snapshot `len` trong critical section; snapshot conn list
dưới `RLock`, write ngoài lock, gom dead conn rồi xóa + `Close()`
trong một critical section ngắn; xóa dead code.

## 🟡 Best practice / robustness

### 4. WebSocket route tạo handler mới mỗi request (`main.go`)

`r.GET("/ws/alerts", func(c) { handlers.NewAlertHandler(engine).HandleWS(c) })`
tạo handler mới cho mỗi connection thay vì reuse instance đã có ở
route group `/api`. Đã dùng chung.

### 5. Graceful shutdown cho alerting engine (`main.go` + `engine.go`)

`engine.Run()` `for range ticker.C` mãi mãi, không nhận signal dừng.
Sửa: `Run(ctx)` với `select { ctx.Done() | ticker.C }`. `main` tạo
`engineCtx/engineCancel`, gọi `engineCancel()` trước `srv.Shutdown`.

### 6. Context timeout cho mọi query ES

Trước đây dùng `context.Background()` → query treo có thể block
goroutine vô hạn. Thêm `WithTimeout(5s)` ở `GetLogs`, `CountLogs`,
`countErrors`.

### 7. Dashboard — URL động + escape mọi field + sync filter (`dashboard/app.js`)

Bốn thay đổi nhỏ:

- **`API_BASE` / `WS_URL` dynamic**: hardcoded `http://localhost:8080`
  → vỡ khi truy cập qua tên server, reverse proxy, hoặc đổi
  `API_PORT`. Đổi sang `window.location.origin` (auto `ws://`/`wss://`).
- **Escape mọi field** trong `renderTable` (không chỉ `log_message`).
  Trước đây `level` và `service` từ ES được nhúng raw qua template
  literal — defense in depth quan trọng vì pipeline có thể tampered
  hoặc thêm service ngoài đẩy log vào.
- **Whitelist level** trước khi dùng làm class CSS (`badge-${level}`)
  để tránh tạo class lạ từ data ngoài.
- **`fetchStats` truyền filter `app`** hiện tại, để 4 con số INFO/
  WARN/ERROR/total khớp với bảng đang lọc bên dưới.
- `escapeHtml` thêm `'` và backtick.

## ⚪ Vấn đề thấy nhưng KHÔNG sửa (out of scope)

- **CORS `Allow-Origin: *`** — ok cho dev, cần config khi prod. Có
  thể đọc từ env `CORS_ALLOWED_ORIGINS`.
- **WebSocket `CheckOrigin: return true`** — tương tự, nên check
  origin trong production.
- **Elasticsearch bind 0.0.0.0** — README đã cảnh báo "không public
  ES ra Internet"; trong dev/local thì giữ nguyên.
- **Demo Go service** `rand.Float64()` không seed — Go 1.20+ tự seed
  nên không sao. Go version trong Dockerfile cần ≥ 1.20 (`go.mod`
  declare `go 1.22` → OK).
- **Filebeat queue size 1000 events** — đủ cho demo, scale thật cần
  tính lại theo throughput.

## Cấu trúc commit

```
ccdbf86 fix(api): graceful shutdown for alerting engine
e219f43 fix(dashboard): drop hardcoded localhost URLs, escape every ES field, sync stats with filter
9046b09 refactor(api): reuse alert handler instance for WebSocket route
05ad899 fix(api): race condition + dead-conn buildup + dead code in alerting
f70ba09 fix(compose): pin internal ES port for api-server
c226de3 fix(api): plug response body leak in CountLogs + add ES query timeouts
```

## Kiểm thử đề xuất sau khi merge

1. **Smoke test docker-compose** với `ES_PORT=9201` trong `.env` —
   trước fix sẽ thấy api-server log lỗi kết nối, sau fix phải ok.
2. **Stress `/api/logs/count`**: `ab -n 1000 -c 20 http://localhost:8080/api/logs/count`
   — quan sát số kết nối tới ES qua `_nodes/stats/http`, không được
   tăng vô hạn.
3. **Dashboard qua reverse proxy** (Nginx forward `/` → api-server) —
   trước fix dashboard không gọi được API; sau fix phải hoạt động.
4. **Race detector**: `cd api-server && go test -race ./...` (nếu có
   thêm test) hoặc `go run -race main.go`.
5. **WS dead conn**: mở dashboard, kill tab, đợi tới khi alert tiếp
   theo trigger — log không được nhả "write error" lặp lại.
