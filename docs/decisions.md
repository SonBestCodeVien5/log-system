# Quyết định kỹ thuật

Tài liệu này ghi lại các lựa chọn đã hình thành MVP, phương án thay thế và hệ
quả. Đây là nguồn ngắn gọn để giải thích thiết kế trong báo cáo và khi bảo vệ.

## 1. Dùng JSON Lines làm contract log

- **Bối cảnh:** Hai demo service viết bằng Node.js và Go phải tạo log mà cùng một pipeline có thể xử lý ổn định.
- **Phương án cân nhắc:** text tự do + Grok; JSON array; JSON Lines.
- **Quyết định:** mỗi dòng là một JSON object gồm `timestamp`, `level`, `service`, `message`, `metadata`.
- **Lý do:** tách record theo dòng, dễ append/tail, parse có cấu trúc và không phụ thuộc Grok cho field cơ bản.
- **Hệ quả:** producer phải giữ đúng schema; malformed JSON được đánh dấu parse error. Grok chỉ enrich `log_message`, không quyết định log có được lưu hay không.
- **Bằng chứng:** `services/demo-node/index.js`, `services/demo-go/main.go`, `logstash/pipeline/logstash.conf`.

## 2. Tách collector và processor bằng Filebeat + Logstash

- **Bối cảnh:** Ứng dụng không nên tự quản retry, offset và kết nối trực tiếp đến Elasticsearch.
- **Phương án cân nhắc:** service gọi Elasticsearch; chỉ dùng Filebeat; Filebeat qua Logstash.
- **Quyết định:** service ghi file, Filebeat tail và ship qua Logstash trước khi index.
- **Lý do:** tách concern; Filebeat giữ registry/offset; Logstash tập trung parse, normalize timestamp và enrich.
- **Hệ quả:** có thêm hai runtime component và cần healthcheck/debug theo từng chặng, đổi lại pipeline dễ quan sát và mở rộng.
- **Bằng chứng:** `filebeat/filebeat.yml`, `logstash/pipeline/logstash.conf`, `docker-compose.yml`.

## 3. Dùng Elasticsearch cho lưu trữ và tìm kiếm log

- **Bối cảnh:** workload chính là full-text search, exact filter và time range trên log.
- **Phương án cân nhắc:** MySQL/PostgreSQL; file scan; Elasticsearch.
- **Quyết định:** lưu theo index ngày `logs-YYYY.MM.dd` và query `logs-*`.
- **Lý do:** inverted index, Query DSL và date math phù hợp dữ liệu log hơn truy vấn `LIKE` hoặc scan file.
- **Hệ quả:** cần nhiều RAM hơn database quan hệ và phải bổ sung ILM/retention nếu triển khai production.
- **Bằng chứng:** `logstash/pipeline/logstash.conf`, `api-server/handlers/logs.go`.

## 4. Dùng Go + Gin cho API và alerting

- **Bối cảnh:** một process cần phục vụ REST, WebSocket và goroutine kiểm tra alert định kỳ.
- **Phương án cân nhắc:** Spring Boot; Node.js; Go + Gin.
- **Quyết định:** dùng Go 1.22, Gin, go-elasticsearch v8 và gorilla/websocket.
- **Lý do:** goroutine/context phù hợp background loop và graceful shutdown; binary nhỏ, dependency trực tiếp và dễ container hóa.
- **Hệ quả:** nhóm phải tự quản concurrency, mutex, timeout và WebSocket lifecycle rõ ràng.
- **Bằng chứng:** `api-server/main.go`, `api-server/alerting/engine.go`, `api-server/handlers/alerts.go`.

## 5. Alert theo sliding window, threshold và cooldown

