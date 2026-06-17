# Roadmap 1 tháng trước báo cáo và bảo vệ

Tài liệu này là kế hoạch làm việc cho giai đoạn còn khoảng 1 tháng trước khi báo cáo/bảo vệ dự án Log System. Mục tiêu không phải thêm thật nhiều tính năng, mà là biến hệ thống hiện tại thành một dự án có thể chạy ổn định, có bằng chứng kiểm thử, có câu chuyện kỹ thuật rõ ràng và có một điểm cá nhân đủ nổi bật để bảo vệ.

## Định hướng tổng thể

Trong 1 tháng tới, ưu tiên theo thứ tự:

1. **Chạy chắc MVP end-to-end**: service sinh log, Filebeat tail, Logstash parse/enrich, Elasticsearch index, Go API query, dashboard hiển thị, WebSocket alert hoạt động.
2. **Hiểu sâu phần đã implement**: không học lan man, chỉ học đúng những gì đang có trong repo và có thể bị hỏi khi bảo vệ.
3. **Ghi bằng chứng thật**: command output, response time, screenshot, demo script, lỗi gặp phải và cách xử lý.
4. **Thêm một điểm cá nhân nhỏ**: chủ động tạo kịch bản incident/error spike để demo alert không phụ thuộc may rủi.
5. **Đóng gói báo cáo**: README, tài liệu, slide, demo 5 phút, Q&A.

Tỉ lệ thời gian khuyến nghị:

| Nhóm việc | Tỉ lệ | Ý nghĩa |
|---|---:|---|
| Verify + evidence | 25% | Chứng minh hệ thống chạy thật |
| Học sâu + đọc code | 35% | Bảo vệ được thứ mình làm |
| Docs + report + slide | 20% | Biến repo thành tài liệu tốt nghiệp |
| Feature cá nhân nhỏ | 10% | Tạo dấu ấn riêng |
| Rehearsal + clean clone | 10% | Giảm rủi ro ngày bảo vệ |

## Bản đồ code cần đọc

Đây là danh sách file nên mở song song khi đọc roadmap:

| Chủ đề | Code/tài liệu cần đọc | Cần hiểu được gì |
|---|---|---|
| Docker Compose stack | [`docker-compose.yml`](../docker-compose.yml#L8-L184) | Các service, network, volume, healthcheck, env |
| Demo service Node.js | [`services/demo-node/index.js`](../services/demo-node/index.js#L13-L106) | Ghi JSON Lines, tỉ lệ INFO/WARN/ERROR, interval |
| Demo service Go | [`services/demo-go/main.go`](../services/demo-go/main.go#L83-L145) | Ticker, ghi file, error handling khi marshal/write |
| Filebeat collector | [`filebeat/filebeat.yml`](../filebeat/filebeat.yml#L9-L77) | Tail file, registry, retry, memory queue |
| Logstash pipeline | [`logstash/pipeline/logstash.conf`](../logstash/pipeline/logstash.conf#L12-L87) | Parse JSON, promote fields, date, Grok enrich |
| API routes | [`api-server/main.go`](../api-server/main.go#L71-L94) | `/api/logs`, `/api/logs/count`, `/api/alerts/config`, `/ws/alerts`, static dashboard |
| ES client + server init | [`api-server/main.go`](../api-server/main.go#L28-L60) | Env config, ES connection, goroutine alerting |
| Logs API | [`api-server/handlers/logs.go`](../api-server/handlers/logs.go#L53-L109) | Pagination, ES search, response format |
| ES query builder | [`api-server/handlers/logs.go`](../api-server/handlers/logs.go#L164-L205) | `term`, `match`, `range`, bool query |
| Count API | [`api-server/handlers/logs.go`](../api-server/handlers/logs.go#L115-L159) | Count theo INFO/WARN/ERROR |
| Alerting engine | [`api-server/alerting/engine.go`](../api-server/alerting/engine.go#L45-L77) | State, config, mutex, clients, sent map |
| Sliding window | [`api-server/alerting/engine.go`](../api-server/alerting/engine.go#L82-L136) | Ticker, count ERROR, threshold, broadcast |
| Dedup alert | [`api-server/alerting/engine.go`](../api-server/alerting/engine.go#L141-L153) | Check + write atomic bằng single `Lock` |
| Dynamic threshold | [`api-server/alerting/engine.go`](../api-server/alerting/engine.go#L158-L187) | Update/read config runtime |
| WebSocket handler | [`api-server/handlers/alerts.go`](../api-server/handlers/alerts.go#L33-L58) | Upgrade HTTP sang WebSocket, register client |
| Alert config API | [`api-server/handlers/alerts.go`](../api-server/handlers/alerts.go#L64-L82) | Validate threshold, update config |
| Dashboard fetch logs | [`dashboard/app.js`](../dashboard/app.js#L42-L57) | Gọi API và render bảng |
| Dashboard stats | [`dashboard/app.js`](../dashboard/app.js#L62-L78) | Count theo level |
| Dashboard pagination | [`dashboard/app.js`](../dashboard/app.js#L103-L123) | Không load toàn bộ log |
| Dashboard WebSocket | [`dashboard/app.js`](../dashboard/app.js#L173-L212) | Nhận alert realtime, reconnect |
| Dashboard threshold | [`dashboard/app.js`](../dashboard/app.js#L235-L255) | POST dynamic threshold |
| API contract | [`docs/api.md`](api.md) | Request/response chính thức |
| Kiến trúc | [`docs/architecture.md`](architecture.md) | Câu chuyện hệ thống tổng thể |
| Evidence | [`docs/testing-evidence.md`](testing-evidence.md) | Nơi ghi output thật |

## Tuần 1: Đóng băng MVP và verify end-to-end

### Mục tiêu

Kết thúc tuần 1, bạn phải trả lời được: “Hệ thống có chạy từ đầu đến cuối không?” bằng bằng chứng thật, không chỉ bằng mô tả.

### Việc cần làm

- Chạy `docker compose config` để xác nhận Compose file hợp lệ.
- Khởi động hệ thống bằng `docker compose up -d`.
- Kiểm tra Elasticsearch health.
- Kiểm tra demo services có ghi log vào `./logs/demo-node/app.log` và `./logs/demo-go/app.log`.
- Kiểm tra Filebeat có tail file.
- Kiểm tra Logstash có nhận event và index vào Elasticsearch.
- Kiểm tra `logs-*` có document.
- Kiểm tra API health, logs list, filter, count.
- Mở dashboard tại `http://localhost:8080`.
- Hạ threshold alert xuống thấp và xác nhận alert có thể bắn.
- Ghi toàn bộ kết quả quan trọng vào [`docs/testing-evidence.md`](testing-evidence.md).

### Command checklist

```bash
docker compose config
docker compose up -d
docker compose ps
```

```bash
curl -u elastic:${ES_PASSWORD:-changeme123} http://localhost:9200/_cluster/health
curl -u elastic:${ES_PASSWORD:-changeme123} "http://localhost:9200/_cat/indices?v"
curl -u elastic:${ES_PASSWORD:-changeme123} "http://localhost:9200/logs-*/_count"
```

```bash
curl http://localhost:8080/api/health
curl "http://localhost:8080/api/logs?size=3"
curl "http://localhost:8080/api/logs?level=ERROR&size=3"
curl "http://localhost:8080/api/logs?app=demo-node&size=3"
curl "http://localhost:8080/api/logs/count"
```

```bash
curl -s -X POST http://localhost:8080/api/alerts/config \
  -H "Content-Type: application/json" \
  -d '{"threshold":5}'

./scripts/trigger-error-spike.sh 20
```

```bash
docker compose logs api-server
docker compose logs filebeat
docker compose logs logstash
```

### Code cần đọc trong tuần 1

- Docker service dependency: [`docker-compose.yml`](../docker-compose.yml#L13-L168)
- Log producer Node.js: [`services/demo-node/index.js`](../services/demo-node/index.js#L77-L106)
- Log producer Go: [`services/demo-go/main.go`](../services/demo-go/main.go#L107-L143)
- Filebeat input/output: [`filebeat/filebeat.yml`](../filebeat/filebeat.yml#L9-L77)
- Logstash parse/enrich/output: [`logstash/pipeline/logstash.conf`](../logstash/pipeline/logstash.conf#L12-L115)
- API routes: [`api-server/main.go`](../api-server/main.go#L71-L94)

### Kiến thức cần nắm

- Container trong cùng Docker network gọi nhau bằng **service name**, không phải `localhost`.
- Filebeat không parse business log chính; nó tail file và ship tới Logstash.
- Logstash parse JSON Lines từ field `message`, sau đó promote field lên root.
- Elasticsearch index theo ngày `logs-YYYY.MM.dd`.
- Dashboard không gọi trực tiếp Elasticsearch, mà gọi Go API.

### Output cuối tuần

- [`docs/testing-evidence.md`](testing-evidence.md) có ít nhất 8 mục evidence thật.
- Có screenshot dashboard sau khi hiển thị log.
- Có screenshot hoặc log output chứng minh alert đã bắn.
- Có ghi lại lỗi gặp phải nếu có, ví dụ ES chưa start, Filebeat chưa tail, API chưa connect được ES.

### Tiêu chí pass/fail

| Tiêu chí | Pass khi |
|---|---|
| Pipeline | `logs-*/_count` trả `count > 0` |
| API | `/api/logs` trả `data`, `total`, `page`, `size` |
| Filter | `level=ERROR` chỉ trả log ERROR |
| Dashboard | Mở được `http://localhost:8080` và thấy log |
| Alert | Có WebSocket/banner hoặc log `[alerting] alert sent` |

## Tuần 2: Học sâu implementation và sửa drift tài liệu

### Mục tiêu

Kết thúc tuần 2, bạn phải giải thích được từng tầng bằng lời của mình, chỉ được mở code để đối chiếu, không cần đọc thuộc lòng.

### Luồng học đề xuất

1. **Demo services**: Vì sao app chỉ ghi file log?
2. **Filebeat**: Vì sao cần agent tail file?
3. **Logstash**: Vì sao parse JSON trước, Grok chỉ enrich phụ?
4. **Elasticsearch**: Vì sao dùng `term`, `match`, `range`?
5. **Go API**: Vì sao không để dashboard query ES trực tiếp?
6. **Alerting**: Vì sao dùng sliding window, cooldown, dedup?
7. **Dashboard**: Vì sao pagination 20 record/trang?

### Việc cần làm

- Đọc lại [`docs/knowledge-base.md`](knowledge-base.md), nhưng đối chiếu với code thật.
- Sửa hoặc ghi chú những đoạn tài liệu đang lệch implementation.
- Đặc biệt kiểm tra phần Alert Deduplication: code hiện tại dùng single `Lock` trong [`shouldAlert`](../api-server/alerting/engine.go#L141-L153), không dùng pattern `RLock` rồi `Lock`.
- Viết thêm câu trả lời ngắn cho 10 câu hỏi bảo vệ quan trọng vào [`docs/report-notes.md`](report-notes.md).
- Ghi các quyết định kỹ thuật vào [`docs/decisions.md`](decisions.md).

### Code cần đọc kỹ

- Query builder: [`api-server/handlers/logs.go`](../api-server/handlers/logs.go#L164-L205)
- Extract ES hits: [`api-server/handlers/logs.go`](../api-server/handlers/logs.go#L210-L255)
- Alert state: [`api-server/alerting/engine.go`](../api-server/alerting/engine.go#L45-L61)
- Sliding window loop: [`api-server/alerting/engine.go`](../api-server/alerting/engine.go#L82-L136)
- Dedup atomic check/write: [`api-server/alerting/engine.go`](../api-server/alerting/engine.go#L141-L153)
- Dynamic threshold: [`api-server/alerting/engine.go`](../api-server/alerting/engine.go#L158-L187)
- Dashboard WebSocket: [`dashboard/app.js`](../dashboard/app.js#L173-L212)

### Câu hỏi phải tự trả lời được

1. Vì sao chọn JSON Lines thay vì log text thô?
2. Vì sao Logstash vẫn có Grok nếu log đã là JSON?
3. Vì sao dùng Elasticsearch thay vì MySQL/PostgreSQL?
4. `term` khác `match` như thế nào trong Elasticsearch?
5. Sliding window khác polling mỗi 5 phút như thế nào?
6. Alert fatigue là gì và dedup giải quyết ra sao?
7. Vì sao `shouldAlert` cần check và ghi `sent[key]` trong cùng một lock?
8. Vì sao dashboard cần pagination?
9. Vì sao Go phù hợp với alerting goroutine?
10. Nếu Logstash chết 1 phút rồi sống lại, hệ thống xử lý log như thế nào?

### Output cuối tuần

- [`docs/knowledge-base.md`](knowledge-base.md) không còn mô tả sai implementation chính.
- [`docs/decisions.md`](decisions.md) có ít nhất 5 quyết định kỹ thuật:
  - Go thay Spring Boot.
  - Elasticsearch thay SQL.
  - Filebeat + Logstash thay app ship trực tiếp.
  - JSON Lines làm format chính.
  - Sliding window + dedup cho alerting.
- [`docs/report-notes.md`](report-notes.md) có phần Q&A dùng được khi luyện bảo vệ.

## Tuần 3: Thêm điểm cá nhân nhỏ - chủ động trigger incident

### Mục tiêu

Thêm một lát cắt nhỏ để demo alerting có thể tái lập. Điểm này giúp dự án có giá trị cá nhân mà không phá scope MVP.

### Tên gợi ý

**Incident Replay / Controlled Error Spike**

### Vấn đề giải quyết

Hiện tại demo services sinh ERROR theo xác suất. Khi bảo vệ, nếu chờ random sinh đủ ERROR thì demo alert có thể thiếu ổn định. Một kịch bản incident chủ động giúp bạn nói:

> Em bổ sung cơ chế tái hiện spike lỗi để kiểm thử alerting. Nhờ vậy demo không phụ thuộc may rủi, và hệ thống chứng minh được khả năng phát hiện incident trong điều kiện có thể lặp lại.

### Scope nên làm

Chọn một trong hai hướng, ưu tiên hướng A nếu muốn ít rủi ro.

#### Hướng A: Script ghi spike log vào file hiện có

Tạo script nhỏ ghi nhiều dòng ERROR JSON Lines vào `./logs/demo-node/app.log` hoặc `./logs/demo-go/app.log`.

Ưu điểm:

- Không cần sửa service đang chạy.
- Không ảnh hưởng format log.
- Dễ giải thích: script giả lập incident từ một service.
- Dễ rollback nếu không ổn.

Điều kiện bắt buộc:

- Dòng log vẫn giữ format JSON Lines giống demo services.
- Field bắt buộc: `timestamp`, `level`, `service`, `message`, `metadata`.
- Không ghi sai path Filebeat đang tail.

Code cần đối chiếu format:

- Node format: [`services/demo-node/index.js`](../services/demo-node/index.js#L81-L93)
- Go format: [`services/demo-go/main.go`](../services/demo-go/main.go#L114-L131)
- Filebeat paths: [`filebeat/filebeat.yml`](../filebeat/filebeat.yml#L14-L18)

#### Hướng B: Thêm endpoint demo-only để trigger incident

Thêm endpoint như `POST /api/demo/error-spike` để API ghi hoặc yêu cầu demo service tạo spike.

Ưu điểm:

- Demo đẹp hơn vì gọi API được.

Rủi ro:

- Mở rộng API surface.
- Cần test kỹ hơn.
- Cần giải thích vì sao endpoint này chỉ phục vụ demo/test, không phải production.

Nếu chọn hướng B, phải cập nhật [`docs/api.md`](api.md), [`docs/testing-evidence.md`](testing-evidence.md), và có kiểm tra route trong [`api-server/main.go`](../api-server/main.go#L71-L86).

### Khuyến nghị

Chọn **hướng A** cho dự án tốt nghiệp 1 tháng. Nó đủ tạo dấu ấn cá nhân, nhưng ít làm hệ thống phức tạp.

### Acceptance criteria

- Chạy một command/script có thể tạo ít nhất 10-20 ERROR log trong thời gian ngắn.
- Filebeat ship được các log đó.
- Elasticsearch count ERROR tăng.
- Alerting engine phát alert sau tối đa `ALERT_CHECK_INTERVAL_SECONDS` cộng thêm thời gian ship/index.
- Dashboard nhận banner hoặc API server log có dòng `[alerting] alert sent`.
- Evidence được ghi vào [`docs/testing-evidence.md`](testing-evidence.md).

### Command verify sau khi implement

```bash
curl -s -X POST http://localhost:8080/api/alerts/config \
  -H "Content-Type: application/json" \
  -d '{"threshold":5}'

./scripts/trigger-error-spike.sh 20
```

```bash
curl -u elastic:${ES_PASSWORD:-changeme123} "http://localhost:9200/logs-*/_count"
curl "http://localhost:8080/api/logs?level=ERROR&size=5"
docker compose logs api-server
```

### Output cuối tuần

- Có một cách trigger incident chủ động.
- Có tài liệu ngắn giải thích kịch bản incident trong [`docs/report-notes.md`](report-notes.md).
- Có evidence thật trong [`docs/testing-evidence.md`](testing-evidence.md).
- Demo alert không còn phụ thuộc random ERROR.

## Tuần 4: Đóng gói báo cáo, slide và rehearsal

### Mục tiêu

Kết thúc tuần 4, dự án phải sẵn sàng cho người khác clone, chạy, xem dashboard và hiểu câu chuyện kỹ thuật.

### Việc cần làm

- Test lại flow clone sạch.
- Chụp screenshot dashboard ở trạng thái bình thường.
- Chụp screenshot khi filter ERROR.
- Chụp screenshot alert banner.
- Đo response time thật.
- Cập nhật README nếu hướng dẫn chưa đủ.
- Chuẩn bị demo script 5 phút.
- Chuẩn bị Q&A 10 phút.
- Tập nói ít nhất 3 lần: một lần nhìn notes, một lần chỉ nhìn slide, một lần giả lập bị hỏi ngang.

### Clean clone checklist

```bash
git clone git@github.com:SonBestCodeVien5/log-system.git fresh-test
cd fresh-test
# Tùy chọn: chỉ cần tạo .env nếu muốn override default trong docker-compose.yml
# cp .env.example .env
sudo sysctl -w vm.max_map_count=262144
docker compose up -d
docker compose ps
curl http://localhost:8080/api/health
```

Pass khi:

- Dashboard mở được tại `http://localhost:8080`.
- Có log trong bảng.
- Filter hoạt động.
- Alert trigger được bằng kịch bản incident.

### Demo script 5 phút

1. **30 giây**: giới thiệu bài toán centralized logging.
2. **45 giây**: chỉ luồng dữ liệu service -> Filebeat -> Logstash -> Elasticsearch -> API -> Dashboard.
3. **60 giây**: mở dashboard, filter `ERROR`, filter `demo-node`, giải thích API query.
4. **75 giây**: trigger incident, hạ threshold, chờ alert banner.
5. **60 giây**: mở [`api-server/alerting/engine.go`](../api-server/alerting/engine.go#L82-L153) giải thích sliding window và dedup.
6. **30 giây**: kết luận: hệ thống đã có log tập trung, search/filter, alert realtime, evidence và hướng production tiếp theo.

### Slide/report nên có

1. Bối cảnh bài toán.
2. Kiến trúc tổng thể.
3. Log format JSON Lines.
4. Pipeline Filebeat + Logstash + Elasticsearch.
5. Go API và dashboard.
6. Alerting engine: sliding window, dedup, dynamic threshold.
7. Demo incident replay.
8. Kết quả kiểm thử và response time.
9. Hạn chế.
10. Hướng phát triển.

### Output cuối tháng

- README chạy được theo flow clone -> compose up; `.env` chỉ cần khi muốn override default.
- [`docs/testing-evidence.md`](testing-evidence.md) có output thật.
- [`docs/report-notes.md`](report-notes.md) có demo script và Q&A.
- [`docs/decisions.md`](decisions.md) có quyết định kỹ thuật.
- Có ít nhất 3 screenshot dùng cho báo cáo/slide.
- Có một kịch bản incident chủ động.

## Kế hoạch theo ngày

| Ngày | Việc chính | Output |
|---:|---|---|
| 1 | Chạy `docker compose config`, start stack | Stack chạy hoặc ghi blocker |
| 2 | Verify ES, Filebeat, Logstash, `logs-*` | Evidence pipeline |
| 3 | Verify API logs/filter/count | Evidence API |
| 4 | Verify dashboard và WebSocket | Screenshot dashboard |
| 5 | Verify alert bằng threshold thấp | Evidence alert |
| 6 | Đo response time | Số liệu thật |
| 7 | Tổng hợp lỗi tuần 1 | Danh sách fix/docs |
| 8 | Đọc demo services + Filebeat | Ghi giải thích source -> collector |
| 9 | Đọc Logstash + ES mapping/query | Ghi giải thích parse/index |
| 10 | Đọc `logs.go` | Ghi giải thích API query |
| 11 | Đọc `engine.go` | Ghi giải thích alerting |
| 12 | Đọc dashboard JS | Ghi giải thích UI/API/WS |
| 13 | Sửa drift `knowledge-base.md` | Docs khớp code |
| 14 | Ghi `decisions.md` | 5 quyết định kỹ thuật |
| 15 | Thiết kế incident replay | Plan nhỏ, không đổi format log |
| 16 | Implement hoặc chuẩn bị script trigger | Script/flow tạo ERROR spike |
| 17 | Verify incident replay | Alert trigger ổn định |
| 18 | Ghi evidence incident | Output + screenshot |
| 19 | Viết runbook/alert explanation | Report notes |
| 20 | Review scope, dọn docs | Không còn TODO lớn |
| 21 | Buffer fix bug nhỏ | Stack ổn định |
| 22 | Test clean clone | Ghi thời gian chạy |
| 23 | Chuẩn bị slide kiến trúc | Draft slide |
| 24 | Chuẩn bị slide demo/evidence | Screenshot + số liệu |
| 25 | Chuẩn bị Q&A | Câu trả lời ngắn |
| 26 | Rehearsal lần 1 | Ghi điểm vấp |
| 27 | Rehearsal lần 2 | Demo dưới 5 phút |
| 28 | Rehearsal lần 3 | Tập bị hỏi ngang |
| 29 | Final docs pass | README/docs sạch |
| 30 | Không thêm code mới | Chỉ kiểm tra và nghỉ nhịp |

## Danh sách không nên làm trong 1 tháng cuối

Không nên thêm các scope sau nếu chưa hoàn thành verify/evidence:

- Authentication/RBAC.
- Kubernetes deployment.
- Distributed tracing.
- AI log analysis.
- Multi-tenant dashboard.
- Đổi dashboard sang framework mới.
- Đổi log format.
- Đổi Elasticsearch sang database khác.
- Refactor lớn API/server.

Các scope này không sai, nhưng dễ tạo thêm bug và khó bảo vệ sâu trong thời gian ngắn.

## Cách kể giá trị cá nhân khi bảo vệ

Bạn có thể đóng khung dự án như sau:

> Dự án không chỉ gom log và hiển thị dashboard. Em tập trung vào luồng vận hành hoàn chỉnh: log từ nhiều service được thu thập, parse, index, query và cảnh báo realtime. Điểm em bổ sung để phục vụ kiểm thử và demo là kịch bản incident replay, giúp tái hiện error spike một cách chủ động thay vì phụ thuộc log random.

Các ý cần nhấn:

- Biết chọn scope vừa đủ cho dự án 1 tháng.
- Biết dùng công cụ đúng domain: Filebeat, Logstash, Elasticsearch.
- Biết chứng minh bằng evidence, không chỉ nói tính năng.
- Biết kiểm soát alert fatigue bằng dedup/cooldown.
- Biết tạo demo incident có thể lặp lại.

## Checklist cuối cùng trước ngày bảo vệ

- [ ] `docker compose up -d` chạy sạch.
- [ ] `curl http://localhost:8080/api/health` trả OK.
- [ ] Dashboard mở được.
- [ ] Log table có data.
- [ ] Filter `ERROR` hoạt động.
- [ ] Filter `demo-node` hoạt động.
- [ ] `/api/logs/count` trả count theo level.
- [ ] Incident replay trigger được ERROR spike.
- [ ] Alert banner xuất hiện hoặc API log ghi alert sent.
- [ ] `docs/testing-evidence.md` có output thật.
- [ ] Slide có kiến trúc, demo, evidence, hạn chế.
- [ ] Bạn giải thích được 10 câu hỏi ở tuần 2.
- [ ] Không còn thêm feature mới trong 48 giờ cuối.
