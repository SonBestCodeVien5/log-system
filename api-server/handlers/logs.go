package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
)

// esQueryTimeout — timeout cho mỗi query ES, tránh treo goroutine khi ES chậm
const esQueryTimeout = 5 * time.Second

// ---------------------------------------------------------------
// LogHandler
// ---------------------------------------------------------------
type LogHandler struct {
	es *elasticsearch.Client
}

func NewLogHandler(es *elasticsearch.Client) *LogHandler {
	return &LogHandler{es: es}
}

// ---------------------------------------------------------------
// LogEntry — cấu trúc 1 document log trong ES
// ---------------------------------------------------------------
type LogEntry struct {
	Timestamp string         `json:"@timestamp"`
	Level     string         `json:"level"`
	Service   string         `json:"service"`
	Message   string         `json:"log_message"`
	Metadata  map[string]any `json:"metadata,omitempty"`
}

// ---------------------------------------------------------------
// Response format chuẩn
// ---------------------------------------------------------------
type LogResponse struct {
	Data  []LogEntry `json:"data"`
	Total int64      `json:"total"`
	Page  int        `json:"page"`
	Size  int        `json:"size"`
}

// ---------------------------------------------------------------
// GET /api/logs
// Query params: level, app, from, to, q, page, size
// ---------------------------------------------------------------
func (h *LogHandler) GetLogs(c *gin.Context) {
	// Parse query params
	level  := c.Query("level")
	app    := c.Query("app")
	from   := c.DefaultQuery("from", "now-1h")
	to     := c.DefaultQuery("to", "now")
	q      := c.Query("q")
	page   := parseIntDefault(c.DefaultQuery("page", "1"), 1)
	size   := parseIntDefault(c.DefaultQuery("size", "20"), 20)

	if page < 1 { page = 1 }
	if size < 1 || size > 100 { size = 20 }
	from_offset := (page - 1) * size

	// Build ES query
	query := buildLogsQuery(level, app, from, to, q)
	body  := map[string]any{
		"from":  from_offset,
		"size":  size,
		"sort":  []map[string]any{{"@timestamp": map[string]string{"order": "desc"}}},
		"query": query,
	}

	bodyBytes, _ := json.Marshal(body)

	ctx, cancel := context.WithTimeout(c.Request.Context(), esQueryTimeout)
	defer cancel()

	res, err := h.es.Search(
		h.es.Search.WithContext(ctx),
		h.es.Search.WithIndex("logs-*"),
		h.es.Search.WithBody(bytes.NewReader(bodyBytes)),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		c.JSON(http.StatusInternalServerError, gin.H{"error": res.Status()})
		return
	}

	// Parse response
	var esResp map[string]any
	if err := json.NewDecoder(res.Body).Decode(&esResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse ES response"})
		return
	}

	entries, total := extractHits(esResp)

	c.JSON(http.StatusOK, LogResponse{
		Data:  entries,
		Total: total,
		Page:  page,
		Size:  size,
	})
}

// ---------------------------------------------------------------
// GET /api/logs/count
// Query params: from, to, app
// ---------------------------------------------------------------
func (h *LogHandler) CountLogs(c *gin.Context) {
	from := c.DefaultQuery("from", "now-1h")
	to   := c.DefaultQuery("to", "now")
	app  := c.Query("app")

	levels   := []string{"INFO", "WARN", "ERROR"}
	counts   := map[string]int64{}
	var total int64

	for _, level := range levels {
		count, err := h.countOneLevel(c.Request.Context(), level, app, from, to)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		counts[level] = count
		total += count
	}

	c.JSON(http.StatusOK, gin.H{
		"INFO":  counts["INFO"],
		"WARN":  counts["WARN"],
		"ERROR": counts["ERROR"],
		"total": total,
		"from":  from,
		"to":    to,
	})
}

// countOneLevel — đếm log của 1 level. Tách ra hàm riêng để response body
// được close ngay trong mỗi vòng lặp (defer trong for sẽ giữ tới khi
// caller return, dễ leak khi index lớn / kết nối ES chậm).
func (h *LogHandler) countOneLevel(parent context.Context, level, app, from, to string) (int64, error) {
	query := buildLogsQuery(level, app, from, to, "")
	bodyBytes, err := json.Marshal(map[string]any{"query": query})
	if err != nil {
		return 0, err
	}

	ctx, cancel := context.WithTimeout(parent, esQueryTimeout)
	defer cancel()

	res, err := h.es.Count(
		h.es.Count.WithContext(ctx),
		h.es.Count.WithIndex("logs-*"),
		h.es.Count.WithBody(bytes.NewReader(bodyBytes)),
	)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return 0, fmt.Errorf("es count error: %s", res.Status())
	}

	var countResp map[string]any
	if err := json.NewDecoder(res.Body).Decode(&countResp); err != nil {
		return 0, err
	}

	if count, ok := countResp["count"].(float64); ok {
		return int64(count), nil
	}
	return 0, nil
}

// ---------------------------------------------------------------
// buildLogsQuery — xây ES bool query từ filter params
// ---------------------------------------------------------------
func buildLogsQuery(level, app, from, to, q string) map[string]any {
	must   := []map[string]any{}
	filter := []map[string]any{}

	// Filter theo level (exact match)
	if level != "" {
		must = append(must, map[string]any{
			"term": map[string]any{"level.keyword": level},
		})
	}

	// Filter theo service/app (exact match)
	if app != "" {
		must = append(must, map[string]any{
			"term": map[string]any{"service.keyword": app},
		})
	}

	// Full-text search trong message
	if q != "" {
		must = append(must, map[string]any{
			"match": map[string]any{"log_message": q},
		})
	}

	// Time range filter
	filter = append(filter, map[string]any{
		"range": map[string]any{
			"@timestamp": map[string]string{
				"gte": from,
				"lte": to,
			},
		},
	})

	return map[string]any{
		"bool": map[string]any{
			"must":   must,
			"filter": filter,
		},
	}
}

// ---------------------------------------------------------------
// extractHits — parse ES response thành []LogEntry
// ---------------------------------------------------------------
func extractHits(esResp map[string]any) ([]LogEntry, int64) {
	entries := []LogEntry{}
	var total int64

	hitsOuter, ok := esResp["hits"].(map[string]any)
	if !ok {
		return entries, 0
	}

	// Total count
	if t, ok := hitsOuter["total"].(map[string]any); ok {
		if v, ok := t["value"].(float64); ok {
			total = int64(v)
		}
	}

	// Hits array
	hits, ok := hitsOuter["hits"].([]any)
	if !ok {
		return entries, total
	}

	for _, h := range hits {
		hit, ok := h.(map[string]any)
		if !ok {
			continue
		}
		src, ok := hit["_source"].(map[string]any)
		if !ok {
			continue
		}

		// Re-marshal và unmarshal vào struct
		data, err := json.Marshal(src)
		if err != nil {
			continue
		}
		var entry LogEntry
		if err := json.Unmarshal(data, &entry); err != nil {
			continue
		}
		entries = append(entries, entry)
	}

	return entries, total
}

// ---------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------
func parseIntDefault(s string, def int) int {
	if n, err := strconv.Atoi(s); err == nil {
		return n
	}
	return def
}
