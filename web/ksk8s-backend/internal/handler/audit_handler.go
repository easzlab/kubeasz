package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/easzlab/ksk8s/internal/repository"
	"github.com/gin-gonic/gin"
)

type AuditHandler struct {
	repo *repository.AuditLogRepository
}

func NewAuditHandler() *AuditHandler {
	return &AuditHandler{
		repo: repository.NewAuditLogRepository(),
	}
}

func (h *AuditHandler) List(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "50")
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 200 {
		pageSize = 50
	}
	offset := (page - 1) * pageSize

	query := parseAuditQuery(c)

	logs, err := h.repo.ListFiltered(query, offset, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch audit logs"})
		return
	}

	total, _ := h.repo.CountFiltered(query)

	c.JSON(http.StatusOK, gin.H{
		"logs":      logs,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func parseAuditQuery(c *gin.Context) *repository.AuditQuery {
	q := &repository.AuditQuery{}

	if start := c.Query("start_time"); start != "" {
		if t, err := time.Parse(time.RFC3339, start); err == nil {
			q.StartTime = &t
		}
	}
	if end := c.Query("end_time"); end != "" {
		if t, err := time.Parse(time.RFC3339, end); err == nil {
			q.EndTime = &t
		}
	}
	if action := c.Query("action"); action != "" {
		q.Action = action
	}
	if username := c.Query("username"); username != "" {
		q.Username = username
	}
	if highRisk := c.Query("is_high_risk"); highRisk != "" {
		v := highRisk == "true"
		q.IsHighRisk = &v
	}

	return q
}
