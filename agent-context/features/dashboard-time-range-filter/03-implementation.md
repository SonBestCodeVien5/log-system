# Implementation

## Summary
- Implemented: Thêm hai input `datetime-local` Từ/Đến với độ chính xác đến giây vào filter bar.
- Implemented: Khởi tạo cửa sổ trượt một giờ gần nhất; auto-refresh cập nhật hai đầu mút khi người dùng chưa tự chỉnh time range.
- Implemented: Chuyển datetime theo timezone trình duyệt sang ISO8601 UTC và truyền cùng `from`/`to` vào `/api/logs` lẫn `/api/logs/count`.
- Implemented: Validate input bắt buộc, ngày hợp lệ và `from <= to`; range sai hiển thị lỗi và không gửi request.
- Implemented: Reset xóa level/app/search và phục hồi cửa sổ trượt một giờ.
- Implemented: Sau review, gom các refresh cycle vào `refreshDashboard()` để tiến sliding range đúng một lần rồi tải đồng bộ logs/count cho Apply, timer và WebSocket `error_spike`.
- Implemented: Version hóa URL `app.js`/`style.css` và thêm `DashboardNoCache` middleware cho `/` cùng `/assets/*`, ngăn browser tiếp tục chạy bundle cũ khi dashboard bind mount thay đổi.
- Implemented: Thêm unit test xác nhận cache headers chỉ áp dụng cho dashboard, không áp lên API responses.
- Changed files: `dashboard/index.html`, `dashboard/app.js`, `dashboard/style.css`, `api-server/main.go`, `api-server/middleware/cache.go`, `api-server/middleware/cache_test.go`.

## Deviations From Plan
- Deviation: Không bổ sung hoặc kiểm tra layout mobile.
- Reason: Người dùng yêu cầu tạm hoãn mobile.
- Deviation: Cửa sổ mặc định tiếp tục trượt theo auto-refresh; custom range mới giữ cố định.
- Reason: Nếu khóa `to` ngay lúc load thì auto-refresh không thể hiển thị log mới.
- Deviation: Input và formatter giữ cả giây thay vì chỉ phút.
- Reason: Giá trị `to` ở đầu phút sẽ làm log mới trễ tối đa 59 giây.

## Notes For Verification
- Behavior to verify: Initial logs/count request dùng cùng range một giờ; custom range được đổi đúng từ Asia/Bangkok/local time sang UTC; invalid range không gọi API; reset và auto-refresh phục hồi sliding range.
- Known limitations: Chưa tối ưu layout input thời gian cho mobile theo yêu cầu; browser vẫn có lỗi 404 `favicon.ico` tồn tại độc lập với feature.
- Behavior to verify: WebSocket alert ở sliding mode phải tạo một cặp logs/count với range mới; ở custom mode phải giữ nguyên range.
- Behavior to verify: Browser warmed-cache phải phục hồi qua URL versioned, tải asset mới với `no-store`, rồi gửi `from/to`; API responses không nhận cache headers dashboard.

## Next Handoff
- Current phase: implementation
- Next phase: verification
- Must read: `agent-context/features/dashboard-time-range-filter/04-verification.md`, `dashboard/index.html`, `dashboard/app.js`, `dashboard/style.css`
- Decisions locked: Cùng một absolute ISO8601 range cho logs/count; default sliding một giờ; custom range cố định; validation ở client.
- Open risks: Entry `/` đã cache từ trước không thể bị server chủ động xóa; client bị ảnh hưởng cần mở URL recovery hoặc clear site data một lần. Responsive/mobile chưa nằm trong scope hiện tại; cần re-review trước Git handoff.
- Validation status: Unit tests, rebuild container, cache headers và warmed-cache browser E2E đã pass; chi tiết trong `04-verification.md`.
