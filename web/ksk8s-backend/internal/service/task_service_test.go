package service

import (
	"context"
	"testing"
	"time"

	"github.com/easzlab/ksk8s/internal/model"
	"github.com/easzlab/ksk8s/internal/repository"
	"github.com/easzlab/ksk8s/internal/websocket"
)

func TestSemaphore(t *testing.T) {
	sem := NewSemaphore(1)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	if err := sem.Acquire(ctx); err != nil {
		t.Fatalf("first acquire should succeed: %v", err)
	}

	ctx2, cancel2 := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel2()

	if err := sem.Acquire(ctx2); err == nil {
		t.Error("second acquire should fail when capacity is 1")
	}

	sem.Release()

	ctx3, cancel3 := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel3()

	if err := sem.Acquire(ctx3); err != nil {
		t.Fatalf("acquire after release should succeed: %v", err)
	}
}

func TestLogRing(t *testing.T) {
	ring := websocket.NewLogRing(5)

	for i := 1; i <= 10; i++ {
		ring.Append(websocket.LogLine{LineNumber: i, Content: "line " + string(rune('0'+i))})
	}

	if ring.Total() != 10 {
		t.Errorf("expected total 10, got %d", ring.Total())
	}

	lines, total := ring.Since(7)
	if total != 10 {
		t.Errorf("expected total 10, got %d", total)
	}
	if len(lines) != 3 {
		t.Errorf("expected 3 lines since 7, got %d", len(lines))
	}

	// Test that old lines are overwritten
	lines, _ = ring.Since(0)
	if len(lines) != 5 {
		t.Errorf("expected 5 lines in ring (capacity), got %d", len(lines))
	}
}

func setupTaskServiceTest(t *testing.T) (*TaskService, *model.Cluster) {
	db := repository.SetupTestDB(t)
	_ = db

	t.Setenv("KSK8S_LOG_DIR", t.TempDir())
	t.Setenv("KSK8S_MOCK_EZCTL", "1")

	clusterRepo := repository.NewClusterRepository()
	cluster := &model.Cluster{Name: t.Name(), CreatedBy: 1, Status: "draft"}
	if err := clusterRepo.Create(cluster); err != nil {
		t.Fatalf("failed to create cluster: %v", err)
	}

	svc := NewTaskService(websocket.NewLogRingMap())
	return svc, cluster
}

func TestTaskService_StartTask(t *testing.T) {
	svc, cluster := setupTaskServiceTest(t)

	task, err := svc.StartTask(cluster.ID, "01", "setup", 1)
	if err != nil {
		t.Fatalf("StartTask failed: %v", err)
	}
	if task.Status != "running" {
		t.Errorf("expected status running, got %s", task.Status)
	}
	if task.LogPath == "" {
		t.Error("expected log path to be set")
	}

	// Clean up: abort the running worker
	_ = svc.AbortTask(task.ID)
}

func TestTaskService_StartTask_DuplicateGuard(t *testing.T) {
	svc, cluster := setupTaskServiceTest(t)

	task1, err := svc.StartTask(cluster.ID, "01", "setup", 1)
	if err != nil {
		t.Fatalf("first StartTask failed: %v", err)
	}

	// Second start for same cluster+step should fail
	_, err = svc.StartTask(cluster.ID, "01", "setup", 1)
	if err == nil {
		t.Error("expected duplicate guard error")
	}

	_ = svc.AbortTask(task1.ID)
}

func TestTaskService_AbortTask(t *testing.T) {
	svc, cluster := setupTaskServiceTest(t)

	task, err := svc.StartTask(cluster.ID, "02", "setup", 1)
	if err != nil {
		t.Fatalf("StartTask failed: %v", err)
	}

	// Give worker a moment to start
	time.Sleep(100 * time.Millisecond)

	if err := svc.AbortTask(task.ID); err != nil {
		t.Fatalf("AbortTask failed: %v", err)
	}

	// Verify task status updated
	updated, _ := svc.GetTask(task.ID)
	if updated.Status != "aborted" {
		t.Errorf("expected status aborted, got %s", updated.Status)
	}
	if updated.ExitCode == nil || *updated.ExitCode != -1 {
		t.Error("expected exit code -1 for aborted task")
	}
}

