package repository

import (
	"github.com/easzlab/ksk8s/internal/model"
	"gorm.io/gorm"
)

type ClusterRepository struct{}

func NewClusterRepository() *ClusterRepository {
	return &ClusterRepository{}
}

func (r *ClusterRepository) Create(cluster *model.Cluster) error {
	return DB.Create(cluster).Error
}

func (r *ClusterRepository) GetByID(id int64) (*model.Cluster, error) {
	var cluster model.Cluster
	err := DB.Preload("Nodes").Preload("Params").Preload("ParamLists").Preload("GlobalVars").First(&cluster, id).Error
	return &cluster, err
}

func (r *ClusterRepository) List() ([]model.Cluster, error) {
	var clusters []model.Cluster
	err := DB.Find(&clusters).Error
	return clusters, err
}

func (r *ClusterRepository) ListByUser(userID int64) ([]model.Cluster, error) {
	var clusters []model.Cluster
	err := DB.Where("created_by = ?", userID).Find(&clusters).Error
	return clusters, err
}

func (r *ClusterRepository) ListByIDs(ids []int64) ([]model.Cluster, error) {
	if len(ids) == 0 {
		return []model.Cluster{}, nil
	}
	var clusters []model.Cluster
	err := DB.Where("id IN ?", ids).Find(&clusters).Error
	return clusters, err
}

func (r *ClusterRepository) IsOwner(clusterID, userID int64) bool {
	var count int64
	DB.Model(&model.Cluster{}).Where("id = ? AND created_by = ?", clusterID, userID).Count(&count)
	return count > 0
}

func (r *ClusterRepository) Update(cluster *model.Cluster) error {
	return DB.Model(&model.Cluster{}).Where("id = ?", cluster.ID).Updates(cluster).Error
}

func (r *ClusterRepository) Delete(id int64) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("cluster_id = ?", id).Delete(&model.ClusterNode{}).Error; err != nil {
			return err
		}
		if err := tx.Where("cluster_id = ?", id).Delete(&model.ClusterParam{}).Error; err != nil {
			return err
		}
		if err := tx.Where("cluster_id = ?", id).Delete(&model.ClusterParamList{}).Error; err != nil {
			return err
		}
		if err := tx.Where("cluster_id = ?", id).Delete(&model.HostsGlobalVars{}).Error; err != nil {
			return err
		}
		if err := tx.Where("cluster_id = ?", id).Delete(&model.ClusterVersion{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.Cluster{}, id).Error
	})
}
