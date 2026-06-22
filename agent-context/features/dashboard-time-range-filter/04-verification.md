# Verification

## Commands
- Command: `GOCACHE=/tmp/log-system-go-build-cache go test ./...` trong `api-server`.
- Result: Pass.
- Important output: Bao gồm `middleware/cache_test.go`; dashboard root/assets nhận no-cache headers, API không bị áp header này.
- Command: `node --check dashboard/app.js`
- Result: Pass.
- Important output: Không có lỗi cú pháp JavaScript.
- Command: `git diff --check`
- Result: Pass.
- Important output: Không có whitespace error.
- Command: `docker compose config`
- Result: Pass.
- Important output: Compose render thành công; dashboard được bind mount read-only vào `/app/static` nên source mới được serve trực tiếp.
- Command: `docker compose ps`
- Result: Pass sau khi cấp quyền Docker socket.
- Important output: API server, Elasticsearch, Filebeat, Logstash và hai demo service đều running/healthy.
- Command: `docker compose up -d --build api-server`
- Result: Pass.
- Important output: API image rebuild thành công, container được recreate và health status trở lại `healthy`.

## Scenarios
- Scenario: Load dashboard lần đầu tại `http://127.0.0.1:8080`.
- Expected: Hiển thị Từ/Đến một giờ gần nhất; logs/count nhận cùng range và trả 200.
- Actual: Pass; ví dụ cả hai request dùng `from=2026-06-22T10:10:17.000Z`, `to=2026-06-22T11:10:17.000Z` và trả HTTP 200.
- Scenario: Chọn custom local range `2026-06-22T16:00` đến `17:00` trong timezone Asia/Bangkok.
- Expected: Cả hai endpoint nhận range UTC tương ứng.
- Actual: Pass; logs/count cùng nhận `09:00:00.000Z` đến `10:00:00.000Z`, HTTP 200.
- Scenario: Chọn `from=18:00`, `to=17:00` rồi nhấn Tìm.
- Expected: Hiện lỗi và không gửi request.
- Actual: Pass; UI hiện `Thời gian Từ phải sớm hơn hoặc bằng thời gian Đến.` và danh sách network không có request logs/count mới.
- Scenario: Nhấn Xóa filter sau range sai.
- Expected: Xóa lỗi/filter và phục hồi một giờ gần nhất.
- Actual: Pass; level/app/search rỗng, Từ/Đến được gán lại và request tải lại thành công.
- Scenario: Để auto-refresh bật ở default sliding range.
- Expected: Hai đầu mút tiến theo thời gian và logs/count luôn dùng cùng range.
- Actual: Pass; input tiến từ `17:10:17–18:10:17` lên `17:10:40–18:10:40`; các request tương ứng đều HTTP 200.
- Scenario: Phát một WebSocket `error_spike` trong sliding mode khi auto-refresh đang tắt để cô lập event.
- Expected: Range tiến ngay tại thời điểm alert và cả logs/count dùng cùng cặp mới.
- Actual: Pass; range tiến từ `10:24:40–11:24:40Z` lên `10:24:56–11:24:56Z`; `/api/logs` và `/api/logs/count` đều HTTP 200 với cùng range mới.
- Scenario: Phát WebSocket `error_spike` sau khi người dùng chọn custom range `16:00–17:00` Asia/Bangkok.
- Expected: Refresh cả logs/count nhưng không thay đổi custom range.
- Actual: Pass; hai input giữ nguyên và cả hai endpoint dùng `09:00:00–10:00:00Z`, HTTP 200.
- Scenario: Reset về sliding mode, chờ rồi nhấn Tìm khi auto-refresh tắt.
- Expected: Apply tiến range một lần và logs/count dùng cùng range mới.
- Actual: Pass; reset dùng `10:26:02–11:26:02Z`, Apply tiến lên `10:26:20–11:26:20Z`; hai endpoint đều HTTP 200.
- Scenario: Browser context đã warmed-cache HTML/JS cũ mở URL recovery `/?v=time-range-20260622` sau cache remediation.
- Expected: HTML mới tham chiếu asset versioned, time inputs được khởi tạo và response dashboard có no-store headers.
- Actual: Pass; DOM có Từ/Đến với default một giờ, asset URLs mang `v=time-range-20260622`; `/`, CSS và JS trả `Cache-Control: no-store, no-cache, must-revalidate, max-age=0`, `Pragma: no-cache`, `Expires: 0`.
- Scenario: Trên recovered warmed-cache client, chọn hôm nay `06:00-06:10`, tắt auto-refresh và nhấn Tìm.
- Expected: Request list/count có cùng ISO UTC range và UI không hiển thị full log.
- Actual: Pass; request dùng `2026-06-21T23:00:00Z-23:10:00Z`, cả hai HTTP 200; stats 0, `Trang 1 / 1`, `0 log`, bảng `Không tìm thấy log nào`.

