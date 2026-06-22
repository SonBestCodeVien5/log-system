// =================================================================
// Log System Dashboard — app.js
// Fetch API + WebSocket + Filter + Pagination + Alert Banner
// =================================================================

// API và WebSocket cùng host/port với trang đang serve dashboard.
// Tránh hardcode "localhost:8080" — sẽ vỡ khi truy cập qua tên server,
// reverse proxy, hoặc đổi API_PORT trong .env.
const API_BASE = window.location.origin;
const WS_URL   =
  (window.location.protocol === "https:" ? "wss://" : "ws://") +
  window.location.host + "/ws/alerts";
const PAGE_SIZE = 20;

// ---------------------------------------------------------------
// State
// ---------------------------------------------------------------
let currentPage  = 1;
let totalLogs    = 0;
let ws           = null;
let refreshTimer = null;
let usesSlidingTimeRange = true;

// ---------------------------------------------------------------
// Init
// ---------------------------------------------------------------
document.addEventListener("DOMContentLoaded", () => {
  setDefaultTimeRange();

  ["filter-from", "filter-to"].forEach(id => {
    document.getElementById(id).addEventListener("input", () => {
      usesSlidingTimeRange = false;
      setFilterError("");
    });
  });

  refreshDashboard({ advanceTimeRange: false });
  connectWebSocket();
  startAutoRefresh();

  // Enter để tìm kiếm
  document.getElementById("filter-search")
    .addEventListener("keydown", e => { if (e.key === "Enter") applyFilter(); });

  ["threshold-input", "window-input", "cooldown-input"].forEach(id => {
    document.getElementById(id)
      .addEventListener("keydown", e => { if (e.key === "Enter") updateAlertConfig(); });
  });

  // Toggle auto-refresh
  document.getElementById("auto-refresh")
    .addEventListener("change", e => {
      if (e.target.checked) startAutoRefresh();
      else stopAutoRefresh();
    });
});

// ---------------------------------------------------------------
// Fetch logs — GET /api/logs
// ---------------------------------------------------------------
async function fetchLogs() {
  const params = buildParams();
  if (!params) return;

  params.set("page", currentPage);
  params.set("size", PAGE_SIZE);

  try {
    const res  = await fetch(`${API_BASE}/api/logs?${params}`);
    const data = await res.json();

    totalLogs = data.total || 0;
    renderTable(data.data || []);
    renderPagination();
  } catch (err) {
    showTableError("Không thể kết nối API. Kiểm tra server đang chạy.");
  }
}

// ---------------------------------------------------------------
// Fetch stats — GET /api/logs/count
// Tôn trọng filter app hiện tại để stats khớp với bảng bên dưới.
// (Backend hỗ trợ filter 'app'; level vẫn hiển thị riêng từng cột.)
// ---------------------------------------------------------------
async function fetchStats() {
  const params = new URLSearchParams();
  const timeRange = getTimeRange();
  if (!timeRange) return;

  params.set("from", timeRange.from);
  params.set("to", timeRange.to);
  const app = document.getElementById("filter-app").value;
  if (app) params.set("app", app);

  try {
    const res  = await fetch(`${API_BASE}/api/logs/count?${params}`);
    const data = await res.json();

    document.getElementById("stat-total").textContent =
      formatNumber(data.total || 0);
    document.getElementById("stat-info").textContent =
      formatNumber(data.INFO  || 0);
    document.getElementById("stat-warn").textContent =
      formatNumber(data.WARN  || 0);
    document.getElementById("stat-error").textContent =
      formatNumber(data.ERROR || 0);
    document.getElementById("stat-updated").textContent =
      new Date().toLocaleTimeString("vi-VN");
  } catch (_) {}
}

// ---------------------------------------------------------------
// Render table
// ---------------------------------------------------------------
function renderTable(logs) {
  const tbody = document.getElementById("log-tbody");

  if (!logs.length) {
    tbody.innerHTML =
      `<tr class="empty-row"><td colspan="4">Không tìm thấy log nào</td></tr>`;
    return;
  }

  // Tất cả field tới từ ES đều phải escape trước khi nhúng vào HTML.
  // Trước đây chỉ message được escape, level/service nhúng raw qua
  // template literal — bất kỳ log nào có level/service chứa ký tự
  // HTML đều có thể phá layout (hoặc tệ hơn nếu pipeline log thay đổi).
  // Level cũng được whitelist để tránh tạo CSS class lạ.
  const VALID_LEVELS = new Set(["INFO", "WARN", "ERROR"]);

  tbody.innerHTML = logs.map(log => {
    const rawLevel = String(log.level || "").toUpperCase();
    const level    = VALID_LEVELS.has(rawLevel) ? rawLevel : "INFO";
    const levelTxt = escapeHtml(log.level || "—");
    const service  = escapeHtml(log.service || "—");
    const ts       = escapeHtml(formatTimestamp(log["@timestamp"]));
    const message  = escapeHtml(log.log_message || "—");

    return `
    <tr>
      <td class="ts-cell" data-label="Thời gian">${ts}</td>
      <td data-label="Level"><span class="badge badge-${level}">${levelTxt}</span></td>
      <td class="svc-cell" data-label="Service">${service}</td>
      <td class="msg-cell" data-label="Message">${message}</td>
    </tr>`;
  }).join("");
}

