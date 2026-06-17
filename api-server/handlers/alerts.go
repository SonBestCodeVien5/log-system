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

type alertConfigRequest struct {
	Threshold       *int `json:"threshold"`
	WindowSeconds   *int `json:"window_seconds"`
	CooldownSeconds *int `json:"cooldown_seconds"`
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
	var req alertConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid config: " + err.Error()})
		return
	}

	if req.Threshold == nil && req.WindowSeconds == nil && req.CooldownSeconds == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "at least one config field is required"})
		return
	}

	var cfg alerting.AlertConfig
	if req.Threshold != nil {
		if *req.Threshold < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "threshold must be >= 1"})
			return
		}
		cfg.Threshold = *req.Threshold
	}
	if req.WindowSeconds != nil {
		if *req.WindowSeconds < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "window_seconds must be >= 1"})
			return
		}
		cfg.WindowSeconds = *req.WindowSeconds
	}
	if req.CooldownSeconds != nil {
		if *req.CooldownSeconds < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "cooldown_seconds must be >= 1"})
			return
		}
		cfg.CooldownSeconds = *req.CooldownSeconds
	}

	h.engine.UpdateConfig(cfg)

	c.JSON(http.StatusOK, gin.H{
		"status": "updated",
		"config": h.engine.GetConfig(),
	})
}
