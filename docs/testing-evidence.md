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

### Bước 10 - End-to-end verification plan

- Feature: Dashboard static serving
- Command or scenario:
  ```bash
  curl -s http://localhost:8080/ | head -5
  ```
- Expected result: output chứa `<!DOCTYPE html>`.
- Actual result: Chưa chạy trong phiên này.
- Status: Pending
- Notes: Chạy sau khi rebuild/start `api-server`.

- Feature: API filter `level=ERROR`
- Command or scenario:
  ```bash
  curl -s "http://localhost:8080/api/logs?level=ERROR&size=3" | python3 -m json.tool
  ```
- Expected result: response có `data`, `total`, `page`, `size`; các record trả về có `level` là `ERROR`.
- Actual result: Chưa chạy trong phiên này.
- Status: Pending
- Notes: Dùng để xác nhận Bước 8 filter theo level.

- Feature: API filter service
- Command or scenario:
  ```bash
  curl -s "http://localhost:8080/api/logs?app=demo-node&size=3" | python3 -m json.tool
  ```
- Expected result: các record trả về thuộc service `demo-node`.
- Actual result: Chưa chạy trong phiên này.
- Status: Pending
- Notes: Dùng để xác nhận filter theo service.

- Feature: API count theo level
- Command or scenario:
  ```bash
  curl -s "http://localhost:8080/api/logs/count" | python3 -m json.tool
  ```
- Expected result: response có `INFO`, `WARN`, `ERROR`, `total`.
- Actual result: Chưa chạy trong phiên này.
- Status: Pending
- Notes: Dữ liệu này feed stats bar trên dashboard.

- Feature: Dynamic threshold + alerting
- Command or scenario:
  ```bash
  curl -s -X POST http://localhost:8080/api/alerts/config \
    -H "Content-Type: application/json" \
    -d '{"threshold": 5}' | python3 -m json.tool

  docker compose logs api-server | grep -i alert
  ```
- Expected result: config update thành công; sau 10-15 giây log có dòng `[alerting] alert sent - count=X threshold=5` khi số ERROR vượt ngưỡng.
- Actual result: Chưa chạy trong phiên này.
- Status: Pending
- Notes: Đây là verification quan trọng nhất cho Bước 9.

- Feature: Response time thực tế
- Command or scenario:
  ```bash
  time curl -s "http://localhost:8080/api/logs?size=100" -o /dev/null
  time curl -s "http://localhost:8080/api/health" -o /dev/null
  ```
- Expected result: Có số liệu `real/user/sys` để điền vào `docs/deployment.md`.
- Actual result: Chưa chạy trong phiên này.
- Status: Pending
- Notes: Không ghi số liệu ước lượng; chỉ điền sau khi đo thật trên máy dev.
