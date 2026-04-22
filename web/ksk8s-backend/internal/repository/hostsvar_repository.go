package repository

import (
	"github.com/easzlab/ksk8s/internal/model"
)

type HostsVarRepository struct{}

func NewHostsVarRepository() *HostsVarRepository {
	return &HostsVarRepository{}
}

func (r *HostsVarRepository) GetOrCreate(clusterID int64) (*model.HostsGlobalVars, error) {
	var vars model.HostsGlobalVars
	err := DB.Where("cluster_id = ?", clusterID).First(&vars).Error
	if err != nil {
		vars.ClusterID = clusterID
		err = DB.Create(&vars).Error
	}
	return &vars, err
}

func (r *HostsVarRepository) Update(vars *model.HostsGlobalVars) error {
	return DB.Save(vars).Error
}

func (r *HostsVarRepository) DeleteByClusterID(clusterID int64) error {
	return DB.Where("cluster_id = ?", clusterID).Delete(&model.HostsGlobalVars{}).Error
}
