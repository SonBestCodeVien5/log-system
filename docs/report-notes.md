# Nội dung báo cáo MVP

Tài liệu này là bản thảo nội dung kỹ thuật có thể chuyển sang mẫu báo cáo của
trường. Bằng chứng lệnh và số đo chi tiết nằm trong
[`testing-evidence.md`](testing-evidence.md); outline thuyết trình nằm trong
[`slide-outline.md`](slide-outline.md).

## 1. Tóm tắt đề tài

Đề tài xây dựng một nền tảng log tập trung giúp thu thập log từ nhiều dịch vụ,
chuẩn hóa dữ liệu, lưu trữ để tìm kiếm nhanh và phát cảnh báo realtime khi số
lượng lỗi vượt ngưỡng. MVP sử dụng hai service mô phỏng bằng Node.js và Go,
Elastic Stack cho ingest/storage, Go API cho truy vấn và alerting, cùng dashboard
HTML/CSS/JavaScript thuần cho người vận hành.

Kết quả chính là một luồng có thể quan sát từ đầu đến cuối: service ghi JSON
Lines, Filebeat thu thập, Logstash parse/enrich, Elasticsearch index, REST API
lọc/phân trang, dashboard hiển thị, và WebSocket đẩy cảnh báo ERROR spike.

## 2. Bài toán và mục tiêu

Khi mỗi service giữ log riêng, người vận hành phải mở nhiều file hoặc container
để tìm nguyên nhân sự cố. Cách làm này tốn thời gian, khó đối chiếu theo mốc thời
gian và không có cơ chế cảnh báo tập trung.

MVP đặt ra năm mục tiêu quan sát được:

1. Thu thập log từ ít nhất hai service dùng công nghệ khác nhau.
2. Chuẩn hóa thành schema chung và lưu theo index thời gian.
3. Tìm kiếm, lọc theo level/service/nội dung/thời gian và phân trang.
4. Cảnh báo realtime khi ERROR trong sliding window vượt threshold.
5. Chạy toàn bộ hệ thống bằng Docker Compose và có kịch bản demo lặp lại được.

## 3. Phạm vi MVP

### Đã hoàn thành

- Hai demo service sinh INFO/WARN/ERROR theo JSON Lines.
- Filebeat tail file và chuyển event đến Logstash.
- Logstash parse JSON, normalize timestamp/level, Grok enrich và index `logs-*`.
- Elasticsearch 8.13.0 lưu index theo ngày.
- Go API cung cấp health, list/filter/search/pagination, count theo level và dynamic alert config.
- Alert engine chạy goroutine, dùng sliding window, cooldown/dedup và WebSocket broadcast.
- Dashboard có stats, filter, search, pagination 20 dòng, auto-refresh 10 giây và alert banner.
- Incident replay tạo ERROR spike chủ động để kiểm thử và demo.

### Ngoài phạm vi và hướng phát triển

- Authentication/RBAC cho API và dashboard.
- TLS, reverse proxy, CORS/origin allowlist và secret manager production.
- Elasticsearch ILM/retention, backup và capacity planning.
- High availability, multi-node Elasticsearch và horizontal scaling.
- Distributed tracing, multi-tenancy, notification email/Slack và AI log analysis.

Các mục trên được trình bày như hạn chế có chủ đích của local MVP, không phải
tính năng đã triển khai.

## 4. Kiến trúc và luồng dữ liệu

```text
demo-node / demo-go
        |
        | JSON Lines: ./logs/<service>/app.log
        v
     Filebeat
        | Beats :5044
        v
     Logstash
        | parse + normalize + enrich
        v
 Elasticsearch (logs-YYYY.MM.dd)
        |
        +--> Go REST API --> Dashboard
        |
        +--> AlertEngine --> WebSocket --> Alert banner
```

Thiết kế tách producer, collector, processor, storage và presentation thành các
tầng độc lập. Demo service chỉ chịu trách nhiệm ghi log. Filebeat quản lý việc
tail và offset. Logstash tập trung xử lý schema. Elasticsearch tối ưu tìm kiếm.
Go API che giấu Query DSL khỏi dashboard và alert engine tái sử dụng cùng nguồn
dữ liệu Elasticsearch để đánh giá lỗi.

