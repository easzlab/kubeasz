package repository

import "github.com/easzlab/ksk8s/internal/model"

type BindingRepository struct{}

func NewBindingRepository() *BindingRepository {
	return &BindingRepository{}
}

func (r *BindingRepository) ListByUser(userID int64) ([]model.UserClusterBinding, error) {
	var bindings []model.UserClusterBinding
	err := DB.Where("user_id = ?", userID).Find(&bindings).Error
	return bindings, err
}

func (r *BindingRepository) ListByCluster(clusterID int64) ([]model.UserClusterBinding, error) {
	var bindings []model.UserClusterBinding
	err := DB.Where("cluster_id = ?", clusterID).Find(&bindings).Error
	return bindings, err
}

func (r *BindingRepository) Create(binding *model.UserClusterBinding) error {
	return DB.Create(binding).Error
}

func (r *BindingRepository) Delete(userID, clusterID int64) error {
	return DB.Where("user_id = ? AND cluster_id = ?", userID, clusterID).Delete(&model.UserClusterBinding{}).Error
}

func (r *BindingRepository) Exists(userID, clusterID int64) bool {
	var count int64
	DB.Model(&model.UserClusterBinding{}).Where("user_id = ? AND cluster_id = ?", userID, clusterID).Count(&count)
	return count > 0
}
