# Review

## Findings
- Severity: Resolved Medium.
- File/line: `dashboard/index.html:7`, `dashboard/index.html:142`, `api-server/main.go:90-96`.
- Issue: Dashboard HTML/CSS/JS từng dùng URL cố định và static responses không có cache policy. Browser có warmed cache chạy HTML/JS từ 2026-06-17 dù host/container đã bind-mount file mới. JS cũ không đọc time inputs nên nút Tìm gửi `/api/logs?page=1&size=20` và count chỉ có `from=now-1h`, đúng hiện tượng “vẫn full log”. Restart WSL/Docker không xóa browser HTTP cache.
- Evidence: Host và container có cùng hash JS mới `f0476e77...`; fetch `cache:no-store` nhận JS mới dài 13.654 byte. Normal navigation trên `http://localhost:8080` nhận cached JS cũ dài 11.838 byte và HTML cũ dài 4.225 byte. Khi cache-bust riêng `/`, HTML mới hiện hai ô thời gian nhưng JS cũ vẫn chạy, inputs trống và network vẫn thiếu `from/to`.
- Resolution: `index.html` đã version hóa `style.css`/`app.js`; API server thêm `DashboardNoCache` cho `/` và `/assets/*` với unit test. API container đã rebuild healthy. Warmed-cache browser phục hồi qua `/?v=time-range-20260622`, nhận no-store headers và exact filter `06:00-06:10` gửi đúng `from/to`, trả 0 log.
- Finding Medium trước đó đã được xử lý tại `dashboard/app.js:253-256` và `dashboard/app.js:287-293`: Apply, timer và WebSocket `error_spike` đều đi qua `refreshDashboard()`, tiến sliding range đúng một lần rồi tải cả logs/count. Browser trace xác nhận hai endpoint nhận cùng range mới; custom range vẫn cố định.
- Severity: Low, UX ambiguity (không phải lỗi range query).
- File/line: `dashboard/index.html:123-139`, `dashboard/app.js:152-160`.
- Issue: Bảng luôn hiển thị tối đa 20 dòng/trang và chưa có nhãn tóm tắt range đang active. Với tốc độ demo hiện tại, một phút vẫn có khoảng 60-70 log nên trang đầu tiếp tục đầy 20 dòng, tạo cảm giác filter không hoạt động dù `total`, số trang và timestamp đã đổi.
- Evidence: Live API trả 3.059 log cho một giờ, 59 log cho một phút và 0 log cho khoảng năm 2099. Trên UI, range local `18:34-18:35` trả 64 log, `Trang 1 / 4`, đúng 20 dòng trên trang đầu và mọi timestamp nằm trong range; range năm 2099 hiển thị 0 log và bảng rỗng.
- Evidence bổ sung theo báo cáo người dùng: range local hôm nay `2026-06-22 06:00-06:10` được đổi đúng thành `2026-06-21T23:00:00Z-23:10:00Z`; API list/count đều HTTP 200 và total 0. Tái hiện trên UI với đúng hai giá trị này cũng cho stats 0, `Trang 1 / 1`, `0 log` và `Không tìm thấy log nào`.
- Recommendation: Nếu muốn tránh hiểu nhầm, thêm dòng/chip `Đang lọc: <from> - <to> · <total> log` cạnh pagination hoặc filter bar. Không cần thay API/Elasticsearch.
- Severity: Medium, pre-existing/outside time-range diff.
- File/line: `dashboard/app.js:69-76`, `dashboard/app.js:95-109`.
- Issue: `fetchLogs()` và `fetchStats()` parse JSON nhưng không kiểm tra `res.ok`; riêng stats còn nuốt toàn bộ exception. Khi API/Elasticsearch trả 500, dashboard có thể hiển thị bảng rỗng hoặc các stat bằng 0 như một response thành công thay vì báo backend lỗi. `git show HEAD:dashboard/app.js` xác nhận hành vi này tồn tại trước feature time-range.
- Recommendation: Tạo follow-up dashboard hardening: kiểm tra `res.ok`, đọc `error` từ response, hiển thị trạng thái lỗi riêng cho table/stats và không ghi đè dữ liệu hợp lệ bằng số 0 khi request thất bại. Finding này không chặn commit time-range nếu staging giữ đúng scope.
- Fixes applied during review: Không sửa source trong lượt re-review này.

## Test Gaps
- Gap: Chưa có test tự động cho state transition `sliding -> custom -> reset -> sliding` hoặc kiểm tra query pair logs/count.
- Risk: Thấp cho MVP vì browser E2E đã bao phủ load, Apply, Reset, timer, invalid range và WebSocket; về lâu dài event path mới vẫn có thể bỏ qua helper chung.
- Gap: Chưa kiểm tra cross-browser và timezone có DST cho `datetime-local`.
- Risk: Thấp trong môi trường hiện tại dùng Asia/Bangkok, nhưng nên kiểm tra trước khi triển khai cho người dùng đa timezone.
- Gap: Mobile/responsive chưa kiểm tra.
- Risk: Được người dùng chủ động hoãn; không chặn scope desktop.

## Residual Risks
- Risk: Các request async không có abort/sequence guard; nếu timer, Apply và WebSocket xảy ra gần nhau trên backend chậm, response cũ có thể về sau response mới. Đây là rủi ro Low đã tồn tại trong cơ chế auto-refresh và chưa quan sát thấy ở E2E hiện tại.
- Risk: HTTP 4xx/5xx đang bị dashboard xử lý như dữ liệu rỗng/0; đây là follow-up pre-existing cần ưu tiên sau feature.
- Risk: Console còn 404 `/favicon.ico`, không liên quan feature.
- Risk: Client/tab cũ có thể giữ bundle dashboard cũ vì chưa có cache-busting/hot reload.
- Risk: Working tree có nhiều thay đổi docs/context không thuộc feature; Git handoff phải stage theo file, không stage toàn bộ repo.

## Commit Readiness
- Ready: Có, cho scope desktop đã thống nhất.
- Reason: Logic/API time range, sliding/custom behavior, cache-busting/no-cache middleware và warmed-cache recovery đều đã pass unit/static/runtime/browser E2E. Client có cached root từ trước chỉ cần one-time recovery URL hoặc clear site data; response mới không còn cacheable.

## Next Handoff
- Current phase: review
- Next phase: log-git
- Must read: `agent-context/features/dashboard-time-range-filter/03-implementation.md`, `agent-context/features/dashboard-time-range-filter/04-verification.md`, `dashboard/index.html`, `dashboard/app.js`, `dashboard/style.css`, `api-server/main.go`
- Decisions locked: Default range trượt một giờ; custom range cố định; mỗi refresh cycle chính dùng cùng cặp `from`/`to`; mobile ngoài scope hiện tại.
- Open risks: Existing cached root cần one-time recovery; UX chưa có active-range summary; pre-existing HTTP error handling Medium; async response ordering Low; chưa có automated UI test; mobile/DST và favicon ngoài scope.
- Validation status: Cache remediation verified: Go tests, JS syntax, Compose config, healthy rebuilt API, no-store headers and warmed-cache browser exact filter pass. Ready for scoped Git handoff.
