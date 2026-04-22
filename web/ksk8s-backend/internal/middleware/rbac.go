package middleware

import (
	"net/http"
	"strconv"

	"github.com/easzlab/ksk8s/internal/model"
	"github.com/easzlab/ksk8s/internal/repository"
	"github.com/gin-gonic/gin"
)

func RequirePlatformAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		if model.NormalizeRole(role.(string)) != model.RolePlatformAdmin {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "platform_admin access required"})
			return
		}
		c.Next()
	}
}

func RequireClusterAdminOrAbove() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		r := model.NormalizeRole(role.(string))
		if r != model.RolePlatformAdmin && r != model.RoleClusterAdmin {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "cluster_admin access required"})
			return
		}
		c.Next()
	}
}

func RequireAuditorOrAbove() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		r := model.NormalizeRole(role.(string))
		if r != model.RolePlatformAdmin && r != model.RoleClusterAdmin && r != model.RoleSecurityAuditor {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denied"})
			return
		}
		c.Next()
	}
}

// RequireClusterAccess ensures the user can access a specific cluster.
// platform_admin: all clusters
// cluster_admin: own clusters + bound clusters
// security_auditor: all clusters (readonly)
func RequireClusterAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		r := model.NormalizeRole(role.(string))
		if r == model.RolePlatformAdmin {
			c.Next()
			return
		}

		userID, _ := c.Get("user_id")
		uid := userID.(int64)
		clusterID := parseClusterID(c)

		clusterRepo := NewClusterRepository()
		bindingRepo := NewBindingRepository()

		// cluster_admin: own or bound
		if r == model.RoleClusterAdmin {
			if clusterRepo.IsOwner(clusterID, uid) || bindingRepo.Exists(uid, clusterID) {
				c.Next()
				return
			}
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "cluster access denied"})
			return
		}

		// security_auditor: can view all clusters (read-only handled by route-level middleware)
		if r == model.RoleSecurityAuditor {
			c.Next()
			return
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "cluster access denied"})
	}
}

// RequireClusterWrite ensures the user can write to a specific cluster.
// security_auditor is blocked.
func RequireClusterWrite() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		r := model.NormalizeRole(role.(string))
		if r == model.RolePlatformAdmin {
			c.Next()
			return
		}
		if r == model.RoleSecurityAuditor {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "read-only access"})
			return
		}

		userID, _ := c.Get("user_id")
		uid := userID.(int64)
		clusterID := parseClusterID(c)

		clusterRepo := NewClusterRepository()
		bindingRepo := NewBindingRepository()

		if clusterRepo.IsOwner(clusterID, uid) || bindingRepo.Exists(uid, clusterID) {
			c.Next()
			return
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "cluster write access denied"})
	}
}

func parseClusterID(c *gin.Context) int64 {
	idStr := c.Param("id")
	if idStr == "" {
		idStr = c.Param("clusterId")
	}
	id, _ := strconv.ParseInt(idStr, 10, 64)
	return id
}

// --- helpers to avoid import cycle ---

type clusterRepo struct{}

func NewClusterRepository() *clusterRepo { return &clusterRepo{} }

func (r *clusterRepo) IsOwner(clusterID, userID int64) bool {
	var count int64
	repository.DB.Model(&model.Cluster{}).Where("id = ? AND created_by = ?", clusterID, userID).Count(&count)
	return count > 0
}

type bindingRepo struct{}

func NewBindingRepository() *bindingRepo { return &bindingRepo{} }

func (r *bindingRepo) Exists(userID, clusterID int64) bool {
	var count int64
	repository.DB.Model(&model.UserClusterBinding{}).Where("user_id = ? AND cluster_id = ?", userID, clusterID).Count(&count)
	return count > 0
}
