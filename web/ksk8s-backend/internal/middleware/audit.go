package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"strings"
	"time"

	"github.com/easzlab/ksk8s/internal/model"
	"github.com/easzlab/ksk8s/internal/repository"
	"github.com/gin-gonic/gin"
)

// AuditMiddleware logs API requests to the audit_logs table.
func AuditMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		// Capture request body for important operations
		var bodyBytes []byte
		if c.Request.Body != nil && shouldCaptureBody(c.Request.Method, c.Request.URL.Path) {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()

		// Skip health checks and websocket endpoints
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/healthz") || strings.HasPrefix(path, "/readyz") || strings.HasPrefix(path, "/ws/") {
			return
		}

		userID, _ := c.Get("user_id")
		username, _ := c.Get("username")
		if userID == nil {
			userID = int64(0)
		}
		if username == nil {
			username = "anonymous"
		}

		action := extractAction(c.Request.Method, path)
		resourceType, resourceID := extractResource(path)

		details := map[string]interface{}{
			"method":   c.Request.Method,
			"path":     path,
			"duration": duration.Milliseconds(),
		}
		if len(bodyBytes) > 0 && len(bodyBytes) < 4096 {
			// Sanitize body: remove password fields
			var bodyMap map[string]interface{}
			if err := json.Unmarshal(bodyBytes, &bodyMap); err == nil {
				delete(bodyMap, "password")
				delete(bodyMap, "password_hash")
				if sanitized, err := json.Marshal(bodyMap); err == nil {
					details["body"] = string(sanitized)
				}
			}
		}

		detailsJSON, _ := json.Marshal(details)

		// Mark as high-risk if X-OTP-Code header is present (OTP-verified operation)
		isHighRisk := c.GetHeader("X-OTP-Code") != ""

		auditLog := &model.Audit{
			UserID:       userID.(int64),
			Username:     username.(string),
			Action:       action,
			ResourceType: resourceType,
			ResourceID:   resourceID,
			Details:      string(detailsJSON),
			IPAddress:    c.ClientIP(),
			StatusCode:   status,
			IsHighRisk:   isHighRisk,
		}

		// Async insert to avoid blocking response
		go func() {
			repo := repository.NewAuditLogRepository()
			if err := repo.Create(auditLog); err != nil {
				log.Printf("[audit] failed to create audit log: %v", err)
			}
		}()
	}
}

func shouldCaptureBody(method string, path string) bool {
	if method == "GET" || method == "DELETE" {
		return false
	}
	// Skip login body (has password)
	if strings.Contains(path, "/auth/login") {
		return false
	}
	return true
}

func extractAction(method string, path string) string {
	// Map HTTP method + path to a human-readable action
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 2 {
		return method + " " + path
	}

	resource := parts[1] // e.g., clusters, templates, auth
	action := method + "_" + resource

	// Special cases for nested resources
	if strings.Contains(path, "/steps/") {
		action = "run_task"
	} else if strings.Contains(path, "/tasks/") && strings.Contains(path, "/abort") {
		action = "abort_task"
	} else if strings.Contains(path, "/tasks/") && strings.Contains(path, "/approve") {
		action = "approve_task"
	} else if strings.Contains(path, "/nodes") {
		action = method + "_node"
	} else if strings.Contains(path, "/config") {
		action = method + "_config"
	} else if strings.Contains(path, "/auth/") {
		action = strings.TrimPrefix(path, "/api/auth/")
	}

	return action
}

func extractResource(path string) (string, string) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 3 {
		return "", ""
	}

	resourceType := parts[1] // clusters, templates, etc.
	resourceID := ""
	if len(parts) >= 3 {
		resourceID = parts[2] // :id
	}

	return resourceType, resourceID
}
