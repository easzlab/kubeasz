package model

import "time"

const (
	RolePlatformAdmin   = "platform_admin"
	RoleClusterAdmin    = "cluster_admin"
	RoleSecurityAuditor = "security_auditor"
)

var ValidRoles = []string{RolePlatformAdmin, RoleClusterAdmin, RoleSecurityAuditor}

func IsValidRole(role string) bool {
	for _, r := range ValidRoles {
		if r == role {
			return true
		}
	}
	return false
}

// NormalizeRole maps legacy role names to current ones.
// admin -> platform_admin, viewer -> security_auditor.
func NormalizeRole(role string) string {
	switch role {
	case "admin":
		return RolePlatformAdmin
	case "viewer":
		return RoleSecurityAuditor
	}
	return role
}

type User struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Username     string    `gorm:"size:64;not null;uniqueIndex" json:"username"`
	PasswordHash string    `gorm:"size:255;not null" json:"-"`
	Email        string    `gorm:"size:128" json:"email"`
	Role         string    `gorm:"size:32;default:'cluster_admin'" json:"role"`
	OTPSecret    string    `gorm:"size:64" json:"-"`
	OTPEnabled   bool      `gorm:"default:false" json:"otp_enabled"`
	Language     string    `gorm:"size:16;default:'en'" json:"language"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (u *User) IsPlatformAdmin() bool   { return NormalizeRole(u.Role) == RolePlatformAdmin }
func (u *User) IsClusterAdmin() bool    { return NormalizeRole(u.Role) == RoleClusterAdmin }
func (u *User) IsSecurityAuditor() bool { return NormalizeRole(u.Role) == RoleSecurityAuditor }

// UserClusterBinding binds a cluster_admin user to additional clusters.
type UserClusterBinding struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64     `gorm:"not null;uniqueIndex:uk_user_cluster" json:"user_id"`
	ClusterID int64     `gorm:"not null;uniqueIndex:uk_user_cluster" json:"cluster_id"`
	CreatedAt time.Time `json:"created_at"`
}

// SystemSetting holds platform-wide configuration flags.
type SystemSetting struct {
	ID                   int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Key                  string    `gorm:"size:64;not null;uniqueIndex" json:"key"`
	Value                string    `gorm:"size:256;not null" json:"value"`
	Description          string    `gorm:"size:256" json:"description"`
	UpdatedAt            time.Time `json:"updated_at"`
}
