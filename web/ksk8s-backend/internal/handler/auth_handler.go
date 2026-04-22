package handler

import (
	"net/http"

	"github.com/easzlab/ksk8s/internal/middleware"
	"github.com/easzlab/ksk8s/internal/model"
	"github.com/easzlab/ksk8s/internal/repository"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userRepo    *repository.UserRepository
	settingRepo *repository.SettingRepository
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		userRepo:    repository.NewUserRepository(),
		settingRepo: repository.NewSettingRepository(),
	}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	OTPCode  string `json:"otp_code"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=64"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userRepo.GetByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// If OTP enabled, require OTP code
	if user.OTPEnabled {
		if req.OTPCode == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "otp_required", "message": "OTP code required"})
			return
		}
		if !middleware.ValidateOTP(user.OTPSecret, req.OTPCode) {
			c.JSON(http.StatusForbidden, gin.H{"error": "invalid_otp", "message": "Invalid OTP code"})
			return
		}
	}

	token, err := middleware.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	regEnabled, _ := h.settingRepo.Get("registration_enabled")

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":          user.ID,
			"username":    user.Username,
			"role":        model.NormalizeRole(user.Role),
			"otp_enabled": user.OTPEnabled,
		},
		"registration_enabled": regEnabled == "true",
	})
}

func (h *AuthHandler) Register(c *gin.Context) {
	regEnabled, _ := h.settingRepo.Get("registration_enabled")
	if regEnabled != "true" {
		c.JSON(http.StatusForbidden, gin.H{"error": "registration is disabled"})
		return
	}

	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	// Self-registration defaults to cluster_admin; explicit role only allowed for platform_admin
	role := model.RoleClusterAdmin
	if req.Role != "" {
		// Only platform_admin can assign roles during registration
		callerRole, _ := c.Get("role")
		if callerRole == model.RolePlatformAdmin && model.IsValidRole(req.Role) {
			role = req.Role
		}
	}

	user := &model.User{
		Username:     req.Username,
		PasswordHash: string(hash),
		Email:        req.Email,
		Role:         role,
	}

	if err := h.userRepo.Create(user); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "username already exists"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": user.ID, "username": user.Username, "role": user.Role})
}

func (h *AuthHandler) Settings(c *gin.Context) {
	regEnabled, _ := h.settingRepo.Get("registration_enabled")
	c.JSON(http.StatusOK, gin.H{
		"registration_enabled": regEnabled == "true",
	})
}

// OTPSetup generates a new OTP secret for the user.
func (h *AuthHandler) OTPSetup(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(int64)

	user, err := h.userRepo.GetByID(uid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	secret, err := middleware.GenerateOTPSecret()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate OTP secret"})
		return
	}

	user.OTPSecret = secret
	user.OTPEnabled = false // Must verify first before enabling
	if err := h.userRepo.Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save OTP secret"})
		return
	}

	issuer := "ksk8s"
	url := middleware.GenerateOTPURL(secret, user.Username, issuer)

	c.JSON(http.StatusOK, gin.H{
		"secret": secret,
		"url":    url,
		"qr":     url, // Frontend can generate QR from this URL
	})
}

// OTPVerify verifies an OTP code and enables OTP for the user.
func (h *AuthHandler) OTPVerify(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(int64)

	var req struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userRepo.GetByID(uid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if user.OTPSecret == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "OTP not set up"})
		return
	}

	if !middleware.ValidateOTP(user.OTPSecret, req.Code) {
		c.JSON(http.StatusForbidden, gin.H{"error": "invalid_otp"})
		return
	}

	user.OTPEnabled = true
	if err := h.userRepo.Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to enable OTP"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP enabled"})
}

// OTPDisable disables OTP for the user (requires password + OTP code).
func (h *AuthHandler) OTPDisable(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(int64)

	var req struct {
		Password string `json:"password" binding:"required"`
		OTPCode  string `json:"otp_code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userRepo.GetByID(uid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if !user.OTPEnabled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "OTP not enabled"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid password"})
		return
	}

	if !middleware.ValidateOTP(user.OTPSecret, req.OTPCode) {
		c.JSON(http.StatusForbidden, gin.H{"error": "invalid_otp"})
		return
	}

	user.OTPEnabled = false
	user.OTPSecret = ""
	if err := h.userRepo.Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to disable OTP"})
		return
	}

	middleware.ClearOTPVerifiedSession(uid)
	c.JSON(http.StatusOK, gin.H{"message": "OTP disabled"})
}