- **Bối cảnh:** cần phát hiện spike ERROR nhưng tránh gửi cảnh báo lặp ở mỗi chu kỳ polling.
- **Phương án cân nhắc:** alert trên từng ERROR; fixed batch; sliding window + cooldown.
- **Quyết định:** định kỳ đếm ERROR trong cửa sổ gần nhất, alert khi `count > threshold`, rồi deduplicate theo cooldown.
- **Lý do:** phản ánh mật độ lỗi trong thời gian gần và kiểm soát alert fatigue.
- **Hệ quả:** alert có độ trễ tối đa xấp xỉ check interval khi không bị cooldown; config phải được bảo vệ bằng mutex và check/write dedup phải atomic.
- **Bằng chứng:** `api-server/alerting/engine.go`; incident replay ngày 2026-06-17 ghi nhận alert sau ~33 giây vì cooldown cũ còn active.

## 6. Cho phép cập nhật alert config từng phần

- **Bối cảnh:** người vận hành cần đổi threshold/window/cooldown mà không restart API.
- **Phương án cân nhắc:** chỉ đọc `.env` lúc start; bắt buộc gửi toàn bộ config; partial update.
- **Quyết định:** `POST /api/alerts/config` nhận một hoặc nhiều field dương; field vắng mặt giữ giá trị hiện tại; body rỗng bị từ chối.
- **Lý do:** thuận tiện cho dashboard và demo, đồng thời validate trước khi mutate để tránh trạng thái cập nhật dở dang.
- **Hệ quả:** config runtime mất khi restart và quay về biến môi trường; persistence là hướng phát triển sau MVP.
- **Bằng chứng:** `api-server/handlers/alerts.go`, `api-server/handlers/alerts_test.go`.

## 7. Giữ dashboard bằng HTML/CSS/JavaScript thuần

- **Bối cảnh:** MVP cần màn hình vận hành trực tiếp, không cần frontend build pipeline.
- **Phương án cân nhắc:** React/Vue; server-rendered templates; vanilla dashboard.
- **Quyết định:** ba file `index.html`, `app.js`, `style.css`, được Go API serve cùng origin.
- **Lý do:** giảm dependency và thao tác triển khai; REST và WebSocket URL suy ra từ `window.location` nên chạy được khi đổi host/protocol.
- **Hệ quả:** state/rendering quản lý thủ công; phù hợp scope MVP nhưng cần refactor nếu UI tăng lớn.
- **Bằng chứng:** `dashboard/`, `api-server/main.go`.

## 8. Dùng incident replay thay vì chờ log ERROR ngẫu nhiên

- **Bối cảnh:** demo alert sẽ thiếu ổn định nếu phụ thuộc xác suất ERROR 15% của demo services.
- **Phương án cân nhắc:** tăng xác suất trong source; thêm demo API; script append JSON Lines hợp lệ.
- **Quyết định:** dùng `scripts/trigger-error-spike.sh` tạo batch ERROR có `batch_id` và fallback ghi qua container khi host log bị giới hạn quyền.
- **Lý do:** deterministic, không đổi contract và đi qua đúng pipeline production-like.
- **Hệ quả:** script là công cụ local/demo; cần chờ cooldown hết trước khi rehearsal.
- **Bằng chứng:** `scripts/trigger-error-spike.sh`, `docs/testing-evidence.md`.

## 9. Giới hạn bảo mật ở mức local MVP

- **Bối cảnh:** stack phục vụ demo nội bộ, trong khi production cần kiểm soát truy cập chặt hơn.
- **Quyết định:** bật Elasticsearch Basic Security nhưng CORS và WebSocket origin vẫn mở cho local dashboard; không triển khai auth/RBAC cho API.
- **Lý do:** giữ scope phù hợp thời gian tốt nghiệp và tập trung chứng minh luồng log/alert.
- **Hệ quả:** không expose trực tiếp ra Internet; production cần TLS/reverse proxy, auth/RBAC, origin allowlist, secret management và ILM.
- **Bằng chứng:** `docker-compose.yml`, `api-server/middleware/cors.go`, `api-server/handlers/alerts.go`.
