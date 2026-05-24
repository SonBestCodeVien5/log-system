# Task: Thêm API endpoint mới

## Yêu cầu
Thêm endpoint: [MÔ TẢ ENDPOINT]

## Checklist
- [ ] Thêm handler trong `api-server/handlers/`
- [ ] Đăng ký route trong `main.go`
- [ ] Đọc config từ `os.Getenv()`, không hardcode
- [ ] Handle error, trả JSON đúng format `{"data":..., "total":...}`
- [ ] Test bằng curl sau khi viết xong
