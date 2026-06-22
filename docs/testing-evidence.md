# Testing Evidence

Record commands, results, and evidence used to validate the system.

Thứ tự verify chi tiết theo tuần nằm ở
[`docs/one-month-defense-roadmap.md`](one-month-defense-roadmap.md).

## Evidence Template

- Feature:
- Command or scenario:
- Expected result:
- Actual result:
- Status:
- Notes:

## Current Evidence

### Bước 10 - Kết quả end-to-end

- Feature: Dashboard static serving
- Command or scenario:
  ```bash
  curl -s http://localhost:8080/ | head -5
  ```
- Expected result: output chứa `<!DOCTYPE html>`.
- Actual result: Pass ngày 2026-06-17. `curl -s -o /tmp/log-system-dashboard-index.html -w "%{http_code} %{time_total}\n" http://localhost:8080/` trả `200 0.012889`; 5 dòng đầu chứa `<!DOCTYPE html>`.
- Status: Pass
- Notes: Xác nhận API container serve được dashboard static tại `/`.

- Feature: API filter `level=ERROR`
- Command or scenario:
  ```bash
  curl -s "http://localhost:8080/api/logs?level=ERROR&size=3" | python3 -m json.tool
  ```
- Expected result: response có `data`, `total`, `page`, `size`; các record trả về có `level` là `ERROR`.
- Actual result: Pass ngày 2026-06-17. Response trả `data`, `total`, `page`, `size`; 3 record mẫu đều có `level: "ERROR"`.
- Status: Pass
- Notes: Query mẫu trả `total: 45` tại thời điểm chạy.

- Feature: API filter service
- Command or scenario:
  ```bash
  curl -s "http://localhost:8080/api/logs?app=demo-node&size=3" | python3 -m json.tool
  ```
- Expected result: các record trả về thuộc service `demo-node`.
- Actual result: Pass ngày 2026-06-17. Response `data` trả 3 record mẫu đều có `service: "demo-node"` và `total: 200`.
- Status: Pass
- Notes: Xác nhận filter `app=demo-node` hoạt động.

- Feature: API count theo level
- Command or scenario:
  ```bash
  curl -s "http://localhost:8080/api/logs/count" | python3 -m json.tool
  ```
- Expected result: response có `INFO`, `WARN`, `ERROR`, `total`.
- Actual result: Pass ngày 2026-06-17. Response mẫu sau incident replay: `{"ERROR":78,"INFO":235,"WARN":102,"from":"now-1h","to":"now","total":415}`.
- Status: Pass
- Notes: Dữ liệu này feed stats bar trên dashboard.

- Feature: Dynamic threshold + alerting
- Command or scenario:
  ```bash
  curl -s -X POST http://localhost:8080/api/alerts/config \
    -H "Content-Type: application/json" \
    -d '{"threshold":5}' | python3 -m json.tool

  ./scripts/trigger-error-spike.sh 20

  docker compose logs api-server | grep -i alert
  ```
- Expected result: config update thành công; script ghi 20 ERROR JSON Lines vào `./logs/demo-node/app.log`; sau khoảng `ALERT_CHECK_INTERVAL_SECONDS` giây log có dòng `[alerting] alert sent - count=X threshold=5`.
- Actual result: Pass ngày 2026-06-17. `POST /api/alerts/config` trả `{"config":{"threshold":5,"window_seconds":300,"cooldown_seconds":60},"status":"updated"}`. Lần chạy đầu của script bắt lỗi `Permission denied` khi ghi trực tiếp vào host log file; script đã được sửa để fallback qua `docker compose exec -T demo-node`. Sau đó `./scripts/trigger-error-spike.sh 20` trả `wrote 20 ERROR logs to ./logs/demo-node/app.log via container=demo-node ...`. API query `level=ERROR&q=INCIDENT_REPLAY&size=5` trả `total: 20`. API server log có `[alerting] alert sent - count=68 threshold=5`.
- Status: Pass
- Notes: Alert sau replay xuất hiện sau khoảng 33 giây vì trước đó engine đã gửi alert và đang trong `cooldown_seconds=60`; đây là hành vi dedup/cooldown kỳ vọng. Khi demo, nên hạ threshold và trigger incident khi cooldown không còn active.

