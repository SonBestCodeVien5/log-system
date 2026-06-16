# AGENTS.md — services

Thư mục này chứa 2 demo service sinh log để test pipeline.

## Mục đích
Không phải nghiệp vụ thật — chỉ để sinh log liên tục với đủ 3 level
(INFO, WARN, ERROR) theo đúng format JSON Lines mà Logstash có thể parse.

## Format log BẮT BUỘC
```json
{"timestamp":"2024-01-15T10:23:11Z","level":"ERROR","service":"demo-node","message":"Payment gateway timeout","metadata":{"order_id":"789"}}
```
Mỗi dòng là một JSON object hoàn chỉnh. `timestamp` phải là ISO8601 UTC,
`level` viết HOA (`INFO`, `WARN`, `ERROR`), `service` không có khoảng trắng,
`message` là nội dung hiển thị trên dashboard, và `metadata` chứa dữ liệu phụ.

## demo-node (Node.js)
- Ghi log ra file `/logs/demo-node/app.log`
- Sinh ngẫu nhiên INFO (60%), WARN (25%), ERROR (15%)
- Mỗi 1-3 giây sinh 1 log
- Có thể tăng tốc sinh ERROR để test alerting

## demo-go (Go)
- Ghi log ra file `/logs/demo-go/app.log`
- Logic tương tự demo-node
- Module: `github.com/SonBestCodeVien5/log-system/demo-go`