// ---------------------------------------------------------------
// Pagination
// ---------------------------------------------------------------
function renderPagination() {
  const totalPages = Math.max(1, Math.ceil(totalLogs / PAGE_SIZE));

  document.getElementById("page-info").textContent =
    `Trang ${currentPage} / ${totalPages}`;
  document.getElementById("total-info").textContent =
    `${formatNumber(totalLogs)} log`;
  document.getElementById("btn-prev").disabled = currentPage <= 1;
  document.getElementById("btn-next").disabled = currentPage >= totalPages;
}

function prevPage() {
  if (currentPage > 1) { currentPage--; fetchLogs(); }
}

function nextPage() {
  const totalPages = Math.ceil(totalLogs / PAGE_SIZE);
  if (currentPage < totalPages) { currentPage++; fetchLogs(); }
}

// ---------------------------------------------------------------
// Filter
// ---------------------------------------------------------------
function buildParams() {
  const params = new URLSearchParams();
  const level  = document.getElementById("filter-level").value;
  const app    = document.getElementById("filter-app").value;
  const q      = document.getElementById("filter-search").value.trim();
  const timeRange = getTimeRange();

  if (!timeRange) return null;

  if (level) params.set("level", level);
  if (app)   params.set("app",   app);
  if (q)     params.set("q",     q);
  params.set("from", timeRange.from);
  params.set("to", timeRange.to);
  return params;
}

function applyFilter() {
  if (!getTimeRange()) return;

  currentPage = 1;
  refreshDashboard();
}

function resetFilter() {
  document.getElementById("filter-level").value  = "";
  document.getElementById("filter-app").value    = "";
  document.getElementById("filter-search").value = "";
  usesSlidingTimeRange = true;
  setDefaultTimeRange();
  currentPage = 1;
  refreshDashboard({ advanceTimeRange: false });
}

function setDefaultTimeRange() {
  const now = new Date();
  const oneHourAgo = new Date(now.getTime() - 60 * 60 * 1000);

  document.getElementById("filter-from").value = toLocalDateTimeValue(oneHourAgo);
  document.getElementById("filter-to").value = toLocalDateTimeValue(now);
  setFilterError("");
}

function getTimeRange() {
  const fromValue = document.getElementById("filter-from").value;
  const toValue = document.getElementById("filter-to").value;

  if (!fromValue || !toValue) {
    setFilterError("Vui lòng chọn đầy đủ thời gian Từ và Đến.");
    return null;
  }

  const from = new Date(fromValue);
  const to = new Date(toValue);

  if (Number.isNaN(from.getTime()) || Number.isNaN(to.getTime())) {
    setFilterError("Khoảng thời gian không hợp lệ.");
    return null;
  }

  if (from > to) {
    setFilterError("Thời gian Từ phải sớm hơn hoặc bằng thời gian Đến.");
    return null;
  }

  setFilterError("");
  return { from: from.toISOString(), to: to.toISOString() };
}

function toLocalDateTimeValue(date) {
  const timezoneOffset = date.getTimezoneOffset() * 60 * 1000;
  return new Date(date.getTime() - timezoneOffset).toISOString().slice(0, 19);
}

function setFilterError(message) {
  document.getElementById("filter-error").textContent = message;
}

function refreshDashboard({ advanceTimeRange = true } = {}) {
  if (advanceTimeRange && usesSlidingTimeRange) setDefaultTimeRange();
  fetchLogs();
  fetchStats();
}

// ---------------------------------------------------------------
// Auto-refresh
// ---------------------------------------------------------------
function startAutoRefresh() {
  stopAutoRefresh();
  refreshTimer = setInterval(refreshDashboard, 10000);
}

function stopAutoRefresh() {
  if (refreshTimer) { clearInterval(refreshTimer); refreshTimer = null; }
}

