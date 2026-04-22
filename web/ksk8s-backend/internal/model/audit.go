package model

import "time"

type Audit struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       int64     `gorm:"not null;index" json:"user_id"`
	Username     string    `gorm:"size:64" json:"username"`
	Action       string    `gorm:"size:64;not null;index" json:"action"`
	ResourceType string    `gorm:"size:32" json:"resource_type,omitempty"`
	ResourceID   string    `gorm:"size:64" json:"resource_id,omitempty"`
	Details      string    `gorm:"type:json" json:"details,omitempty"`
	IPAddress    string    `gorm:"size:64" json:"ip_address,omitempty"`
	StatusCode   int       `json:"status_code"`
	IsHighRisk   bool      `gorm:"default:false;index" json:"is_high_risk"`
	CreatedAt    time.Time `json:"created_at"`
}
