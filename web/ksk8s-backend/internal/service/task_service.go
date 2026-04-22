package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/easzlab/ksk8s/internal/model"
	"github.com/easzlab/ksk8s/internal/repository"
	"github.com/easzlab/ksk8s/internal/websocket"
)

// TaskService orchestrates task lifecycle: start, abort, approve, retry.
type TaskService struct {
	taskRepo   *repository.TaskRepository
	guardRepo  *repository.GuardRepository
	clusterRepo *repository.ClusterRepository
	logRepo    *repository.LogRepository
	logRings   *websocket.LogRingMap
	sem        *Semaphore
	workers    map[int64]*Worker
}

// NewTaskService creates a TaskService with the shared log ring map.
func NewTaskService(logRings *websocket.LogRingMap) *TaskService {
	return &TaskService{
		taskRepo:    repository.NewTaskRepository(),
		guardRepo:   repository.NewGuardRepository(),
		clusterRepo: repository.NewClusterRepository(),
		logRepo:     repository.NewLogRepository(),
		logRings:    logRings,
		sem:         NewSemaphore(2),
		workers:     make(map[int64]*Worker),
	}
}

// StartTask creates and starts a new task for a cluster+step.
func (s *TaskService) StartTask(clusterID int64, step string, taskType string, userID int64, extraArgs ...string) (*model.Task, error) {
	cluster, err := s.clusterRepo.GetByID(clusterID)
	if err != nil {
		return nil, fmt.Errorf("cluster not found: %w", err)
	}

	// Check guard table
	if _, err := s.guardRepo.GetByClusterStep(clusterID, step); err == nil {
		return nil, fmt.Errorf("step %s is already running for cluster %d", step, clusterID)
	}

	task := &model.Task{
		ClusterID:  clusterID,
		TaskType:   taskType,
		StepNumber: &step,
		Status:     "pending",
		LogPath:    fmt.Sprintf("/var/log/ksk8s/%d/task_pending.log", clusterID),
	}
	if err := s.taskRepo.Create(task); err != nil {
		return nil, fmt.Errorf("create task: %w", err)
	}

	// Acquire semaphore (with timeout)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	if err := s.sem.Acquire(ctx); err != nil {
		task.Status = "failed"
		task.ErrorMessage = "could not acquire execution slot (max 2 concurrent tasks)"
		_ = s.taskRepo.Update(task)
		return task, err
	}

	// Set guard
	if err := s.guardRepo.Acquire(clusterID, step, task.ID); err != nil {
		s.sem.Release()
		task.Status = "failed"
		task.ErrorMessage = "duplicate task guard"
		_ = s.taskRepo.Update(task)
		return task, fmt.Errorf("guard acquire failed: %w", err)
	}

	// Create log writer
	logWriter, err := NewLogWriter(task.ID, clusterID, s.logRepo, s.logRings)
	if err != nil {
		_ = s.guardRepo.Release(clusterID, step)
		s.sem.Release()
		task.Status = "failed"
		task.ErrorMessage = err.Error()
		_ = s.taskRepo.Update(task)
		return task, fmt.Errorf("log writer: %w", err)
	}

	// Update task status
	now := time.Now()
	task.Status = "running"
	task.LogPath = logWriter.Path()
	task.StartedAt = &now
	_ = s.taskRepo.Update(task)

	// Write cluster config files to disk so ezctl can read them
	if os.Getenv("KSK8S_MOCK_EZCTL") != "1" {
		if err := WriteClusterConfigFiles(cluster); err != nil {
			_ = s.guardRepo.Release(clusterID, step)
			s.sem.Release()
			_ = logWriter.Close()
			task.Status = "failed"
			task.ErrorMessage = err.Error()
			_ = s.taskRepo.Update(task)
			return task, fmt.Errorf("write cluster config files: %w", err)
		}
	}

	// Build worker
	var worker *Worker
	var workerArgs []string
	workerArgs = append(workerArgs, extraArgs...)
	// Auto-skip cilium kernel check on older kernels (CentOS 7)
	if step == "06" {
		workerArgs = append(workerArgs, "-e", "skip_cilium_kernel_check=true")
	}
	if os.Getenv("KSK8S_MOCK_EZCTL") == "1" {
		worker = NewWorkerWithCommand(task.ID, clusterID, "/bin/bash", []string{"-c", mockScript(step)}, logWriter)
	} else {
		worker = NewWorker(task.ID, clusterID, cluster.Name, step, logWriter, workerArgs...)
	}

	s.workers[task.ID] = worker

	// Start worker
	if err := worker.Start(); err != nil {
		_ = s.guardRepo.Release(clusterID, step)
		s.sem.Release()
		_ = logWriter.Close()
		task.Status = "failed"
		task.ErrorMessage = err.Error()
		_ = s.taskRepo.Update(task)
		return task, fmt.Errorf("worker start: %w", err)
	}

	// Save worker PID for recovery
	pid := worker.PID()
	task.WorkerPID = &pid
	_ = s.taskRepo.Update(task)

	// Goroutine to wait for completion
	go s.watchWorker(worker, task, step)

	return task, nil
}

