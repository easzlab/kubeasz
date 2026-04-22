package repository

import (
	"github.com/easzlab/ksk8s/internal/model"
)

type NodeRepository struct{}

func NewNodeRepository() *NodeRepository {
	return &NodeRepository{}
}

func (r *NodeRepository) Create(node *model.ClusterNode) error {
	return DB.Create(node).Error
}

func (r *NodeRepository) GetByClusterID(clusterID int64) ([]model.ClusterNode, error) {
	var nodes []model.ClusterNode
	err := DB.Where("cluster_id = ?", clusterID).Order("sort_order asc").Find(&nodes).Error
	return nodes, err
}

func (r *NodeRepository) Update(node *model.ClusterNode) error {
	return DB.Save(node).Error
}

func (r *NodeRepository) Delete(id int64) error {
	return DB.Delete(&model.ClusterNode{}, id).Error
}

func (r *NodeRepository) DeleteByClusterID(clusterID int64) error {
	return DB.Where("cluster_id = ?", clusterID).Delete(&model.ClusterNode{}).Error
}

func (r *NodeRepository) GetMaxSortOrder(clusterID int64, group string) (int, error) {
	var maxOrder int
	err := DB.Model(&model.ClusterNode{}).
		Where("cluster_id = ? AND group_name = ?", clusterID, group).
		Select("COALESCE(MAX(sort_order), 0)").
		Scan(&maxOrder).Error
	return maxOrder, err
}

func (r *NodeRepository) DeleteByClusterGroupIP(clusterID int64, group string, ip string) error {
	return DB.Where("cluster_id = ? AND group_name = ? AND ip_address = ?", clusterID, group, ip).
		Delete(&model.ClusterNode{}).Error
}
