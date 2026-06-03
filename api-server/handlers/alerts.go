package handlers

import (
	"net/http"

	"github.com/SonBestCodeVien5/log-system/api-server/alerting"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// ---------------------------------------------------------------
// AlertHandler
// ---------------------------------------------------------------
type AlertHandler struct {
	engine *alerting.AlertEngine
}

func NewAlertHandler(engine *alerting.AlertEngine) *AlertHandler {
	return &AlertHandler{engine: engine}
}

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // cho phép mọi origin trong dev
	},
}

// ---------------------------------------------------------------
// HandleWS — GET /ws/alerts
// Dashboard kết nối WebSocket để nhận alert real-time
// ---------------------------------------------------------------
func (h *AlertHandler) HandleWS(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	// Đăng ký client
	h.engine.RegisterClient(conn)
	defer h.engine.UnregisterClient(conn)

	// Gửi config hiện tại ngay khi connect
	cfg := h.engine.GetConfig()
	conn.WriteJSON(map[string]any{
		"type":   "config",
		"config": cfg,
	})

	// Giữ connection — đọc message từ client (ping/pong)
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

// ---------------------------------------------------------------
// UpdateConfig — POST /api/alerts/config
// Dashboard gửi threshold mới — Dynamic Threshold
// ---------------------------------------------------------------
func (h *AlertHandler) UpdateConfig(c *gin.Context) {
	var cfg alerting.AlertConfig
	if err := c.ShouldBindJSON(&cfg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid config: " + err.Error()})
		return
	}

	if cfg.Threshold < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "threshold must be >= 1"})
		return
	}

	h.engine.UpdateConfig(cfg)

	c.JSON(http.StatusOK, gin.H{
		"status":  "updated",
		"config":  h.engine.GetConfig(),
	})
}
