# Roadmap hiện tại

Tài liệu này phản ánh trạng thái repo sau commit `0ceef5f`
(`feat: add dashboard and deployment guide`). Mục tiêu hiện tại là không thêm
scope lớn, ưu tiên xác minh end-to-end, ghi bằng chứng, hoàn thiện tài liệu và
chuẩn bị bảo vệ. Nếu thêm code, chỉ nên là một kịch bản incident nhỏ để chủ động
trigger alert khi demo.

Roadmap chi tiết cho 1 tháng trước báo cáo/bảo vệ nằm ở
[`docs/one-month-defense-roadmap.md`](one-month-defense-roadmap.md).

## Tổng quan trạng thái

| Nhóm việc | Trạng thái | Ghi chú |
|---|---|---|
| Infrastructure | Hoàn thành | `docker-compose.yml`, Elasticsearch, Logstash, Filebeat |
| Pipeline | Hoàn thành | JSON Lines -> Logstash parse/enrich -> `logs-*` |
| Demo services | Hoàn thành | Node.js + Go sinh log liên tục vào `./logs` |
| Go API | Hoàn thành | REST API, WebSocket, alerting engine |
| Dashboard | Hoàn thành, cần verify runtime | Log viewer, filter, pagination, stats, alert banner |
| Docs | Đang hoàn thiện | README đã có hướng dẫn triển khai; còn số liệu thực tế |
| Bảo vệ | Chưa làm | Cần demo script, câu trả lời, test clone sạch |

## Bước 10 - End-to-end test

Mục tiêu: xác nhận toàn bộ luồng từ demo services đến dashboard và alerting
hoạt động đúng trên môi trường Docker Compose.

### 10.1 Verify dashboard serve được

```bash
curl -s http://localhost:8080/ | head -5
```

Kết quả đúng: thấy `<!DOCTYPE html>`.

### 10.2 Test filter API

```bash
# Filter ERROR
curl -s "http://localhost:8080/api/logs?level=ERROR&size=3" | python3 -m json.tool

# Filter theo service
curl -s "http://localhost:8080/api/logs?app=demo-node&size=3" | python3 -m json.tool

# Đếm theo level
curl -s "http://localhost:8080/api/logs/count" | python3 -m json.tool
```

Kết quả đúng:

- `/api/logs` trả `data`, `total`, `page`, `size`.
- Filter `level=ERROR` chỉ trả log ERROR.
- Filter `app=demo-node` chỉ trả log của service `demo-node`.
- `/api/logs/count` trả `INFO`, `WARN`, `ERROR`, `total`.

### 10.3 Test alerting

Giảm threshold để dễ trigger:

```bash
curl -s -X POST http://localhost:8080/api/alerts/config \
  -H "Content-Type: application/json" \
  -d '{"threshold": 5}' | python3 -m json.tool
```

Chờ 10-15 giây rồi kiểm tra log:

```bash
docker compose logs api-server | grep -i alert
```

Kết quả mong đợi:

```text
[alerting] alert sent - count=X threshold=5
```

Đồng thời dashboard nên hiển thị banner đỏ khi WebSocket nhận message
`error_spike`.

### 10.4 Đo response time thực tế

```bash
time curl -s "http://localhost:8080/api/logs?size=100" -o /dev/null
time curl -s "http://localhost:8080/api/health" -o /dev/null
```

Ghi lại kết quả vào `docs/testing-evidence.md` và bảng hiệu năng trong
`docs/deployment.md`.

### 10.5 Git sau khi verify

Nếu Bước 10 pass, commit phần tài liệu bằng message riêng, ví dụ:

```bash
git add docs/project-roadmap.md docs/testing-evidence.md docs/deployment.md docs/architecture.md
git commit -m "docs: record end-to-end verification plan"
git push
```

## Bước 11 - Hoàn thiện docs

Mục tiêu: biến repo thành tài liệu có thể clone và chạy được ngay.

Checklist chi tiết theo tuần nằm ở
[`docs/one-month-defense-roadmap.md`](one-month-defense-roadmap.md).

Việc cần làm:

- Điền số liệu thực tế từ Bước 10.4 vào `docs/deployment.md`.
- Ghi output quan trọng từ Bước 10.1-10.3 vào `docs/testing-evidence.md`.
- Cập nhật `docs/report-notes.md` bằng kịch bản demo và ảnh chụp màn hình.
- Đảm bảo README root vẫn chạy được theo flow clone -> copy `.env` -> compose up.

Acceptance:

- Người khác đọc README có thể chạy hệ thống mà không cần hỏi thêm.
- Tài liệu có bằng chứng command/output, không chỉ mô tả.
- Các số liệu hiệu năng là số đo thật, không phải ước lượng.

## Bước 12 - Chuẩn bị bảo vệ

Demo script 5 phút:

1. Mở dashboard tại `http://localhost:8080`, cho thấy log đang cập nhật.
2. Filter `ERROR`, sau đó filter `demo-node`, cho thấy API search/filter hoạt động.
3. Giảm threshold xuống `5`, chờ khoảng 10-15 giây.
4. Khi banner đỏ xuất hiện, giải thích WebSocket alert.
5. Mở nhanh `api-server/alerting/engine.go` để giải thích sliding window,
   deduplication và dynamic threshold.

Cần chuẩn bị:

- Đọc lại `docs/knowledge-base.md`.
- Chuẩn bị câu trả lời cho lựa chọn Go thay Spring, Elasticsearch thay SQL,
  Filebeat/Logstash thay ship trực tiếp từ app, sliding window, dedup alert.
- Test clone sạch:

```bash
git clone git@github.com:SonBestCodeVien5/log-system.git fresh
cd fresh
cp .env.example .env
sudo sysctl -w vm.max_map_count=262144
docker compose up -d
```

Kết quả tốt nhất: dashboard mở được trong dưới 5 phút và có log hiển thị.

## Thứ tự ưu tiên tiếp theo

1. Chạy Bước 10.1-10.3 để xác nhận dashboard/API/alerting.
2. Chạy Bước 10.4 để lấy số liệu thật.
3. Điền số liệu và output vào docs.
4. Commit riêng phần docs verification.
5. Test clone sạch trước buổi bảo vệ.
