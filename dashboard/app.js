// =================================================================
// Log System Dashboard — app.js
// Fetch API + WebSocket + Filter + Pagination + Alert Banner
// =================================================================

const API_BASE = "http://localhost:8080";
const WS_URL   = "ws://localhost:8080/ws/alerts";
const PAGE_SIZE = 20;

// ---------------------------------------------------------------
// State
// ---------------------------------------------------------------
let currentPage  = 1;
let totalLogs    = 0;
let ws           = null;
let refreshTimer = null;

// ---------------------------------------------------------------
// Init
// ---------------------------------------------------------------
document.addEventListener("DOMContentLoaded", () => {
  fetchLogs();
  fetchStats();
  connectWebSocket();
  startAutoRefresh();

  // Enter để tìm kiếm
  document.getElementById("filter-search")
    .addEventListener("keydown", e => { if (e.key === "Enter") applyFilter(); });

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
// ---------------------------------------------------------------
async function fetchStats() {
  try {
    const res  = await fetch(`${API_BASE}/api/logs/count?from=now-1h`);
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

  tbody.innerHTML = logs.map(log => `
    <tr>
      <td class="ts-cell">${formatTimestamp(log["@timestamp"])}</td>
      <td><span class="badge badge-${log.level || "INFO"}">${log.level || "—"}</span></td>
      <td class="svc-cell">${log.service || "—"}</td>
      <td class="msg-cell">${escapeHtml(log.log_message || "—")}</td>
    </tr>
  `).join("");
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

  if (level) params.set("level", level);
  if (app)   params.set("app",   app);
  if (q)     params.set("q",     q);
  return params;
}

function applyFilter() {
  currentPage = 1;
  fetchLogs();
  fetchStats();
}

function resetFilter() {
  document.getElementById("filter-level").value  = "";
  document.getElementById("filter-app").value    = "";
  document.getElementById("filter-search").value = "";
  currentPage = 1;
  fetchLogs();
  fetchStats();
}

// ---------------------------------------------------------------
// Auto-refresh
// ---------------------------------------------------------------
function startAutoRefresh() {
  stopAutoRefresh();
  refreshTimer = setInterval(() => {
    fetchLogs();
    fetchStats();
  }, 10000);
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
        // Refresh stats ngay khi nhận alert
        fetchStats();
      }

      if (msg.type === "config") {
        // Cập nhật input threshold theo config hiện tại
        document.getElementById("threshold-input").value =
          msg.config.threshold || 10;
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
async function updateThreshold() {
  const val = parseInt(document.getElementById("threshold-input").value);
  if (!val || val < 1) {
    alert("Ngưỡng phải là số nguyên dương");
    return;
  }

  try {
    const res = await fetch(`${API_BASE}/api/alerts/config`, {
      method:  "POST",
      headers: { "Content-Type": "application/json" },
      body:    JSON.stringify({ threshold: val }),
    });
    const data = await res.json();
    if (data.status === "updated") {
      alert(`Đã cập nhật ngưỡng cảnh báo: ${val} lỗi`);
    }
  } catch (err) {
    alert("Lỗi khi cập nhật ngưỡng: " + err.message);
  }
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
    .replace(/"/g, "&quot;");
}

function showTableError(msg) {
  document.getElementById("log-tbody").innerHTML =
    `<tr class="empty-row"><td colspan="4" style="color:#dc2626">${msg}</td></tr>`;
}