func (s *TaskService) watchWorker(worker *Worker, task *model.Task, step string) {
	<-worker.Done()

	now := time.Now()
	task.CompletedAt = &now

	// Check Ansible PLAY RECAP results from log output
	hosts, failed := worker.logWriter.AnsibleResult()
	exitCode := 0
	if worker.Err() != nil {
		exitCode = ParseExitCode(worker.Err())
	}
	task.ExitCode = &exitCode

	// Write explicit PLAY RECAP analysis to log before closing
	if hosts > 0 {
		var resultLine string
		if failed > 0 {
			resultLine = fmt.Sprintf("[ksk8s] PLAY RECAP analyzed: %d host(s), %d failure(s) -> FAILED", hosts, failed)
		} else {
			resultLine = fmt.Sprintf("[ksk8s] PLAY RECAP analyzed: %d host(s), 0 failure(s) -> SUCCESS", hosts)
		}
		_ = worker.logWriter.WriteLine(resultLine, "stdout")
	} else {
		_ = worker.logWriter.WriteLine("[ksk8s] No PLAY RECAP found in output", "stdout")
	}

	_ = worker.logWriter.Close()

	if worker.Err() != nil {
		task.Status = "failed"
		task.ErrorMessage = worker.Err().Error()
	} else if failed > 0 {
		task.Status = "failed"
		task.ErrorMessage = fmt.Sprintf("ansible PLAY RECAP shows %d failure(s) across %d host(s)", failed, hosts)
	} else {
		task.Status = "success"
	}

	_ = s.taskRepo.Update(task)
	_ = s.guardRepo.Release(task.ClusterID, step)
	s.sem.Release()
	delete(s.workers, task.ID)
}

// AbortTask kills the worker process group and marks task aborted.
func (s *TaskService) AbortTask(taskID int64) error {
	worker, ok := s.workers[taskID]
	if !ok {
		return fmt.Errorf("task %d is not running", taskID)
	}

	if err := worker.Kill(); err != nil {
		return fmt.Errorf("kill worker: %w", err)
	}

	task, err := s.taskRepo.GetByID(taskID)
	if err != nil {
		return err
	}

	task.Status = "aborted"
	code := -1
	task.ExitCode = &code
	now := time.Now()
	task.CompletedAt = &now
	_ = s.taskRepo.Update(task)

	_ = s.guardRepo.ReleaseByTaskID(taskID)
	delete(s.workers, taskID)
	return nil
}

// ApproveTask marks a task as approved and advances to next step if applicable.
func (s *TaskService) ApproveTask(taskID int64, userID int64) error {
	task, err := s.taskRepo.GetByID(taskID)
	if err != nil {
		return err
	}
	if task.Status != "awaiting_approval" {
		return fmt.Errorf("task is not awaiting approval (status=%s)", task.Status)
	}

	task.Status = "success"
	now := time.Now()
	task.ApprovedAt = &now
	if userID > 0 {
		task.ApprovedBy = &userID
	}
	return s.taskRepo.Update(task)
}