func TestTaskService_ApproveTask(t *testing.T) {
	svc, cluster := setupTaskServiceTest(t)

	task := &model.Task{
		ClusterID:  cluster.ID,
		TaskType:   "setup",
		StepNumber: strPtr("03"),
		Status:     "awaiting_approval",
	}
	repository.NewTaskRepository().Create(task)

	if err := svc.ApproveTask(task.ID, 1); err != nil {
		t.Fatalf("ApproveTask failed: %v", err)
	}

	updated, _ := svc.GetTask(task.ID)
	if updated.Status != "success" {
		t.Errorf("expected status success, got %s", updated.Status)
	}
	if updated.ApprovedBy == nil || *updated.ApprovedBy != 1 {
		t.Error("expected approved_by to be set")
	}
}

func TestTaskService_ApproveTask_WrongStatus(t *testing.T) {
	svc, cluster := setupTaskServiceTest(t)

	task := &model.Task{
		ClusterID:  cluster.ID,
		TaskType:   "setup",
		StepNumber: strPtr("04"),
		Status:     "running",
	}
	repository.NewTaskRepository().Create(task)

	err := svc.ApproveTask(task.ID, 1)
	if err == nil {
		t.Error("expected error when approving non-awaiting task")
	}
}

func TestTaskService_GetTask(t *testing.T) {
	svc, cluster := setupTaskServiceTest(t)

	task := &model.Task{
		ClusterID:  cluster.ID,
		TaskType:   "setup",
		StepNumber: strPtr("05"),
		Status:     "pending",
	}
	repository.NewTaskRepository().Create(task)

	found, err := svc.GetTask(task.ID)
	if err != nil {
		t.Fatalf("GetTask failed: %v", err)
	}
	if found.ID != task.ID {
		t.Errorf("expected task ID %d, got %d", task.ID, found.ID)
	}
}

func TestTaskService_ListTasksByCluster(t *testing.T) {
	svc, cluster := setupTaskServiceTest(t)

	repo := repository.NewTaskRepository()
	repo.Create(&model.Task{ClusterID: cluster.ID, TaskType: "setup", Status: "pending"})
	repo.Create(&model.Task{ClusterID: cluster.ID, TaskType: "setup", Status: "success"})

	tasks, err := svc.ListTasksByCluster(cluster.ID)
	if err != nil {
		t.Fatalf("ListTasksByCluster failed: %v", err)
	}
	if len(tasks) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(tasks))
	}
}

func TestTaskService_RecoverOnStartup(t *testing.T) {
	svc, cluster := setupTaskServiceTest(t)

	// Simulate a task that was running before restart
	repo := repository.NewTaskRepository()
	task := &model.Task{
		ClusterID:  cluster.ID,
		TaskType:   "setup",
		StepNumber: strPtr("06"),
		Status:     "running",
	}
	repo.Create(task)

	// Also create a guard entry
	guardRepo := repository.NewGuardRepository()
	_ = guardRepo.Acquire(cluster.ID, "06", task.ID)

	svc.RecoverOnStartup()

	updated, _ := repo.GetByID(task.ID)
	if updated.Status != "failed" {
		t.Errorf("expected status failed after recovery, got %s", updated.Status)
	}
	if updated.ErrorMessage != "backend restarted while task was running" {
		t.Errorf("expected recovery error message, got %s", updated.ErrorMessage)
	}

	// Guard should be released
	_, err := guardRepo.GetByClusterStep(cluster.ID, "06")
	if err == nil {
		t.Error("expected guard to be released after recovery")
	}
}

func TestTaskService_RetryTask(t *testing.T) {
	svc, cluster := setupTaskServiceTest(t)

	task, err := svc.RetryTask(cluster.ID, "07", "setup", 1)
	if err != nil {
		t.Fatalf("RetryTask failed: %v", err)
	}
	if task.Status != "running" {
		t.Errorf("expected status running, got %s", task.Status)
	}

	_ = svc.AbortTask(task.ID)
}

func TestTaskService_RunningWorkers(t *testing.T) {
	svc, cluster := setupTaskServiceTest(t)

	task, err := svc.StartTask(cluster.ID, "08", "setup", 1)
	if err != nil {
		t.Fatalf("StartTask failed: %v", err)
	}

	workers := svc.RunningWorkers()
	if len(workers) != 1 {
		t.Errorf("expected 1 running worker, got %d", len(workers))
	}
	if workers[task.ID] == nil {
		t.Error("expected worker for task to be present")
	}

	_ = svc.AbortTask(task.ID)
}

func strPtr(s string) *string {
	return &s
}

func TestTaskService_NewLogRepo(t *testing.T) {
	svc := NewTaskService(websocket.NewLogRingMap())
	if svc.NewLogRepo() == nil {
		t.Error("expected non-nil log repo")
	}
}

func TestMockScript(t *testing.T) {
	script := mockScript("01")
	if script == "" {
		t.Error("expected non-empty mock script")
	}
}
