# Task: Thêm alert rule mới

## Yêu cầu
Thêm rule: [MÔ TẢ RULE — ví dụ: cảnh báo khi 1 service cụ thể có quá nhiều WARN]

## Cần sửa
- `api-server/alerting/engine.go` — thêm logic check mới
- `dashboard/app.js` — hiển thị loại alert mới
- `dashboard/index.html` — thêm UI element nếu cần

## Lưu ý
- Dùng RWMutex khi đọc/ghi shared state
- Deduplication key phải unique cho từng loại rule
- Test bằng cách tăng tốc sinh ERROR trong demo service
