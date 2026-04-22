package repository

import (
	"github.com/easzlab/ksk8s/internal/model"
)

type SettingRepository struct{}

func NewSettingRepository() *SettingRepository {
	return &SettingRepository{}
}

func (r *SettingRepository) Get(key string) (string, error) {
	var s model.SystemSetting
	err := DB.Where("`key` = ?", key).First(&s).Error
	if err != nil {
		return "", err
	}
	return s.Value, nil
}

func (r *SettingRepository) Set(key string, value string, description string) error {
	var s model.SystemSetting
	err := DB.Where("`key` = ?", key).First(&s).Error
	if err != nil {
		// create
		s = model.SystemSetting{Key: key, Value: value, Description: description}
		return DB.Create(&s).Error
	}
	s.Value = value
	return DB.Save(&s).Error
}

func (r *SettingRepository) InitDefaults() {
	// Ensure registration setting exists
	var count int64
	DB.Model(&model.SystemSetting{}).Where("`key` = ?", "registration_enabled").Count(&count)
	if count == 0 {
		_ = r.Set("registration_enabled", "false", "Allow self-registration on login page")
	}
}