## 5. Contract dữ liệu và pipeline

Mỗi log là một JSON object trên một dòng:

```json
{"timestamp":"2026-06-17T10:23:11Z","level":"ERROR","service":"demo-node","message":"Payment gateway timeout","metadata":{"order_id":"789"}}
```

Logstash parse JSON vào object tạm, promote các field thành `@timestamp`,
`level`, `service`, `log_message`, `metadata`, chuyển level sang chữ hoa và thêm
`environment=development`, `stack=log-system`. Grok là bước enrich tùy chọn để
tách `error_code`/`error_detail`; không match Grok không làm mất log.

Index theo ngày `logs-YYYY.MM.dd` giúp giới hạn phạm vi vận hành theo thời gian
và tạo nền tảng cho retention sau này. API query qua alias pattern `logs-*`.

## 6. API và dashboard

`GET /api/logs` xây Elasticsearch bool query từ exact filter `level`, `app`,
full-text `q` và range `from/to`; kết quả sort mới nhất trước và trả response
`data`, `total`, `page`, `size`. Mỗi trang mặc định 20 và tối đa 100 record.

`GET /api/logs/count` đếm riêng INFO/WARN/ERROR trong một giờ gần nhất theo mặc
định. Dashboard dùng endpoint này cho stats bar, gọi list endpoint cho bảng và
tự refresh mỗi 10 giây. REST/WebSocket URL được suy ra từ origin hiện tại, giúp
dashboard chạy qua cùng API host mà không hardcode `localhost`.

## 7. Alerting engine

Alert engine chạy trong một goroutine và kiểm tra Elasticsearch theo chu kỳ.
Với cấu hình mặc định, engine đếm ERROR trong 300 giây gần nhất mỗi 10 giây.
Khi `count > threshold`, engine thực hiện check-and-write atomic trên map
deduplication. Nếu alert cùng loại đã được gửi trong cooldown 60 giây, lần mới
bị bỏ qua; ngược lại message `error_spike` được broadcast qua WebSocket.

Threshold, window và cooldown có thể cập nhật từng phần bằng REST mà không
restart. Mutex bảo vệ config và dedup; danh sách WebSocket client có mutex riêng.
Khi broadcast, engine snapshot client dưới read lock rồi thực hiện I/O ngoài lock
để tránh giữ lock trong thao tác mạng có thể block.

## 8. Incident replay

Demo services tạo ERROR theo xác suất nên không đảm bảo alert xuất hiện đúng lúc
bảo vệ. Script `scripts/trigger-error-spike.sh` giải quyết vấn đề này bằng cách
append một batch ERROR hợp lệ vào đúng file Filebeat đang tail. Mỗi record có
`source=incident-replay`, `batch_id` và sequence nên có thể tìm lại bằng API.

Trong lần verify ngày 2026-06-17, script tạo 20 ERROR, API tìm thấy đủ 20 record
và API server ghi `alert sent`. Alert được quan sát sau khoảng 33 giây vì engine
đã gửi một alert trước đó và cooldown vẫn còn hiệu lực. Kết quả này vừa chứng
minh incident replay hoạt động, vừa minh họa đúng cơ chế chống alert fatigue.

## 9. Kết quả kiểm thử

| Hạng mục | Kết quả |
|---|---|
| Compose runtime | Sáu service healthy ngày 2026-06-22 |
| Dashboard | WebSocket Connected; stats/table/pagination có dữ liệu |
| API filter ERROR | Pass; ba record mẫu đều là ERROR |
| API filter service | Pass; ba record mẫu đều là demo-node |
| Count theo level | Pass; response có INFO/WARN/ERROR/total/from/to |
| Incident replay | Pass; 20 record tìm thấy và alert sent |
| `/api/logs?size=100` | HTTP 200 trong `0.033992s` |
| `/api/health` | HTTP 200 trong `0.007250s` |
| Dashboard `/` | HTTP 200 trong `0.012889s` |
| Static/build checks | Compose config, Node, shell và hai Go module đều pass |

