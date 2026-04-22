package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/easzlab/ksk8s/internal/config"
	"github.com/easzlab/ksk8s/internal/handler"
	"github.com/easzlab/ksk8s/internal/middleware"
	"github.com/easzlab/ksk8s/internal/repository"
	"github.com/easzlab/ksk8s/internal/service"
	"github.com/easzlab/ksk8s/internal/tls"
	"github.com/easzlab/ksk8s/internal/websocket"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	// Init database
	db, err := repository.InitDB()
	if err != nil {
		log.Fatalf("failed to init database: %v", err)
	}
	_ = db
	if err := repository.AutoMigrate(); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	// Init default settings
	settingRepo := repository.NewSettingRepository()
	settingRepo.InitDefaults()

	// Boot-time JWT secret generation
	jwtSecret := os.Getenv("KSK8S_JWT_SECRET")
	if jwtSecret == "" {
		secretPath := filepath.Join(os.Getenv("KSK8S_DATA_DIR"), ".jwt_secret")
		if secretPath == "/.jwt_secret" {
			secretPath = "/var/lib/ksk8s/.jwt_secret"
		}
		generated, err := config.GenerateJWTSecret(secretPath)
		if err != nil {
			log.Printf("failed to generate JWT secret: %v", err)
		} else {
			jwtSecret = generated
			log.Println("generated new JWT secret")
		}
	}
	middleware.InitJWTSecret(jwtSecret)

	// Self-signed TLS generation
	certPath := filepath.Join(os.Getenv("KSK8S_DATA_DIR"), "tls", "cert.pem")
	keyPath := filepath.Join(os.Getenv("KSK8S_DATA_DIR"), "tls", "key.pem")
	if certPath == "/tls/cert.pem" {
		certPath = "/var/lib/ksk8s/tls/cert.pem"
		keyPath = "/var/lib/ksk8s/tls/key.pem"
	}
	if err := tls.GenerateSelfSigned(certPath, keyPath); err != nil {
		log.Printf("failed to generate TLS certs: %v", err)
	}

	// WebSocket hub
	hub := websocket.NewHub()
	go hub.Run()

	// Log ring map (shared between task service and websocket handler)
	logRings := websocket.NewLogRingMap()

	// Task service (singleton, shared across handlers)
	taskService := service.NewTaskService(logRings)

	// Backend restart recovery: mark orphaned running tasks as failed
	go taskService.RecoverOnStartup()

	// Start ring broadcaster: watches all log rings and broadcasts new lines to WS clients
	go startRingBroadcaster(hub, logRings)

	// Log retention cron: compress 30d+ logs, delete 90d+ logs
	go startLogRetention()

	// Handlers
	authHandler := handler.NewAuthHandler()
	clusterHandler := handler.NewClusterHandler()
	templateHandler := handler.NewTemplateHandler()
	taskHandler := handler.NewTaskHandler(taskService)
	wsHandler := handler.NewWebSocketHandler(hub, logRings)
	sshHandler := handler.NewSSHHandler()
	nodeOpsHandler := handler.NewNodeOpsHandler(taskService)
	auditHandler := handler.NewAuditHandler()
	userHandler := handler.NewUserHandler()

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(middleware.AuditMiddleware())

	// Public routes
	r.POST("/api/auth/login", authHandler.Login)
	r.POST("/api/auth/register", authHandler.Register)
	r.GET("/api/auth/settings", authHandler.Settings)
	r.GET("/healthz", handler.Healthz)
	r.GET("/readyz", handler.Readyz)

	// Protected routes
	api := r.Group("/api")
	api.Use(middleware.JWTAuth())
	{
		// Auth / OTP
		api.POST("/auth/otp/setup", authHandler.OTPSetup)
		api.POST("/auth/otp/verify", authHandler.OTPVerify)
		api.POST("/auth/otp/disable", authHandler.OTPDisable)

		// User management (platform_admin only)
		api.GET("/users", middleware.RequirePlatformAdmin(), userHandler.List)
		api.POST("/users", middleware.RequirePlatformAdmin(), userHandler.Create)
		api.GET("/users/:id", middleware.RequirePlatformAdmin(), userHandler.Get)
		api.PUT("/users/:id/role", middleware.RequirePlatformAdmin(), userHandler.UpdateRole)
		api.PUT("/users/:id/password", middleware.RequirePlatformAdmin(), userHandler.ResetPassword)
		api.PUT("/users/:id/otp", middleware.RequirePlatformAdmin(), userHandler.ToggleOTP)
		api.PUT("/users/:id/language", middleware.RequirePlatformAdmin(), userHandler.UpdateLanguage)
		api.GET("/settings/registration", middleware.RequirePlatformAdmin(), userHandler.GetRegistrationSetting)
		api.PUT("/settings/registration", middleware.RequirePlatformAdmin(), userHandler.SetRegistrationSetting)
		api.POST("/users/bind-cluster", middleware.RequirePlatformAdmin(), userHandler.BindCluster)
		api.POST("/users/unbind-cluster", middleware.RequirePlatformAdmin(), userHandler.UnbindCluster)
		api.GET("/users/:id/bindings", middleware.RequirePlatformAdmin(), userHandler.ListBindings)

		// Clusters
		api.POST("/clusters", middleware.RequireClusterAdminOrAbove(), clusterHandler.Create)
		api.GET("/clusters", clusterHandler.List)
		api.GET("/clusters/:id", middleware.RequireClusterAccess(), clusterHandler.Get)
		api.PUT("/clusters/:id", middleware.RequireClusterWrite(), clusterHandler.Update)
		api.GET("/clusters/:id/config", middleware.RequireClusterAccess(), clusterHandler.GetConfig)
		api.PUT("/clusters/:id/config", middleware.RequireClusterWrite(), clusterHandler.SaveConfig)
		api.POST("/clusters/:id/generate-config", middleware.RequireClusterAccess(), clusterHandler.GenerateConfig)
		api.DELETE("/clusters/:id", middleware.OTPVerifyMiddleware(), middleware.RequireClusterWrite(), clusterHandler.Delete)

		// Node ops (high-risk)
		api.POST("/clusters/:id/nodes", middleware.OTPVerifyMiddleware(), middleware.RequireClusterWrite(), nodeOpsHandler.AddNode)
		api.DELETE("/clusters/:id/nodes", middleware.OTPVerifyMiddleware(), middleware.RequireClusterWrite(), nodeOpsHandler.RemoveNode)

		// Templates
		api.POST("/templates", middleware.RequirePlatformAdmin(), templateHandler.Create)
		api.GET("/templates", templateHandler.List)
		api.GET("/templates/:id", templateHandler.Get)
		api.GET("/templates/:id/parsed", templateHandler.GetParsed)
		api.PUT("/templates/:id", middleware.RequirePlatformAdmin(), templateHandler.Update)
		api.DELETE("/templates/:id", middleware.RequirePlatformAdmin(), templateHandler.Delete)
		api.POST("/templates/:id/set-default", middleware.RequirePlatformAdmin(), templateHandler.SetDefault)

		// Tasks
		api.GET("/clusters/:id/tasks", middleware.RequireClusterAccess(), taskHandler.List)
		api.GET("/clusters/:id/tasks/:taskId", middleware.RequireClusterAccess(), taskHandler.Get)
		api.POST("/clusters/:id/steps/:step/run", middleware.OTPVerifyMiddleware(), middleware.RequireClusterWrite(), taskHandler.RunStep)
		api.POST("/clusters/:id/tasks/:taskId/abort", middleware.OTPVerifyMiddleware(), middleware.RequireClusterWrite(), taskHandler.Abort)
		api.POST("/clusters/:id/tasks/:taskId/approve", middleware.OTPVerifyMiddleware(), middleware.RequireClusterWrite(), taskHandler.Approve)
		api.POST("/clusters/:id/steps/:step/retry", middleware.OTPVerifyMiddleware(), middleware.RequireClusterWrite(), taskHandler.Retry)
		api.GET("/clusters/:id/tasks/:taskId/logs", middleware.RequireClusterAccess(), taskHandler.Logs)
		api.GET("/clusters/:id/tasks/:taskId/status", middleware.RequireClusterAccess(), taskHandler.Status)

		// Audit logs
		api.GET("/audit-logs", middleware.RequireAuditorOrAbove(), auditHandler.List)
	}

	// WebSocket endpoints
	r.GET("/ws/tasks/:id/logs", wsHandler.ServeWS)
	r.GET("/ws/ssh", sshHandler.ServeWS)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("ksk8s backend starting on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

// startRingBroadcaster polls log rings and broadcasts new lines to WebSocket clients.
func startRingBroadcaster(hub *websocket.Hub, rings *websocket.LogRingMap) {
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	lastTotals := make(map[int64]int)

	for range ticker.C {
		rings.Iterate(func(taskID int64, ring *websocket.LogRing) {
			_, total := ring.Since(0)
			last, ok := lastTotals[taskID]
			if !ok {
				lastTotals[taskID] = total
				return
			}
			if total > last {
				lines, newTotal := ring.Since(last)
				for _, line := range lines {
					payload, _ := json.Marshal(line)
					hub.Broadcast(taskID, payload)
				}
				lastTotals[taskID] = newTotal
			}
		})
	}
}

// startLogRetention runs a daily job to compress logs older than 30 days and delete logs older than 90 days.
func startLogRetention() {
	logDir := "/var/log/ksk8s"
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	// Run once on startup
	runRetention(logDir)

	for range ticker.C {
		runRetention(logDir)
	}
}

func runRetention(logDir string) {
	// This is a placeholder. In production, implement with filepath.Walk
	// to compress .log files older than 30d and delete .gz files older than 90d.
	log.Printf("[retention] scanning %s for old logs", logDir)
}
