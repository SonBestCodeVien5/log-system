# Plan

## Summary
- Goal: Bổ sung bộ lọc thời gian `from`/`to` trên dashboard và dùng cùng một time-range cho bảng log lẫn stats.
- Success criteria: Người dùng chọn thời gian bắt đầu/kết thúc, nhấn Tìm, và cả `/api/logs` lẫn `/api/logs/count` nhận cùng hai giá trị ISO8601; reset trở về cửa sổ mặc định một giờ gần nhất.

## Key Changes
- Implementation: Thêm hai input `datetime-local` có label `Từ` và `Đến` trong filter bar của `dashboard/index.html`.
- Implementation: Khởi tạo input thành thời điểm hiện tại trừ một giờ và hiện tại; tạo helper đọc, validate và đổi local datetime sang ISO8601.
- Implementation: Mở rộng query builder để thêm `from`/`to`; tái sử dụng cùng helper trong `fetchLogs()` và `fetchStats()` để stats/table không lệch phạm vi.
- Implementation: Khi Apply/Reset, đưa page về 1; Reset phục hồi time range một giờ gần nhất. Auto-refresh giữ nguyên range người dùng đã chọn.
- Implementation: Thêm style desktop cho label/input thời gian trong `dashboard/style.css`.
- Public API/UI/data contracts: Không đổi backend contract; tiếp tục dùng `GET /api/logs?from=<ISO8601>&to=<ISO8601>` và `GET /api/logs/count?from=<ISO8601>&to=<ISO8601>`.
- Public API/UI/data contracts: Nếu `from > to`, không gửi request mới và hiển thị lỗi validation gần filter bar; hai đầu mút được tính inclusive theo API hiện tại (`gte`/`lte`).
- Out of scope: Không thêm date-picker library, preset nâng cao, timezone selector, hay thay đổi Elasticsearch query.

## Acceptance Criteria
- Scenario: Mở dashboard lần đầu.
- Expected result: Hai input hiển thị khoảng một giờ gần nhất; bảng và stats gọi API với cùng `from`/`to` hợp lệ.
- Scenario: Chọn một custom range rồi nhấn Tìm.
- Expected result: Page trở về 1; cả hai endpoint nhận cùng ISO8601 range; bảng và stats phản ánh cùng khoảng thời gian.
- Scenario: Chọn `from` muộn hơn `to`.
- Expected result: UI báo lỗi rõ ràng và không gọi API bằng range sai.
- Scenario: Nhấn Xóa filter.
- Expected result: Level/app/search được xóa, time range trở lại một giờ gần nhất, bảng và stats được tải lại.

## Assumptions
- Assumption: UI nhập thời gian theo timezone trình duyệt; request chuyển sang ISO8601 UTC để tránh mơ hồ với Elasticsearch.
- Assumption: Khoảng mặc định một giờ gần nhất khớp backend hiện tại và stats hiện tại.
- Assumption: Custom absolute range được giữ nguyên khi auto-refresh; muốn trở về cửa sổ gần nhất thì dùng Reset.

## Test Plan
- Static: Kiểm tra mọi ID mới tồn tại và `buildParams`/stats dùng cùng time-range helper.
- Browser: Quan sát Network khi load, Apply và Reset; so sánh `from`/`to` của `/api/logs` và `/api/logs/count`.
- Functional: Dùng một range có dữ liệu, một range không dữ liệu, và một range đảo ngược.
- Layout: Kiểm tra filter bar ở desktop.
- Regression: Xác nhận level/app/search, pagination, auto-refresh, WebSocket và alert config vẫn hoạt động.

## Next Handoff
- Current phase: plan
- Next phase: implementation
- Must read: `agent-context/features/dashboard-time-range-filter/01-discovery.md`, `dashboard/AGENTS.md`, `dashboard/index.html`, `dashboard/app.js`, `dashboard/style.css`
- Decisions locked: Dùng hai `datetime-local`, chuyển sang ISO8601 UTC, đồng bộ range giữa logs và count, reset về một giờ gần nhất.
- Open risks: Responsive/mobile được hoãn theo yêu cầu người dùng.
- Validation status: Plan đã đối chiếu contract hiện tại; chưa sửa hoặc chạy application source.