## Failures Or Skips
- Failure/skip: Entry `/` đã nằm trong browser cache trước khi server có no-store headers vẫn có thể được dùng lại một lần.
- Reason: Server không thể retroactively xóa response đã lưu trong cache client; cần mở `/?v=time-range-20260622`, hard reload/clear site data một lần. Asset versioning và no-store ngăn tái diễn sau recovery.
- Failure/skip: Responsive/mobile test.
- Reason: Được hoãn theo yêu cầu người dùng.
- Failure/skip: Go unit tests.
- Reason: Feature chỉ thay đổi static dashboard; API Go contract không đổi.
- Failure/skip: Console có một lỗi 404 `/favicon.ico`.
- Reason: Lỗi asset có sẵn, không liên quan time-range và không ảnh hưởng feature.

## Next Handoff
- Current phase: verification
- Next phase: review
- Must read: `agent-context/features/dashboard-time-range-filter/03-implementation.md`, `dashboard/index.html`, `dashboard/app.js`, `dashboard/style.css`
- Decisions locked: Feature hoạt động trên desktop, cùng range cho logs/count, validation chặn request sai, auto-refresh giữ default sliding window.
- Open risks: Client có cached root từ trước cần one-time recovery URL/clear site data; mobile chưa được kiểm tra/tối ưu; favicon 404 ngoài scope.
- Validation status: Cache remediation pass trên warmed-cache browser qua versioned recovery URL; exact 06:00-06:10 filter gửi `from/to` và trả 0. Sẵn sàng re-review.

## Runtime Cache Diagnosis
- Host hashes: `app.js=f0476e77...`, `index.html=5002a594...`, `style.css=eadd1c7d...`.
- Container hashes: Khớp hoàn toàn host cho cả ba bind-mounted files; chỉ có một Docker listener trên port 8080.
- Fresh/no-store fetch: `/assets/app.js` có hash mới `f0476e77...`, dài 13.654 byte và chứa `params.set("from", timeRange.from)`.
- Normal cached navigation: Browser dùng HTML cũ dài 4.225 byte và JS cũ dài 11.838 byte, last-modified ngày 2026-06-17; DOM cũ không có time inputs và request là `/api/logs?page=1&size=20` cùng `/api/logs/count?from=now-1h`.
- Mixed-cache reproduction: Thêm query vào `/` tải HTML mới nên time inputs xuất hiện, nhưng fixed `/assets/app.js` vẫn lấy JS cũ; inputs không được khởi tạo và nút Tìm tiếp tục không gửi `from/to`.
- Exact API control: `2026-06-22 06:00-06:10` Asia/Bangkok tương ứng `2026-06-21T23:00Z-23:10Z`; list/count đều HTTP 200, total 0.
- Remediation: `index.html` dùng versioned CSS/JS URLs; API server thêm `DashboardNoCache`; browser warmed-cache phục hồi thành công qua `/?v=time-range-20260622` và các response mới không còn cacheable.
