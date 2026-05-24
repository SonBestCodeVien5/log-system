# Task: Debug khi log không vào Elasticsearch

## Checklist kiểm tra theo thứ tự

1. ES có đang chạy không?
   `curl http://localhost:9200/_cluster/health`

2. Logstash có nhận được gì từ Filebeat không?
   `docker compose logs filebeat`
   `docker compose logs logstash`

3. Filebeat có đang tail đúng file không?
   Kiểm tra `filebeat/filebeat.yml` → paths có trỏ đúng `/logs/**/*.log`?

4. Log file có tồn tại không?
   `ls -la ./logs/demo-node/` và `./logs/demo-go/`

5. Grok parse có lỗi không?
   Xem logstash log tìm dòng `_grokparsefailure`

6. ES có index nào chưa?
   `curl http://localhost:9200/_cat/indices?v`
