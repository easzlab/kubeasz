package repository

import (
	"github.com/easzlab/ksk8s/internal/model"
)

type TaskRepository struct{}

func NewTaskRepository() *TaskRepository {
	return &TaskRepository{}
}

func (r *TaskRepository) Create(task *model.Task) error {
	return DB.Create(task).Error
}

func (r *TaskRepository) GetByID(id int64) (*model.Task, error) {
	var task model.Task
	err := DB.First(&task, id).Error
	return &task, err
}

func (r *TaskRepository) GetRunningByCluster(clusterID int64) ([]model.Task, error) {
	var tasks []model.Task
	err := DB.Where("cluster_id = ? AND status IN ?", clusterID, []string{"running", "awaiting_approval"}).Find(&tasks).Error
	return tasks, err
}

func (r *TaskRepository) Update(task *model.Task) error {
	return DB.Save(task).Error
}

func (r *TaskRepository) ListByCluster(clusterID int64) ([]model.Task, error) {
	var tasks []model.Task
	err := DB.Where("cluster_id = ?", clusterID).Order("created_at desc").Find(&tasks).Error
	return tasks, err
}
