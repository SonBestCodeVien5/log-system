# Verification

## Commands
- Command: `git diff --check`
- Result: Pass.
- Important output: Không có whitespace/patch formatting error.

- Command: `docker compose config --quiet`, `node --check services/demo-node/index.js`, `bash -n scripts/trigger-error-spike.sh`
- Result: Pass.
- Important output: Compose, Node demo và incident replay script vẫn hợp lệ sau docs pass.

- Command: `GOCACHE=/tmp/log-system-go-build-cache go test ./...` trong `api-server`
- Result: Pass.
- Important output: Handler tests pass; các package api-server, alerting và middleware chưa có test file riêng.

- Command: `GOCACHE=/tmp/log-system-demo-go-cache go test ./...` trong `services/demo-go`
- Result: Pass.
- Important output: Module compile; không có test file.

- Command: Node script quét Markdown links trong `README.md` và `docs/*.md`
- Result: Pass.
- Important output: `all local markdown links resolve`.

- Command: `rg` tìm các cụm stale như `cần verify runtime`, `sẽ được cập nhật sau`, `Add decisions here`, `Evidence Placeholders`.
- Result: Pass.
- Important output: Không còn match trong README/docs.

## Scenarios
- Scenario: Người đọc mở architecture, deployment, roadmap và evidence.
- Expected: Tất cả cùng phản ánh MVP đã pass E2E, nhưng clean-clone/screenshots/rehearsal chưa hoàn thành.
- Actual: Trạng thái đã được đồng bộ và không overclaim các phần còn mở.

- Scenario: Người viết báo cáo cần nội dung nền.
- Expected: Có narrative từ bài toán, scope, architecture, pipeline, API, alerting, evidence đến hạn chế.
- Actual: `docs/report-notes.md` có 14 mục report-ready, demo script, danh sách hình và Q&A.

- Scenario: Người làm slide cần khung 10 trang.
- Expected: Mỗi slide có một message, nội dung/visual và dữ liệu traceable.
- Actual: `docs/slide-outline.md` có đủ 10 slide và asset checklist.

## Failures Or Skips
- Failure: Ruby link checker không chạy vì môi trường không cài `ruby`.
- Resolution: Chạy checker tương đương bằng Node; tất cả local Markdown link resolve.
- Skip: Không chạy clean-clone, không trigger incident mới và không chụp screenshot trong docs pass này.
- Reason: Các việc đó cần evidence phase riêng; docs hiện ghi chúng là chưa hoàn thành.

## Next Handoff
- Current phase: verification
- Next phase: evidence capture / clean-clone, sau đó Git handoff
- Must read: `agent-context/features/final-mvp-report-audit/05-review.md`, `docs/report-notes.md`, `docs/slide-outline.md`, `docs/testing-evidence.md`
- Decisions locked: scope ứng dụng đóng băng; docs phân biệt E2E evidence, phép đo dev và production limitations.
- Open risks: clean-clone, H1-H3 screenshots, university template và rehearsal chưa hoàn thành.
- Validation status: docs/static checks pass; local link check pass; no application source changed.
