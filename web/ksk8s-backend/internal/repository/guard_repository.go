package repository

import "github.com/easzlab/ksk8s/internal/model"

type GuardRepository struct{}

func NewGuardRepository() *GuardRepository {
	return &GuardRepository{}
}

// Acquire tries to insert a running_tasks row. Returns error if duplicate (cluster_id, step_number).
func (r *GuardRepository) Acquire(clusterID int64, stepNumber string, taskID int64) error {
	return DB.Create(&model.RunningTask{
		ClusterID:  clusterID,
		StepNumber: stepNumber,
		TaskID:     taskID,
	}).Error
}

// Release deletes the guard row for a cluster+step.
func (r *GuardRepository) Release(clusterID int64, stepNumber string) error {
	return DB.Where("cluster_id = ? AND step_number = ?", clusterID, stepNumber).Delete(&model.RunningTask{}).Error
}

// ReleaseByTaskID deletes all guard rows for a given task.
func (r *GuardRepository) ReleaseByTaskID(taskID int64) error {
	return DB.Where("task_id = ?", taskID).Delete(&model.RunningTask{}).Error
}

// GetByClusterStep returns the running task guard for a cluster+step, if any.
func (r *GuardRepository) GetByClusterStep(clusterID int64, stepNumber string) (*model.RunningTask, error) {
	var rt model.RunningTask
	err := DB.Where("cluster_id = ? AND step_number = ?", clusterID, stepNumber).First(&rt).Error
	return &rt, err
}
