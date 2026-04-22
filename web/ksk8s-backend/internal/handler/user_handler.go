package handler

import (
	"net/http"
	"strconv"

	"github.com/easzlab/ksk8s/internal/model"
	"github.com/easzlab/ksk8s/internal/repository"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	userRepo    *repository.UserRepository
	clusterRepo *repository.ClusterRepository
	bindingRepo *repository.BindingRepository
	settingRepo *repository.SettingRepository
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		userRepo:    repository.NewUserRepository(),
		clusterRepo: repository.NewClusterRepository(),
		bindingRepo: repository.NewBindingRepository(),
		settingRepo: repository.NewSettingRepository(),
	}
}

func (h *UserHandler) List(c *gin.Context) {
	users, err := h.userRepo.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) Create(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required,min=3,max=64"`
		Password string `json:"password" binding:"required,min=6"`
		Email    string `json:"email"`
		Role     string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !model.IsValidRole(req.Role) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	user := &model.User{
		Username:     req.Username,
		PasswordHash: string(hash),
		Email:        req.Email,
		Role:         req.Role,
	}
	if err := h.userRepo.Create(user); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "username already exists"})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) UpdateRole(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	user, err := h.userRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	var req struct {
		Role string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !model.IsValidRole(req.Role) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role"})
		return
	}

	user.Role = req.Role
	if err := h.userRepo.Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update role"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetRegistrationSetting(c *gin.Context) {
	regEnabled, _ := h.settingRepo.Get("registration_enabled")
	c.JSON(http.StatusOK, gin.H{"registration_enabled": regEnabled == "true"})
}

func (h *UserHandler) SetRegistrationSetting(c *gin.Context) {
	var req struct {
		Enabled bool `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	val := "false"
	if req.Enabled {
		val = "true"
	}
	if err := h.settingRepo.Set("registration_enabled", val, "Allow self-registration on login page"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update setting"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"registration_enabled": req.Enabled})
}

// BindCluster binds a cluster_admin user to a cluster.
func (h *UserHandler) BindCluster(c *gin.Context) {
	var req struct {
		UserID    int64 `json:"user_id" binding:"required"`
		ClusterID int64 `json:"cluster_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userRepo.GetByID(req.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	if model.NormalizeRole(user.Role) != model.RoleClusterAdmin {
		c.JSON(http.StatusBadRequest, gin.H{"error": "only cluster_admin users can be bound to clusters"})
		return
	}

	if h.bindingRepo.Exists(req.UserID, req.ClusterID) {
		c.JSON(http.StatusConflict, gin.H{"error": "binding already exists"})
		return
	}

	if err := h.bindingRepo.Create(&model.UserClusterBinding{UserID: req.UserID, ClusterID: req.ClusterID}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create binding"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "cluster bound"})
}

// UnbindCluster removes a user-cluster binding.
func (h *UserHandler) UnbindCluster(c *gin.Context) {
	var req struct {
		UserID    int64 `json:"user_id" binding:"required"`
		ClusterID int64 `json:"cluster_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.bindingRepo.Delete(req.UserID, req.ClusterID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete binding"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "cluster unbound"})
}

// Get returns a single user by ID.
func (h *UserHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	user, err := h.userRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// ResetPassword allows platform_admin to reset a user's password.
func (h *UserHandler) ResetPassword(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	user, err := h.userRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	var req struct {
		Password string `json:"password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	user.PasswordHash = string(hash)
	if err := h.userRepo.Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update password"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "password updated"})
}

// ToggleOTP allows platform_admin to enable or disable a user's OTP.
func (h *UserHandler) ToggleOTP(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	user, err := h.userRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	var req struct {
		Enabled bool `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.OTPEnabled = req.Enabled
	if !req.Enabled {
		user.OTPSecret = ""
	}
	if err := h.userRepo.Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update OTP"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "OTP updated", "otp_enabled": user.OTPEnabled})
}

// UpdateLanguage allows platform_admin to set a user's UI language.
func (h *UserHandler) UpdateLanguage(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	user, err := h.userRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	var req struct {
		Language string `json:"language" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.Language = req.Language
	if err := h.userRepo.Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update language"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "language updated", "language": user.Language})
}

// ListBindings returns clusters bound to a user.
func (h *UserHandler) ListBindings(c *gin.Context) {
	userID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	bindings, err := h.bindingRepo.ListByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list bindings"})
		return
	}
	c.JSON(http.StatusOK, bindings)
}
