package repository

import (
	"github.com/easzlab/ksk8s/internal/model"
)

type ParamRepository struct{}

func NewParamRepository() *ParamRepository {
	return &ParamRepository{}
}

func (r *ParamRepository) Create(param *model.ClusterParam) error {
	return DB.Create(param).Error
}

func (r *ParamRepository) GetByClusterID(clusterID int64) ([]model.ClusterParam, error) {
	var params []model.ClusterParam
	err := DB.Where("cluster_id = ?", clusterID).Find(&params).Error
	return params, err
}

func (r *ParamRepository) Update(param *model.ClusterParam) error {
	return DB.Save(param).Error
}

func (r *ParamRepository) DeleteByClusterID(clusterID int64) error {
	return DB.Where("cluster_id = ?", clusterID).Delete(&model.ClusterParam{}).Error
}

func (r *ParamRepository) CreateListItem(item *model.ClusterParamList) error {
	return DB.Create(item).Error
}

func (r *ParamRepository) GetListByClusterID(clusterID int64) ([]model.ClusterParamList, error) {
	var items []model.ClusterParamList
	err := DB.Where("cluster_id = ?", clusterID).Order("sort_order asc").Find(&items).Error
	return items, err
}

func (r *ParamRepository) DeleteListByClusterID(clusterID int64) error {
	return DB.Where("cluster_id = ?", clusterID).Delete(&model.ClusterParamList{}).Error
}
