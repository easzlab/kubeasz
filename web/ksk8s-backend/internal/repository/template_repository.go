package repository

import (
	"github.com/easzlab/ksk8s/internal/model"
)

type TemplateRepository struct{}

func NewTemplateRepository() *TemplateRepository {
	return &TemplateRepository{}
}

func (r *TemplateRepository) Create(t *model.Template) error {
	return DB.Create(t).Error
}

func (r *TemplateRepository) GetByID(id int64) (*model.Template, error) {
	var t model.Template
	err := DB.First(&t, id).Error
	return &t, err
}

func (r *TemplateRepository) GetDefault() (*model.Template, error) {
	var t model.Template
	err := DB.Where("is_default = ?", true).First(&t).Error
	return &t, err
}

func (r *TemplateRepository) List() ([]model.Template, error) {
	var templates []model.Template
	err := DB.Find(&templates).Error
	return templates, err
}

func (r *TemplateRepository) Update(t *model.Template) error {
	return DB.Model(t).Updates(t).Error
}

func (r *TemplateRepository) Delete(id int64) error {
	return DB.Delete(&model.Template{}, id).Error
}

func (r *TemplateRepository) ClearDefault() error {
	return DB.Model(&model.Template{}).Where("is_default = ?", true).Update("is_default", false).Error
}
