package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/SonBestCodeVien5/log-system/api-server/alerting"
	"github.com/gin-gonic/gin"
)

func TestUpdateConfigRejectsInvalidConfigWithoutPartialUpdate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name string
		body string
	}{
		{
			name: "invalid threshold",
			body: `{"threshold":0,"window_seconds":300,"cooldown_seconds":60}`,
		},
		{
			name: "invalid window",
			body: `{"threshold":5,"window_seconds":-1,"cooldown_seconds":60}`,
		},
		{
			name: "invalid cooldown",
			body: `{"threshold":5,"window_seconds":300,"cooldown_seconds":0}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setDefaultAlertEnv(t)

			engine := alerting.NewEngine(nil)
			handler := NewAlertHandler(engine)
			router := gin.New()
			router.POST("/api/alerts/config", handler.UpdateConfig)

			req := httptest.NewRequest(
				http.MethodPost,
				"/api/alerts/config",
				strings.NewReader(tt.body),
			)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want %d; body=%s", w.Code, http.StatusBadRequest, w.Body.String())
			}

			cfg := engine.GetConfig()
			if cfg.Threshold != 10 || cfg.WindowSeconds != 300 || cfg.CooldownSeconds != 60 {
				t.Fatalf("config changed after invalid request: %+v", cfg)
			}
		})
	}
}

func TestUpdateConfigAcceptsCompletePositiveConfig(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setDefaultAlertEnv(t)

	engine := alerting.NewEngine(nil)
	handler := NewAlertHandler(engine)
	router := gin.New()
	router.POST("/api/alerts/config", handler.UpdateConfig)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/alerts/config",
		strings.NewReader(`{"threshold":5,"window_seconds":120,"cooldown_seconds":30}`),
	)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", w.Code, http.StatusOK, w.Body.String())
	}

	cfg := engine.GetConfig()
	if cfg.Threshold != 5 || cfg.WindowSeconds != 120 || cfg.CooldownSeconds != 30 {
		t.Fatalf("config = %+v, want threshold=5 window=120 cooldown=30", cfg)
	}
}

func TestUpdateConfigAcceptsPartialThresholdConfig(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setDefaultAlertEnv(t)

	engine := alerting.NewEngine(nil)
	handler := NewAlertHandler(engine)
	router := gin.New()
	router.POST("/api/alerts/config", handler.UpdateConfig)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/alerts/config",
		strings.NewReader(`{"threshold":5}`),
	)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", w.Code, http.StatusOK, w.Body.String())
	}

	cfg := engine.GetConfig()
	if cfg.Threshold != 5 || cfg.WindowSeconds != 300 || cfg.CooldownSeconds != 60 {
		t.Fatalf("config = %+v, want threshold=5 window=300 cooldown=60", cfg)
	}
}

func TestUpdateConfigPreservesPriorPartialUpdates(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setDefaultAlertEnv(t)

	engine := alerting.NewEngine(nil)
	handler := NewAlertHandler(engine)
	router := gin.New()
	router.POST("/api/alerts/config", handler.UpdateConfig)

	sendConfig := func(body string) {
		t.Helper()
		req := httptest.NewRequest(
			http.MethodPost,
			"/api/alerts/config",
			strings.NewReader(body),
		)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d; body=%s", w.Code, http.StatusOK, w.Body.String())
		}
	}

	sendConfig(`{"threshold":5}`)
	sendConfig(`{"window_seconds":120}`)

	cfg := engine.GetConfig()
	if cfg.Threshold != 5 || cfg.WindowSeconds != 120 || cfg.CooldownSeconds != 60 {
		t.Fatalf("config = %+v, want threshold=5 window=120 cooldown=60", cfg)
	}
}

func TestUpdateConfigRejectsEmptyConfig(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setDefaultAlertEnv(t)

	engine := alerting.NewEngine(nil)
	handler := NewAlertHandler(engine)
	router := gin.New()
	router.POST("/api/alerts/config", handler.UpdateConfig)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/alerts/config",
		strings.NewReader(`{}`),
	)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body=%s", w.Code, http.StatusBadRequest, w.Body.String())
	}

	cfg := engine.GetConfig()
	if cfg.Threshold != 10 || cfg.WindowSeconds != 300 || cfg.CooldownSeconds != 60 {
		t.Fatalf("config changed after empty request: %+v", cfg)
	}
}

func setDefaultAlertEnv(t *testing.T) {
	t.Helper()
	t.Setenv("ALERT_THRESHOLD", "10")
	t.Setenv("ALERT_WINDOW_SECONDS", "300")
	t.Setenv("ALERT_COOLDOWN_SECONDS", "60")
	t.Setenv("ALERT_CHECK_INTERVAL_SECONDS", "10")
}
