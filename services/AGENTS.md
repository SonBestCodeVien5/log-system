# AGENTS.md — services

Thư mục này chứa 2 demo service sinh log để test pipeline.

## Mục đích
Không phải nghiệp vụ thật — chỉ để sinh log liên tục với đủ 3 level
(INFO, WARN, ERROR) theo đúng format Grok có thể parse.

## Format log BẮT BUỘC
```
[2024-01-15T10:23:11Z] [ERROR] [tên-service] nội dung message
```
Timestamp phải là ISO8601 UTC. Level viết HOA. Tên service không có khoảng trắng.

## demo-node (Node.js)
- Ghi log ra file `/logs/demo-node/app.log`
- Sinh ngẫu nhiên INFO (60%), WARN (25%), ERROR (15%)
- Mỗi 1-3 giây sinh 1 log
- Có thể tăng tốc sinh ERROR để test alerting

## demo-go (Go)
- Ghi log ra file `/logs/demo-go/app.log`
- Logic tương tự demo-node
- Module: `github.com/SonBestCodeVien5/log-system/demo-go`
