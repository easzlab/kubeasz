package model

import "time"

type Task struct {
	ID            int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	ClusterID     int64      `gorm:"not null;index:idx_cluster_status" json:"cluster_id"`
	TaskType      string     `gorm:"size:16;not null" json:"task_type"`
	StepNumber    *string    `gorm:"size:8" json:"step_number,omitempty"`
	TargetNodeIP  *string    `gorm:"size:64" json:"target_node_ip,omitempty"`
	Status        string     `gorm:"size:16;default:'pending'" json:"status"`
	WorkerPID     *int       `json:"worker_pid,omitempty"`
	LogPath       string     `gorm:"size:512" json:"log_path"`
	StartedAt     *time.Time `json:"started_at,omitempty"`
	CompletedAt   *time.Time `json:"completed_at,omitempty"`
	ApprovedBy    *int64     `json:"approved_by,omitempty"`
	ApprovedAt    *time.Time `json:"approved_at,omitempty"`
	ExitCode      *int       `json:"exit_code,omitempty"`
	ErrorMessage  string     `gorm:"type:text" json:"error_message,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
}

type Log struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	TaskID      int64     `gorm:"not null;index:idx_task_line" json:"task_id"`
	LineNumber  int       `gorm:"not null;index:idx_task_line" json:"line_number"`
	Content     string    `gorm:"type:text;not null" json:"content"`
	Timestamp   time.Time `json:"timestamp"`
	Stream      string    `gorm:"size:8;default:'stdout'" json:"stream"`
}

type RunningTask struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ClusterID   int64     `gorm:"not null;uniqueIndex:uk_cluster_step" json:"cluster_id"`
	StepNumber  string    `gorm:"size:8;not null;uniqueIndex:uk_cluster_step" json:"step_number"`
	TaskID      int64     `gorm:"not null" json:"task_id"`
	CreatedAt   time.Time `json:"created_at"`
}
