# AGENTS.md — dashboard

## Stack
HTML + Vanilla JS + CSS thuần. Không dùng bất kỳ framework hay thư viện nào.

## Cấu trúc
- `index.html` — layout, không có logic
- `app.js`     — toàn bộ logic: fetch API, filter, WebSocket, render
- `style.css`  — styling

## API Base URL
```js
const API_BASE = 'http://localhost:8080'
const WS_URL   = 'ws://localhost:8080/ws/alerts'
```

## Tính năng cần có
1. Bảng danh sách log — phân trang, 20 dòng/trang
2. Filter bar — chọn level (ALL/INFO/WARN/ERROR), chọn app, search text
3. Alert banner — hiện banner đỏ + số lượng khi nhận WebSocket alert, tắt được
4. Auto-refresh — tự động fetch log mới mỗi 10 giây
5. Threshold control — input chỉnh ngưỡng alert, gửi POST /api/alerts/config

## Màu sắc log level
- INFO  → #2563eb (xanh)
- WARN  → #d97706 (vàng)
- ERROR → #dc2626 (đỏ)
