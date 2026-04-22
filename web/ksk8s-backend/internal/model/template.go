package model

import "time"

type Template struct {
	ID             int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name           string    `gorm:"size:64;not null;uniqueIndex" json:"name"`
	Description    string    `gorm:"size:512" json:"description"`
	IsDefault      bool      `gorm:"default:false" json:"is_default"`
	HostsContent   string    `gorm:"type:longtext;not null" json:"hosts_content"`
	ConfigContent  string    `gorm:"type:longtext;not null" json:"config_content"`
	CreatedBy      int64     `gorm:"not null" json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
