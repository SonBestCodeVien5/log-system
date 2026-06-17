# Plan

## Summary
- Goal: Lam cum chinh thong so alert hoa cung ngon ngu UI voi dashboard, bot cam giac la mot card rieng chen giua stats va filter.
- Success criteria:
  - Alert config nhin nhu mot phan cua control surface chung, khong chiem hang/card rieng qua nang.
  - Desktop giu duoc mat do van hanh: stats, filter, auto-refresh, va alert config doc nhanh trong mot vung toolbar hop ly.
  - Mobile khong bi day bang log xuong qua sau; control wrap co trat tu va khong tran text.
  - Chuc nang `POST /api/alerts/config`, WebSocket config sync, va validation input giu nguyen.

## Key Changes
- Implementation:
  - Chuyen `Alert config` tu panel/card full-width thanh cum "alert tuning" compact nam gan filter/ops controls.
  - De xuat cau truc UI:
    - Mot `toolbar` hoac `control-strip` duy nhat gom filter chinh o trai, auto-refresh/status/config o phai.
    - Alert config la mot inline group co label ngan, 3 numeric inputs nho, nut update cung style voi button/filter.
    - `config-status` thanh text status nho trong cung group, khong lam header rieng.
  - Reuse visual tokens co san: border `#d1d5db`, radius `6px`, input height xap xi 29-31px, button style tu filter bar.
  - Bo hoac giam vai tro `.alert-config-panel` card: khong full-width white panel rieng neu khong can.
  - Mobile:
    - Filter controls full width nhu hien tai.
    - Alert config wrap thanh 2 hang co logic: row input numeric, row nut update + status; hoac dat trong details/disclosure neu muon tiet kiem chieu cao.
- Public API/UI/data contracts:
  - Khong doi endpoint: `POST /api/alerts/config`.
  - Khong doi request body: `threshold`, `window_seconds`, `cooldown_seconds`.
  - Khong doi WebSocket config message handling.
  - Khong doi API base/WS same-origin logic.
- Out of scope:
  - Khong doi backend alerting.
  - Khong them framework/icon lib/build step.
  - Khong doi log table, level colors, hoac pagination behavior.

## Acceptance Criteria
- Scenario: Mo dashboard desktop 1280px.
- Expected result:
  - Khong con mot card Alert config full-width nam rieng giua stats va filter.
  - Config controls can hang voi filter/ops controls, spacing va button/input dong bo.
  - Status config khong tao khoang trong lon.
- Scenario: Mo dashboard mobile 390px.
- Expected result:
  - Config controls wrap gon, khong tran ngang, khong de text chong len nhau.
  - Bang log van xuat hien sau control area voi khoang cach hop ly.
- Scenario: Cap nhat alert config.
- Expected result:
  - Input validation trong `app.js` van hien status loi/thanh cong.
  - Request gui dung `POST /api/alerts/config` voi body hien tai.
  - Khi WebSocket gui `type: "config"`, input van sync dung.

## Assumptions
- Assumption: Nguoi dung muon UI operational dense hon, khong phai them mot trang settings rieng.
- Assumption: Van nen hien 3 thong so trong dashboard vi day la tinh nang demo quan trong.
- Assumption: Neu can tiep tuc nhanh, phuong an tot nhat la "compact inline config trong toolbar", khong phai collapsible an mac dinh.

## Next Handoff
- Current phase: plan
- Next phase: implementation
- Must read: `agent-context/features/dashboard-alert-config-ui/01-discovery.md`, `dashboard/index.html`, `dashboard/style.css`, `dashboard/app.js`
- Decisions locked: Giu vanilla dashboard va hop dong API hien tai; thay doi chu yeu la layout/style cua Alert config.
- Open risks: Can verify bang Playwright desktop/mobile sau khi edit; neu API server khong chay, dung static render cho layout va test JS bang inspect/manual.
- Validation status: Plan dua tren source inspection va render tinh desktop/mobile; chua implement.