// RetryTask creates a new task for the same cluster+step.
func (s *TaskService) RetryTask(clusterID int64, step string, taskType string, userID int64) (*model.Task, error) {
	return s.StartTask(clusterID, step, taskType, userID)
}

// GetTask returns a task by ID.
func (s *TaskService) GetTask(taskID int64) (*model.Task, error) {
	return s.taskRepo.GetByID(taskID)
}

// ListTasksByCluster returns all tasks for a cluster.
func (s *TaskService) ListTasksByCluster(clusterID int64) ([]model.Task, error) {
	return s.taskRepo.ListByCluster(clusterID)
}

// NewLogRepo returns a new LogRepository (convenience for handlers).
func (s *TaskService) NewLogRepo() *repository.LogRepository {
	return repository.NewLogRepository()
}

// GetWorker returns the running worker for a task, if any.
func (s *TaskService) GetWorker(taskID int64) *Worker {
	return s.workers[taskID]
}

// GetLogRing returns the log ring for a task.
func (s *TaskService) GetLogRing(taskID int64) *websocket.LogRing {
	return s.logRings.Get(taskID)
}

// RunningWorkers returns a snapshot of current workers.
func (s *TaskService) RunningWorkers() map[int64]*Worker {
	result := make(map[int64]*Worker)
	for k, v := range s.workers {
		result[k] = v
	}
	return result
}

// RecoverOnStartup scans for tasks that were running before a backend restart.
// If the worker PID is still alive, it re-attaches (best effort).
// Otherwise, marks the task as failed and releases guards.
func (s *TaskService) RecoverOnStartup() {
	// Find all tasks in running or awaiting_approval state
	var tasks []model.Task
	_ = repository.DB.Where("status IN ?", []string{"running", "awaiting_approval"}).Find(&tasks)
	for _, task := range tasks {
		if task.WorkerPID != nil && Reattach(*task.WorkerPID) {
			log.Printf("[recovery] re-attached to task %d (pid %d)", task.ID, *task.WorkerPID)
			// We can't fully reconstruct the worker object without stdout/stderr pipes,
			// but we mark it as still running. The user can abort and retry if needed.
			continue
		}
		// PID dead or missing — mark failed
		task.Status = "failed"
		task.ErrorMessage = "backend restarted while task was running"
		code := -1
		task.ExitCode = &code
		now := time.Now()
		task.CompletedAt = &now
		_ = s.taskRepo.Update(&task)
		_ = s.guardRepo.ReleaseByTaskID(task.ID)
		log.Printf("[recovery] marked task %d as failed (backend restart)", task.ID)
	}
}

func mockScript(step string) string {
	return fmt.Sprintf(`
echo "[MOCK] Starting step %s..."
for i in {1..20}; do
  echo "[MOCK] Progress $i/20"
  sleep 0.3
done
echo "[MOCK] Step %s completed successfully"
`, step, step)
}

// WriteClusterConfigFiles writes a cluster's hosts and config content to
// /etc/kubeasz/clusters/<name>/ so ezctl can read them.
func WriteClusterConfigFiles(cluster *model.Cluster) error {
	baseDir := os.Getenv("KSK8S_KUBEASZ_DIR")
	if baseDir == "" {
		baseDir = "/etc/kubeasz"
	}
	clusterDir := filepath.Join(baseDir, "clusters", cluster.Name)

	if err := os.MkdirAll(clusterDir, 0755); err != nil {
		return fmt.Errorf("create cluster dir %s: %w", clusterDir, err)
	}

	hostsPath := filepath.Join(clusterDir, "hosts")
	if err := os.WriteFile(hostsPath, []byte(cluster.HostsContent), 0644); err != nil {
		return fmt.Errorf("write hosts file: %w", err)
	}

	configPath := filepath.Join(clusterDir, "config.yml")
	if err := os.WriteFile(configPath, []byte(cluster.ConfigContent), 0644); err != nil {
		return fmt.Errorf("write config.yml: %w", err)
	}

	return nil
}