// ---------------------------------------------------------------
// WebSocket — nhận alert real-time
// ---------------------------------------------------------------
function connectWebSocket() {
  try {
    ws = new WebSocket(WS_URL);
  } catch (_) { return; }

  ws.onopen = () => {
    setWsStatus(true);
  };

  ws.onmessage = (event) => {
    try {
      const msg = JSON.parse(event.data);

      if (msg.type === "error_spike") {
        showAlert(
          `⚠ ${msg.count} lỗi trong ${msg.window} vừa qua` +
          ` (ngưỡng: ${msg.threshold})`
        );
        // Tiến sliding range một lần rồi refresh đồng bộ bảng và stats.
        refreshDashboard();
      }

      if (msg.type === "config") {
        // Cập nhật các input config theo trạng thái hiện tại
        setAlertConfigInputs(msg.config);
        setConfigStatus("Đã đồng bộ", "ok");
      }
    } catch (_) {}
  };

  ws.onclose = () => {
    setWsStatus(false);
    // Reconnect sau 5 giây
    setTimeout(connectWebSocket, 5000);
  };

  ws.onerror = () => {
    setWsStatus(false);
  };
}

function setWsStatus(connected) {
  const el = document.getElementById("ws-status");
  el.textContent  = connected ? "● Connected" : "● Disconnected";
  el.className    = "ws-badge " + (connected ? "connected" : "disconnected");
}

// ---------------------------------------------------------------
// Alert Banner
// ---------------------------------------------------------------
function showAlert(message) {
  document.getElementById("alert-text").textContent = message;
  document.getElementById("alert-banner").classList.remove("hidden");
}

function closeAlert() {
  document.getElementById("alert-banner").classList.add("hidden");
}

// ---------------------------------------------------------------
// Dynamic Threshold — POST /api/alerts/config
// ---------------------------------------------------------------
async function updateAlertConfig() {
  const threshold = parseInt(document.getElementById("threshold-input").value);
  const windowSeconds = parseInt(document.getElementById("window-input").value);
  const cooldownSeconds = parseInt(document.getElementById("cooldown-input").value);

  if (!threshold || threshold < 1) {
    setConfigStatus("Ngưỡng phải là số nguyên dương", "error");
    return;
  }
  if (!windowSeconds || windowSeconds < 1) {
    setConfigStatus("Window phải là số nguyên dương", "error");
    return;
  }
  if (!cooldownSeconds || cooldownSeconds < 1) {
    setConfigStatus("Cooldown phải là số nguyên dương", "error");
    return;
  }

  try {
    setConfigStatus("Đang cập nhật...", "pending");
    const res = await fetch(`${API_BASE}/api/alerts/config`, {
      method:  "POST",
      headers: { "Content-Type": "application/json" },
      body:    JSON.stringify({
        threshold,
        window_seconds: windowSeconds,
        cooldown_seconds: cooldownSeconds,
      }),
    });
    const data = await res.json();
    if (!res.ok) {
      setConfigStatus(data.error || "Không thể cập nhật cấu hình", "error");
      return;
    }
    if (data.status === "updated") {
      setAlertConfigInputs(data.config || {});
      setConfigStatus("Đã cập nhật", "ok");
    }
  } catch (err) {
    setConfigStatus("Lỗi khi cập nhật cấu hình: " + err.message, "error");
  }
}

function setAlertConfigInputs(config) {
  document.getElementById("threshold-input").value =
    config.threshold || 10;
  document.getElementById("window-input").value =
    config.window_seconds || 300;
  document.getElementById("cooldown-input").value =
    config.cooldown_seconds || 60;
}

function setConfigStatus(message, state) {
  const el = document.getElementById("config-status");
  el.textContent = message;
  el.className = `config-status ${state || ""}`.trim();
}

// ---------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------
function formatTimestamp(ts) {
  if (!ts) return "—";
  try {
    return new Date(ts).toLocaleString("vi-VN", {
      year:   "numeric", month:  "2-digit", day:    "2-digit",
      hour:   "2-digit", minute: "2-digit", second: "2-digit",
      hour12: false,
    });
  } catch (_) { return ts; }
}

function formatNumber(n) {
  return Number(n).toLocaleString("vi-VN");
}

function escapeHtml(str) {
  return String(str)
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#39;")
    .replace(/`/g, "&#96;");
}

function showTableError(msg) {
  document.getElementById("log-tbody").innerHTML =
    `<tr class="empty-row"><td colspan="4" style="color:#dc2626">${msg}</td></tr>`;
}
