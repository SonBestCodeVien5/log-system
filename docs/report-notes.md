# Report Notes

Working notes for the graduation report. Use this file to collect report-ready explanations from implemented features, architecture decisions, testing evidence, and known limitations.

Roadmap 1 tháng cuối để biến các ghi chú này thành demo/report hoàn chỉnh:
[`docs/one-month-defense-roadmap.md`](one-month-defense-roadmap.md).

## Feature Summaries

Record each feature with:

- Feature name:
- Purpose:
- User/system problem solved:
- Main implementation points:
- Evidence:
- Limitations:

### Centralized Logging Platform MVP

- Feature name: Centralized Logging Platform MVP
- Purpose: Gom log từ nhiều service, tìm kiếm/filter qua dashboard và cảnh báo realtime khi ERROR vượt ngưỡng.
- User/system problem solved: Người vận hành không cần mở từng file log riêng lẻ; có một nơi tập trung để xem, lọc và nhận cảnh báo.
- Main implementation points:
  - Demo services Node.js và Go ghi JSON Lines vào `./logs`.
  - Filebeat tail file log và gửi tới Logstash.
  - Logstash parse JSON, enrich message và index vào Elasticsearch `logs-*`.
  - Go API dùng gin query Elasticsearch, expose REST endpoint và WebSocket.
  - Alerting engine dùng sliding window, deduplication và dynamic threshold.
  - Dashboard HTML/CSS/JS thuần hiển thị bảng log, stats, filter, pagination và alert banner.
- Evidence:
  - Cần bổ sung output Bước 10 trong `docs/testing-evidence.md`.
  - Cần screenshot dashboard sau khi mở `http://localhost:8080`.
- Limitations:
  - Chưa có authentication cho dashboard/API.
  - Chưa cấu hình retention/ILM production cho Elasticsearch.
  - Số liệu hiệu năng cần đo trên môi trường thật trước khi đưa vào báo cáo.

## Architecture Explanation Drafts

Add reusable paragraphs that explain the centralized logging pipeline, API server, dashboard, and alerting engine in report-friendly language.

## Evidence Placeholders

Track screenshots, command outputs, demo flows, or diagrams that should be captured later for the report.

## Demo Script 5 Phút

1. Mở dashboard tại `http://localhost:8080` và chỉ ra log đang được refresh tự động.
2. Filter `ERROR`, sau đó filter `demo-node`, để chứng minh API query và filter hoạt động.
3. Giảm threshold cảnh báo xuống `5` bằng input trên dashboard.
4. Chờ 10-15 giây, khi banner đỏ xuất hiện thì giải thích WebSocket nhận alert realtime.
5. Mở `api-server/alerting/engine.go` và giải thích ngắn:
   - Sliding window đếm ERROR trong khoảng thời gian gần nhất.
   - Deduplication tránh gửi alert lặp liên tục.
   - Dynamic threshold cho phép đổi ngưỡng mà không restart server.

## Câu Hỏi Cần Ôn Trước Khi Bảo Vệ

- Vì sao chọn Go thay vì Spring Boot cho API/alerting?
- Vì sao dùng Elasticsearch thay vì MySQL/PostgreSQL để tìm kiếm log?
- Filebeat và Logstash giải quyết vấn đề gì trong pipeline?
- Sliding window khác gì polling cố định?
- Alert deduplication tránh Alert Fatigue như thế nào?

Nguồn ôn tập chính: [`docs/knowledge-base.md`](knowledge-base.md).
Lịch ôn theo tuần: [`docs/one-month-defense-roadmap.md`](one-month-defense-roadmap.md).
