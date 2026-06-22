# Outline slide bảo vệ — 10 trang

Outline này dành cho bài trình bày khoảng 7–10 phút, trong đó demo trực tiếp
khoảng 5 phút có thể được chạy sau slide 7. Mỗi slide chỉ giữ một thông điệp
chính; chi tiết kỹ thuật để dành cho lời nói và Q&A.

## Slide 1 — Đề tài và bài toán

**Tiêu đề:** Centralized Logging Platform

**Thông điệp:** Log phân tán làm chậm việc phát hiện và điều tra lỗi.

**Nội dung trên slide:**

- Log nằm ở nhiều service/file/container.
- Khó lọc theo thời gian, level và service.
- Không có cảnh báo tập trung khi ERROR tăng đột biến.
- Mục tiêu: collect → search → visualize → alert.

**Lời nói gợi ý:** giới thiệu bài toán vận hành trước khi nói công nghệ.

## Slide 2 — Yêu cầu và phạm vi MVP

**Thông điệp:** MVP chứng minh một luồng vận hành hoàn chỉnh, không cố mô phỏng toàn bộ production.

**Nội dung trên slide:**

- Hai producer: Node.js và Go.
- JSON Lines thống nhất.
- Filter/search/pagination và stats.
- WebSocket alert + dynamic config.
- Docker Compose và incident replay.

**Không ghi như tính năng:** auth/RBAC, Kubernetes, AI analysis, multi-tenancy.

## Slide 3 — Kiến trúc tổng thể

**Thông điệp:** Mỗi tầng có một trách nhiệm rõ ràng.

**Visual chính:**

```text
Services → Filebeat → Logstash → Elasticsearch → Go API → Dashboard
                                      └────────→ AlertEngine → WebSocket
```

**Gợi ý trình bày:** nói một câu cho mỗi arrow; không đọc danh sách công nghệ.

## Slide 4 — Contract log và ingest pipeline

**Thông điệp:** JSON Lines là contract ổn định; Grok chỉ enrich, không quyết định việc giữ log.

**Nội dung trên slide:**

- JSON mẫu với năm field chính.
- Filebeat tail + registry.
- Logstash parse, uppercase level, normalize `@timestamp`.
- Index `logs-YYYY.MM.dd`.

**Visual:** trước/sau Logstash, highlight `log_message`, `error_code`, `environment`.

## Slide 5 — Elasticsearch Query và Go API

**Thông điệp:** API chuyển nhu cầu vận hành thành Query DSL có giới hạn và timeout.

**Nội dung trên slide:**

- Exact: `level.keyword`, `service.keyword`.
- Full-text: match `log_message`.
- Range: `@timestamp`; sort newest first.
- Pagination 20, tối đa 100; query timeout 5 giây.
- Endpoint: health, logs, count, alert config.

**Visual:** request `/api/logs?...` → bool query → response chuẩn.

## Slide 6 — Dashboard và realtime flow

**Thông điệp:** Một màn hình cho quan sát log và trạng thái alert.

**Nội dung trên slide:**

- Stats INFO/WARN/ERROR.
- Filter, search, pagination, auto-refresh 10 giây.
- WebSocket Connected và alert banner.
- Alert config threshold/window/cooldown.

**Visual:** screenshot H1 hoặc H2; không nhồi source code.

## Slide 7 — Alerting engine

**Thông điệp:** Sliding window phát hiện spike; cooldown kiểm soát alert fatigue.

**Visual chính:** timeline ba lần check.

```text
t=0   count=3   → no alert
t=10  count=15  → alert + record sent time
t=20  count=14  → skip vì cooldown
```

**Điểm kỹ thuật để nói:**

- Điều kiện `count > threshold`.
- Atomic check/write bằng một mutex.
- Snapshot WebSocket clients, network I/O ngoài lock.
- Dynamic config không cần restart.

## Slide 8 — Incident replay và demo

**Thông điệp:** Alert được kiểm thử deterministic, không phụ thuộc ERROR random.

**Nội dung trên slide:**

1. Đặt threshold phù hợp.
2. Chạy `./scripts/trigger-error-spike.sh 20`.
3. Filebeat → Logstash → ES.
4. API tìm thấy batch; WebSocket hiện banner.

**Visual:** screenshot H3 và một dòng evidence `alert sent`.

**Demo transition:** chuyển sang browser sau slide này hoặc chạy demo ngay trong slide.

## Slide 9 — Kết quả kiểm thử

**Thông điệp:** MVP có evidence từ runtime, không chỉ có source code.

**Bảng ngắn:**

| Hạng mục | Kết quả |
|---|---|
| Compose | 6/6 service healthy |
| API filters/count | Pass |
| Incident replay | 20 ERROR được tìm thấy |
| Alert | `alert sent`, cooldown hoạt động |
| Logs API | `0.033992s` cho size=100 |
| Dashboard | `0.012889s` |

**Chú thích bắt buộc:** số đo dev bằng request đơn, không phải load benchmark.

## Slide 10 — Kết luận, hạn chế và hướng phát triển

**Thông điệp:** MVP hoàn thành luồng centralized logging và có đường nâng cấp rõ ràng.

**Kết luận:**

- Thu thập đa service.
- Search/filter tập trung.
- Alert realtime có dedup.
- Incident demo lặp lại được.

**Hạn chế/hướng phát triển:**

- Auth/RBAC, TLS, origin allowlist.
- ILM/retention, backup.
- HA/scale và load testing.
- Notification channels và observability sâu hơn.

**Câu chốt gợi ý:** “Giá trị của MVP không nằm ở số lượng màn hình, mà ở việc em
đã xây dựng và chứng minh được toàn bộ đường đi của một log và một cảnh báo.”

## Asset checklist cho deck

- [ ] H1: dashboard bình thường.
- [ ] H2: dashboard filter ERROR/demo-node.
- [ ] H3: alert banner incident replay.
- [ ] Sơ đồ kiến trúc đã vẽ lại bằng shape/vector.
- [ ] JSON before/after pipeline.
- [ ] Timeline sliding window.
- [ ] Bảng evidence slide 9.
- [ ] Logo trường, tên sinh viên, GVHD và template chính thức.