Ba số response time là request đơn trên máy dev ngày 2026-06-17, chỉ dùng để
chứng minh MVP phản hồi được trong môi trường thử nghiệm; không phải benchmark
đồng thời hoặc SLA production.

## 10. Đóng góp kỹ thuật và bài học

- Xây dựng contract JSON Lines dùng chung giữa hai ngôn ngữ.
- Thiết kế và debug pipeline theo từng chặng thay vì xem hệ thống như hộp đen.
- Chuyển yêu cầu filter/pagination thành Elasticsearch Query DSL có timeout.
- Xử lý concurrency cho alert config, dedup và WebSocket clients.
- Tránh giữ mutex trong network I/O và loại bỏ dead connection khi broadcast lỗi.
- Tạo incident replay deterministic và ghi evidence thay vì dựa vào demo ngẫu nhiên.
- Phân biệt rõ cơ chế thiết kế, số đo dev và tuyên bố production.

## 11. Kịch bản demo 5 phút

1. **0:00–0:30:** giới thiệu vấn đề log phân tán và mục tiêu centralized logging.
2. **0:30–1:15:** chỉ sơ đồ service → Filebeat → Logstash → Elasticsearch → API/dashboard.
3. **1:15–2:15:** mở dashboard, chỉ WebSocket Connected, stats và pagination; filter ERROR rồi demo-node.
4. **2:15–3:30:** đặt threshold phù hợp, chạy `./scripts/trigger-error-spike.sh 20`, tìm `INCIDENT_REPLAY` và chờ banner.
5. **3:30–4:30:** giải thích sliding window, điều kiện `count > threshold`, atomic dedup và cooldown.
6. **4:30–5:00:** nêu evidence, giới hạn local MVP và hướng production.

Trước demo cần đảm bảo cooldown cũ đã hết. Nếu banner đến muộn, dùng API result
và dòng `alert sent` làm bằng chứng dự phòng, đồng thời giải thích cooldown.

## 12. Danh sách hình cần chụp

| Mã hình | Nội dung | Trạng thái | Dùng ở đâu |
|---|---|---|---|
| H1 | Dashboard bình thường: Connected, stats, bảng log | Chưa chụp file report-ready | Kiến trúc/kết quả |
| H2 | Filter ERROR hoặc demo-node | Chưa chụp | API/dashboard |
| H3 | Alert banner sau incident replay | Chưa chụp | Alerting/demo |
| H4 | `docker compose ps` sáu service healthy | Tùy chọn | Triển khai |
| H5 | API JSON hoặc testing evidence | Tùy chọn | Kiểm thử |

## 13. Câu hỏi bảo vệ trọng tâm

1. Vì sao dùng Elasticsearch thay SQL cho log search?
2. Vì sao cần cả Filebeat và Logstash?
3. JSON Lines mang lại lợi ích và ràng buộc gì?
4. Exact filter khác full-text search thế nào trong API?
5. Sliding window hoạt động ra sao và độ trễ đến từ đâu?
6. Vì sao điều kiện là `count > threshold` thay vì alert từng ERROR?
7. Atomic dedup tránh double alert như thế nào?
8. Vì sao broadcast không giữ mutex trong lúc ghi WebSocket?
9. Incident replay chứng minh được phần nào của hệ thống?
10. Muốn đưa production cần bổ sung bảo mật, retention và scale ra sao?

Nguồn trả lời chi tiết: [`knowledge-base.md`](knowledge-base.md),
[`decisions.md`](decisions.md) và source files được dẫn trong hai tài liệu đó.

## 14. Checklist chốt báo cáo

- [x] Nội dung kiến trúc, pipeline, API, alerting và incident replay.
- [x] Số đo và bằng chứng E2E có ngày, ngữ cảnh và giới hạn diễn giải.
- [x] Demo script 5 phút và 10 câu hỏi trọng tâm.
- [ ] Clean-clone độc lập và thời gian khởi động.
- [ ] Ba screenshot H1–H3.
- [ ] Chuyển nội dung sang template chính thức của trường.
- [ ] Rehearsal ít nhất ba lần và ghi điểm vấp.