- Feature: Response time thực tế
- Command or scenario:
  ```bash
  time curl -s "http://localhost:8080/api/logs?size=100" -o /dev/null
  time curl -s "http://localhost:8080/api/health" -o /dev/null
  ```
- Expected result: Có số liệu `real/user/sys` để điền vào `docs/deployment.md`.
- Actual result: Pass ngày 2026-06-17. `curl -s -o /tmp/log-system-api-logs-100.json -w "%{http_code} %{time_total}\n" "http://localhost:8080/api/logs?size=100"` trả `200 0.033992`; `curl -s -o /tmp/log-system-api-health.json -w "%{http_code} %{time_total}\n" http://localhost:8080/api/health` trả `200 0.007250`.
- Status: Pass
- Notes: Số liệu đo bằng `curl` trên máy dev hiện tại; không phải benchmark tải nặng.

### Rà soát chốt MVP ngày 2026-06-22

- Feature: Compose runtime health
- Command or scenario: `docker compose ps`
- Expected result: Elasticsearch, Logstash, Filebeat, demo-node, demo-go và api-server đều chạy; service có healthcheck ở trạng thái healthy.
- Actual result: Cả sáu service đều `Up`; sáu service đều báo `healthy`.
- Status: Pass
- Notes: Đây là kiểm tra đọc trạng thái, không rebuild hoặc thay đổi container.

- Feature: Dashboard runtime
- Command or scenario: Mở `http://localhost:8080` bằng browser và kiểm tra accessibility snapshot, console, network.
- Expected result: WebSocket connected; stats, table và pagination có data; REST request thành công.
- Actual result: Dashboard hiển thị `Connected`, 1.019 log trong cửa sổ thống kê, 20 dòng ở trang 1/51; `/api/logs` và `/api/logs/count` trả HTTP 200.
- Status: Pass
- Notes: Console chỉ có lỗi không ảnh hưởng chức năng: `GET /favicon.ico` trả 404.

- Feature: Static/build checks
- Command or scenario:
  ```bash
  docker compose config --quiet
  node --check services/demo-node/index.js
  bash -n scripts/trigger-error-spike.sh
  cd api-server && GOCACHE=/tmp/log-system-go-build-cache go test ./...
  cd services/demo-go && GOCACHE=/tmp/log-system-demo-go-cache go test ./...
  ```
- Expected result: config hợp lệ; JavaScript/shell hợp lệ; hai Go module compile và test pass.
- Actual result: Tất cả lệnh pass. API có 5 test cho alert config; các package còn lại không có test file. Demo Go không có test file.
- Status: Pass
- Notes: Automated coverage còn mỏng; E2E evidence là bằng chứng chính cho logs query, pipeline và WebSocket alert.

## Khoảng trống bằng chứng còn lại

| Khoảng trống | Mức độ với MVP | Hành động |
|---|---|---|
| Clean-clone từ checkout độc lập | Bắt buộc trước bảo vệ | Chạy README flow, ghi thời gian và blocker môi trường |
| Screenshot dashboard bình thường | Bắt buộc cho slide/report | Chụp ở viewport rõ bảng và stats |
| Screenshot filter ERROR | Bắt buộc cho slide/report | Chụp sau khi áp dụng filter |
| Screenshot alert banner | Bắt buộc cho slide/demo backup | Trigger incident khi cooldown đã hết rồi chụp banner |
| Unit test logs query, cooldown/dedup, WebSocket | Nên có, không chặn MVP | Bổ sung sau khi hoàn tất report/slide nếu còn thời gian |
| Load/stress benchmark | Ngoài phạm vi MVP | Chỉ đề xuất ở hướng phát triển, không suy diễn từ số `curl` |
