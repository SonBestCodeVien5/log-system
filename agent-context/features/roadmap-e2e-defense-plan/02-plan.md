# Plan

## Summary
- Goal: Có roadmap/plan đọc được, đẹp hơn khi mở trực tiếp, và đủ chi tiết để chạy Bước 10-12 mà không phải lục lại chat.
- Success criteria:
  - `actual_progress_roadmap.html` là standalone HTML có giao diện rõ ràng.
  - `docs/project-roadmap.md` mô tả Bước 10-12 bằng lệnh cụ thể.
  - `docs/testing-evidence.md` có checklist để điền Actual result sau khi test.
  - `docs/architecture.md` và `docs/deployment.md` phản ánh trạng thái repo hiện tại.
  - `docs/report-notes.md` có demo script 5 phút và câu hỏi ôn bảo vệ.

## Key Changes
- Implementation:
  - Rebuild `actual_progress_roadmap.html` thành trang HTML đầy đủ với hero, timeline, checklist, command blocks và demo plan.
  - Add `docs/project-roadmap.md` làm roadmap markdown chính.
  - Update docs hiện có để trỏ sang Bước 10 end-to-end verification.
- Public API/UI/data contracts:
  - Không đổi API contract.
  - Không đổi dashboard application behavior.
  - Chỉ đổi tài liệu và roadmap static.
- Out of scope:
  - Không chạy Docker Compose trong phase này.
  - Không sửa application source.
  - Không ghi số liệu hiệu năng khi chưa đo thật.

## Acceptance Criteria
- Scenario: Người dùng mở `actual_progress_roadmap.html` trực tiếp trong browser.
- Expected result: Thấy giao diện roadmap có màu sắc, panel, timeline và checklist thay vì chữ đen cơ bản.
- Scenario: Người dùng đọc `docs/project-roadmap.md`.
- Expected result: Biết cần chạy chính xác các lệnh Bước 10.1-10.4, ghi evidence vào đâu, và chuẩn bị Bước 12 như thế nào.

## Assumptions
- Assumption: Commit `0ceef5f` là baseline đã push sạch trước khi tạo roadmap/docs mới.
- Assumption: Bước 8-9 được xem là hoàn thành về code theo repo hiện tại; Bước 10 vẫn cần verify runtime.
- Assumption: Các docs mới sẽ được review trước khi commit/push tiếp.

## Next Handoff
- Current phase: plan
- Next phase: verification
- Must read:
  - `actual_progress_roadmap.html`
  - `docs/project-roadmap.md`
  - `docs/testing-evidence.md`
- Decisions locked:
  - Không thêm code mới trước khi chạy Bước 10, trừ khi test phát hiện bug.
  - Số liệu hiệu năng chỉ ghi sau khi đo thật.
- Open risks:
  - Cần mở browser local để xác nhận visual.
  - Cần chạy Docker Compose để xác nhận dashboard/API/alerting.
- Validation status:
  - Static diff check pass.
