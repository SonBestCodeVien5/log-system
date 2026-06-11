# Discovery

## Request
- Feature/request: Cập nhật roadmap/plan theo hướng Bước 10 end-to-end test, Bước 11 docs, Bước 12 bảo vệ; chỉ cập nhật tài liệu dựa trên repo hiện tại.
- Feature slug: roadmap-e2e-defense-plan

## Repo Facts
- Current state:
  - Commit `0ceef5f` (`feat: add dashboard and deployment guide`) đã được push lên `origin/main`.
  - Repo sau commit có dashboard HTML/CSS/JS, README triển khai chi tiết, Go API, alerting engine, pipeline Docker Compose.
  - Sau khi push, đã mở thêm diff tài liệu/roadmap mới: `actual_progress_roadmap.html`, `docs/architecture.md`, `docs/deployment.md`, `docs/report-notes.md`, `docs/testing-evidence.md`, `docs/project-roadmap.md`.
- Relevant files:
  - `actual_progress_roadmap.html`
  - `docs/project-roadmap.md`
  - `docs/architecture.md`
  - `docs/deployment.md`
  - `docs/report-notes.md`
  - `docs/testing-evidence.md`
  - `README.md`
  - `dashboard/index.html`, `dashboard/app.js`, `dashboard/style.css`
  - `api-server/alerting/engine.go`
- Existing constraints:
  - Dashboard dùng HTML/CSS/JS thuần.
  - Không hardcode secret; `.env` không commit.
  - Tài liệu không được ghi số liệu hiệu năng giả; chỉ điền sau khi đo thực tế.

## Applicable Instructions
- Root `AGENTS.md`:
  - Giữ kiến trúc Centralized Logging Platform: demo services -> Filebeat -> Logstash -> Elasticsearch -> Go API -> Dashboard/WebSocket.
  - README/docs cần hướng dẫn chạy và test rõ ràng.
- Area `AGENTS.md`:
  - Tác vụ này chỉ chạm tài liệu/roadmap, không sửa application source.
- Skill references:
  - `$log-plan`: tạo discovery/plan context.
  - `$log-git`: kiểm tra status/diff, ghi git handoff, chỉ commit khi user yêu cầu.

## Unknowns And Risks
- Unknowns:
  - Bước 10 end-to-end test chưa được chạy trong phiên này.
  - Chưa có output thật để điền bảng hiệu năng.
  - Playwright không có browser `chrome-for-testing`, nên chưa chụp screenshot render roadmap.
- Risks:
  - Nếu chưa rebuild `api-server`, dashboard mới có thể chưa được serve từ container đang chạy.
  - Nếu ERROR rate thấp, threshold `5` có thể vẫn cần chờ lâu hơn 10-15 giây hoặc cần tăng log rate để trigger alert.
  - `actual_progress_roadmap.html` đã là standalone HTML, nhưng cần người dùng mở bằng browser để xác nhận thẩm mỹ.

## Next Handoff
- Current phase: discovery
- Next phase: plan
- Must read:
  - `docs/project-roadmap.md`
  - `docs/testing-evidence.md`
  - `actual_progress_roadmap.html`
- Decisions locked:
  - Bước 10 là checkpoint tiếp theo.
  - Bước 11-12 là docs và bảo vệ, không thêm code mới trừ khi test phát hiện lỗi.
- Open risks:
  - Chưa verify runtime sau dashboard commit.
  - Chưa có số liệu response time thật.
- Validation status:
  - `git diff --check` pass cho diff tài liệu hiện tại.
  - Playwright visual check bị blocked vì thiếu browser.
