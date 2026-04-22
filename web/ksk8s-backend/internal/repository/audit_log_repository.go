package repository

import (
	"time"

	"github.com/easzlab/ksk8s/internal/model"
	"gorm.io/gorm"
)

type AuditLogRepository struct{}

func NewAuditLogRepository() *AuditLogRepository {
	return &AuditLogRepository{}
}

func (r *AuditLogRepository) Create(log *model.Audit) error {
	return DB.Create(log).Error
}

func (r *AuditLogRepository) List(offset int, limit int) ([]model.Audit, error) {
	var logs []model.Audit
	err := DB.Order("created_at DESC").Offset(offset).Limit(limit).Find(&logs).Error
	return logs, err
}

func (r *AuditLogRepository) Count() (int64, error) {
	var count int64
	err := DB.Model(&model.Audit{}).Count(&count).Error
	return count, err
}

func (r *AuditLogRepository) ListByUser(userID int64, offset int, limit int) ([]model.Audit, error) {
	var logs []model.Audit
	err := DB.Where("user_id = ?", userID).Order("created_at DESC").Offset(offset).Limit(limit).Find(&logs).Error
	return logs, err
}

// ListFiltered returns audit logs with filtering options.
func (r *AuditLogRepository) ListFiltered(query *AuditQuery, offset int, limit int) ([]model.Audit, error) {
	var logs []model.Audit
	db := buildAuditQuery(query)
	err := db.Order("created_at DESC").Offset(offset).Limit(limit).Find(&logs).Error
	return logs, err
}

// CountFiltered returns total count with filtering options.
func (r *AuditLogRepository) CountFiltered(query *AuditQuery) (int64, error) {
	var count int64
	db := buildAuditQuery(query)
	err := db.Model(&model.Audit{}).Count(&count).Error
	return count, err
}

// AuditQuery holds filter parameters.
type AuditQuery struct {
	StartTime  *time.Time
	EndTime    *time.Time
	Action     string
	Username   string
	IsHighRisk *bool
}

func buildAuditQuery(q *AuditQuery) *gorm.DB {
	db := DB
	if q == nil {
		return db
	}
	if q.StartTime != nil {
		db = db.Where("created_at >= ?", *q.StartTime)
	}
	if q.EndTime != nil {
		db = db.Where("created_at <= ?", *q.EndTime)
	}
	if q.Action != "" {
		db = db.Where("action = ?", q.Action)
	}
	if q.Username != "" {
		db = db.Where("username = ?", q.Username)
	}
	if q.IsHighRisk != nil {
		db = db.Where("is_high_risk = ?", *q.IsHighRisk)
	}
	return db
}
