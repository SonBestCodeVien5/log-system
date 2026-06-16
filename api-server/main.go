// =================================================================
// Go API Server — log-system
// gin framework + Elasticsearch client + WebSocket alerting
// =================================================================
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SonBestCodeVien5/log-system/api-server/alerting"
	"github.com/SonBestCodeVien5/log-system/api-server/handlers"
	"github.com/SonBestCodeVien5/log-system/api-server/middleware"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
)

func main() {
	// ---------------------------------------------------------------
	// Config từ môi trường
	// ---------------------------------------------------------------
	esHost   := getEnv("ES_HOST", "localhost")
	esPort   := getEnv("ES_PORT", "9200")
	esPass   := getEnv("ES_PASSWORD", "changeme123")
	apiPort  := getEnv("API_PORT", "8080")

	// ---------------------------------------------------------------
	// Elasticsearch client
	// ---------------------------------------------------------------
	esCfg := elasticsearch.Config{
		Addresses: []string{fmt.Sprintf("http://%s:%s", esHost, esPort)},
		Username:  "elastic",
		Password:  esPass,
	}

	esClient, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		log.Fatalf("[main] failed to create ES client: %v", err)
	}

	// Kiểm tra kết nối ES
	res, err := esClient.Info()
	if err != nil {
		log.Fatalf("[main] ES connection failed: %v", err)
	}
	defer res.Body.Close()
	log.Printf("[main] Elasticsearch connected")

	// ---------------------------------------------------------------
	// Alerting Engine — chạy trong goroutine nền
	// ---------------------------------------------------------------
	engine := alerting.NewEngine(esClient)
	go engine.Run()
	log.Printf("[main] Alerting engine started")

	// ---------------------------------------------------------------
	// Gin router
	// ---------------------------------------------------------------
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())

	// API routes
	logHandler   := handlers.NewLogHandler(esClient)
	alertHandler := handlers.NewAlertHandler(engine)

	api := r.Group("/api")
	{
		api.GET("/health",           healthCheck(esClient))
		api.GET("/logs",             logHandler.GetLogs)
		api.GET("/logs/count",       logHandler.CountLogs)
		api.POST("/alerts/config",   alertHandler.UpdateConfig)
	}

	// WebSocket — dùng chung alertHandler ở trên, không tạo lại mỗi request
	r.GET("/ws/alerts", alertHandler.HandleWS)

	// Dashboard — serve index.html tại root
	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	// Static assets — JS, CSS
	r.Static("/assets", "./static")
	// ---------------------------------------------------------------
	// HTTP Server với graceful shutdown
	// ---------------------------------------------------------------
	srv := &http.Server{
		Addr:         ":" + apiPort,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		log.Printf("[main] API server listening on :%s", apiPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[main] server error: %v", err)
		}
	}()

	// Đợi signal dừng
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("[main] Shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("[main] shutdown error: %v", err)
	}
	log.Println("[main] Server stopped")
}

// ---------------------------------------------------------------
// Health check handler
// ---------------------------------------------------------------
func healthCheck(es *elasticsearch.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := es.Cluster.Health()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":        "error",
				"elasticsearch": "unreachable",
			})
			return
		}
		defer res.Body.Close()
		c.JSON(http.StatusOK, gin.H{
			"status":        "ok",
			"elasticsearch": "connected",
		})
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
