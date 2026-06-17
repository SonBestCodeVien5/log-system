# Discovery

## Request
- Feature/request: Ra soat UI dashboard vi cum chinh thong so alert dang lac que so voi cac phan con lai.
- Feature slug: dashboard-alert-config-ui

## Repo Facts
- Current state:
  - Dashboard la HTML + vanilla JS + CSS thuan, duoc Go API server serve `dashboard/index.html` tai `/` va static assets tai `/assets`.
  - Cum Alert config hien nam rieng giua stats bar va filter bar trong `dashboard/index.html`, gom `threshold`, `window`, `cooldown`, va nut `Cap nhat`.
  - CSS hien tai tao `.alert-config-panel` nhu mot panel/card rieng, full width, co header rieng va border/background rieng.
  - Render check bang static server tam thoi xac nhan desktop: panel chiem mot hang full width khoang 84px, noi dung control chi khoang nua hang; mobile: panel cao khoang 184px, roi filter bar tiep tuc thanh block rieng khoang 220px.
  - API/WS khong chay khi render tam thoi nen console co loi ket noi; day la runtime dependency, khong phai loi UI tinh.
- Relevant files:
  - `dashboard/index.html`: Alert config tai dong 54-77; filter bar tai dong 79-108; stats bar tai dong 30-52.
  - `dashboard/style.css`: `.ops-strip` va `.alert-config-*` tai dong 94-184; stats bar tai dong 189-224; filter bar tai dong 229-288; mobile rules tai dong 398-463.
  - `dashboard/app.js`: alert config state/update tai dong 264-326; WebSocket config sync tai dong 227-231.
- Existing constraints:
  - Khong dung framework/build step.
  - Giu logic trong `app.js`, structure trong `index.html`, styling trong `style.css`.
  - Dashboard can operational, dense, phu hop tac vu soi log lap lai.
  - Threshold control van phai goi `POST /api/alerts/config`.

## Applicable Instructions
- Root `AGENTS.md`:
  - Dashboard dung HTML + Vanilla JS thuan.
  - Pagination 20 record/trang, alert qua WebSocket, config alert qua endpoint co san.
  - Khong hardcode password/port; dashboard dang dung same-origin trong `app.js`.
- Area `dashboard/AGENTS.md`:
  - Bat buoc co threshold control gui `POST /api/alerts/config`.
  - File dashboard chi gom `index.html`, `app.js`, `style.css`.
  - Level colors giu dung INFO/WARN/ERROR.
- Skill references:
  - `phase-discovery.md`: xac dinh UI behavior truoc khi edit, giu MVP architecture.
  - `phase-dashboard.md`: UI dau tien la dashboard van hanh, dense, khong landing/marketing.

## Unknowns And Risks
- Unknowns:
  - Muc "chinh cac thong so" co nen hien day du 3 input luon hay nen thu gon bang disclosure/compact row.
  - Co can doi label sang tieng Viet dong bo voi phan con lai hay giu "Alert config / Threshold / Window / Cooldown" de dung thuat ngu ky thuat.
- Risks:
  - Neu chi doi mau/border ma van de panel full-width, cam giac lac que van con.
  - Neu gop qua nhieu control vao filter bar tren mobile, phan dau trang co the dai va lam bang log bi day xuong.
  - Neu an bot config mac dinh, tinh nang threshold control co the kho tim trong demo bao ve.

## Next Handoff
- Current phase: discovery
- Next phase: plan
- Must read: `dashboard/index.html`, `dashboard/style.css`, `dashboard/app.js`, `dashboard/AGENTS.md`, `.agents/skills/log-system-dev/references/phase-dashboard.md`
- Decisions locked: Chi lap plan, khong sua source trong phase nay; context luu o fallback `agent-context` vi `.agents/context` read-only.
- Open risks: Can chon compact inline hay collapsible config truoc khi implement neu muon thay doi hanh vi hien/doi.
- Validation status: Da inspect source va render desktop/mobile bang static server tam thoi; API/WS khong chay nen chi validate duoc layout tinh.
