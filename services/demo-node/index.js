// =================================================================
// Demo Service A — Node.js
// Sinh log JSON Lines liên tục với tỉ lệ INFO/WARN/ERROR ngẫu nhiên
// Ghi ra file để Filebeat tail
// =================================================================

const fs   = require("fs");
const path = require("path");

// ---------------------------------------------------------------
// Config từ môi trường
// ---------------------------------------------------------------
const SERVICE_NAME = process.env.SERVICE_NAME || "demo-node";
const LOG_FILE     = process.env.LOG_FILE     || "/var/log/app/app.log";
const INTERVAL_MS  = parseInt(process.env.LOG_INTERVAL_MS || "2000");

// ---------------------------------------------------------------
// Đảm bảo thư mục log tồn tại
// ---------------------------------------------------------------
const logDir = path.dirname(LOG_FILE);
if (!fs.existsSync(logDir)) {
  fs.mkdirSync(logDir, { recursive: true });
}

const logStream = fs.createWriteStream(LOG_FILE, { flags: "a" });

// ---------------------------------------------------------------
// Danh sách message mẫu theo từng level
// ---------------------------------------------------------------
const messages = {
  INFO: [
    "User login successful",
    "Order created successfully",
    "Payment processed",
    "Cache refreshed",
    "Health check passed",
    "Config reloaded",
    "Session started",
    "Request completed in 45ms",
  ],
  WARN: [
    "Response time exceeded 500ms",
    "Retry attempt 2/3 for order",
    "Cache miss — falling back to DB",
    "Rate limit approaching threshold",
    "Disk usage at 75%",
    "Connection pool running low",
  ],
  ERROR: [
    "Payment gateway timeout after 30s",
    "Database connection refused",
    "Failed to process order: invalid data",
    "Authentication service unreachable",
    "Null pointer exception in OrderService",
    "Unhandled exception: out of memory",
  ],
};

// ---------------------------------------------------------------
// Chọn level theo tỉ lệ: INFO 60%, WARN 25%, ERROR 15%
// ---------------------------------------------------------------
function pickLevel() {
  const r = Math.random();
  if (r < 0.60) return "INFO";
  if (r < 0.85) return "WARN";
  return "ERROR";
}

function pickMessage(level) {
  const list = messages[level];
  return list[Math.floor(Math.random() * list.length)];
}

// ---------------------------------------------------------------
// Ghi 1 dòng log JSON Lines
// ---------------------------------------------------------------
function writeLog() {
  const level   = pickLevel();
  const message = pickMessage(level);

  const entry = {
    timestamp: new Date().toISOString(),
    level,
    service: SERVICE_NAME,
    message,
    metadata: {
      pid:      process.pid,
      sequence: ++sequence,
    },
  };

  // JSON Lines — mỗi dòng là 1 JSON object hoàn chỉnh
  logStream.write(JSON.stringify(entry) + "\n");

  // Cũng in ra stdout để xem trong docker logs
  console.log(`[${entry.level}] ${entry.message}`);
}

// ---------------------------------------------------------------
// Main loop
// ---------------------------------------------------------------
let sequence = 0;

console.log(`[${SERVICE_NAME}] Starting — writing to ${LOG_FILE} every ${INTERVAL_MS}ms`);

setInterval(writeLog, INTERVAL_MS);

// Graceful shutdown
process.on("SIGTERM", () => {
  console.log(`[${SERVICE_NAME}] Shutting down...`);
  logStream.end();
  process.exit(0);
});
