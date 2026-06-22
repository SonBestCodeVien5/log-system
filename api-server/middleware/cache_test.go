package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestDashboardNoCache(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		path        string
		wantNoCache bool
	}{
		{name: "dashboard root", path: "/", wantNoCache: true},
		{name: "dashboard asset", path: "/assets/app.js", wantNoCache: true},
		{name: "api response", path: "/api/health", wantNoCache: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(DashboardNoCache())
			router.GET(tt.path, func(c *gin.Context) {
				c.Status(http.StatusNoContent)
			})

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, tt.path, nil)
			router.ServeHTTP(recorder, request)

			cacheControl := recorder.Header().Get("Cache-Control")
			if tt.wantNoCache {
				if cacheControl != "no-store, no-cache, must-revalidate, max-age=0" {
					t.Fatalf("Cache-Control = %q", cacheControl)
				}
				if recorder.Header().Get("Pragma") != "no-cache" {
					t.Fatalf("Pragma = %q", recorder.Header().Get("Pragma"))
				}
				if recorder.Header().Get("Expires") != "0" {
					t.Fatalf("Expires = %q", recorder.Header().Get("Expires"))
				}
				return
			}

			if cacheControl != "" {
				t.Fatalf("API Cache-Control = %q, want empty", cacheControl)
			}
		})
	}
}
