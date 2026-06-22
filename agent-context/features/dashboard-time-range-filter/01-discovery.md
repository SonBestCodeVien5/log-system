# Discovery

## Request
- Feature/request: Xác định vì sao API `/api/logs` có query `from` và `to` nhưng dashboard không hiển thị bộ lọc thời gian, đồng thời lập kế hoạch bổ sung UI.
- Feature slug: `dashboard-time-range-filter`

## Repo Facts
- Current state: API đã đọc `from`/`to` cho cả `/api/logs` và `/api/logs/count`, mặc định `now-1h` đến `now`, rồi dùng chúng làm range `@timestamp` trong Elasticsearch.
- Current state: `dashboard/index.html` chỉ có level, service và search; không có input hoặc preset thời gian.
- Current state: `dashboard/app.js::buildParams()` chỉ thêm `level`, `app`, `q`; request bảng log không gửi `from`/`to` nên hoàn toàn dựa vào mặc định backend.
- Current state: `dashboard/app.js::fetchStats()` hardcode `from=now-1h`, không gửi `to`, nên stats cũng dựa vào `to=now` mặc định nhưng không nhận được lựa chọn từ UI.
- Relevant files: `dashboard/index.html`, `dashboard/app.js`, `dashboard/style.css`, `api-server/handlers/logs.go`, `docs/api.md`.
- Existing constraints: Dashboard phải giữ HTML/Vanilla JS/CSS, pagination 20 dòng, dùng `fetch()`, và logic/structure/style nằm đúng ba file dashboard.

## Applicable Instructions
- Root `AGENTS.md`: Dashboard không dùng framework; không hardcode password/port; pagination 20 record/trang.
- Area `AGENTS.md`: Filter bar hiện yêu cầu level/app/search; dashboard phải giữ auto-refresh và responsive behavior.
- Skill references: `phase-discovery.md` yêu cầu xác nhận UI/API contract; `phase-dashboard.md` cho phép time filter khi phù hợp và yêu cầu UI vận hành tốt cho việc kiểm tra log lặp lại.

## Unknowns And Risks
- Unknowns: Chưa có yêu cầu riêng về timezone hoặc preset thời gian; kế hoạch chọn `datetime-local` và chuyển sang ISO8601 UTC trước khi gọi API.
- Risks: Nếu chỉ nối time range vào `/api/logs` mà không nối `/api/logs/count`, số liệu stats sẽ không khớp bảng.
- Risks: `datetime-local` không chứa timezone; phải chuyển bằng `Date#toISOString()` và kiểm tra ngày bắt đầu không lớn hơn ngày kết thúc.
- Risks: Auto-refresh với khoảng kết thúc cố định sẽ refresh dữ liệu nhưng không dịch chuyển cửa sổ thời gian; đây là hành vi đúng cho custom absolute range và cần thể hiện rõ trong acceptance criteria.

## Next Handoff
- Current phase: discovery
- Next phase: plan
- Must read: `agent-context/features/dashboard-time-range-filter/02-plan.md`, `dashboard/index.html`, `dashboard/app.js`, `dashboard/style.css`
- Decisions locked: Nguyên nhân nằm ở dashboard chưa có control và chưa truyền `from`/`to`, không phải API thiếu hỗ trợ.
- Open risks: Timezone conversion và đồng bộ stats/table.
- Validation status: Đã kiểm tra tĩnh source và API documentation; chưa chạy browser/E2E vì đây là phase planning.
