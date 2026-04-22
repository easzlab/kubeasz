package service

import (
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/easzlab/ksk8s/internal/model"
	"github.com/easzlab/ksk8s/internal/repository"
	"github.com/easzlab/ksk8s/internal/websocket"
)

func TestWorker_MockCommand(t *testing.T) {
	db := repository.SetupTestDB(t)
	_ = db

	t.Setenv("KSK8S_LOG_DIR", t.TempDir())

	logRepo := repository.NewLogRepository()
	rings := websocket.NewLogRingMap()

	logWriter, err := NewLogWriter(1, 1, logRepo, rings)
	if err != nil {
		t.Fatalf("failed to create log writer: %v", err)
	}
	defer logWriter.Close()

	// Use a simple echo command as mock
	worker := NewWorkerWithCommand(1, 1, "/bin/echo", []string{"hello", "from", "worker"}, logWriter)

	if err := worker.Start(); err != nil {
		t.Fatalf("failed to start worker: %v", err)
	}

	// Wait for completion
	done := make(chan struct{})
	go func() {
		worker.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Expected
	case <-time.After(5 * time.Second):
		t.Fatal("worker did not complete in time")
	}

	if worker.Err() != nil {
		t.Fatalf("worker returned error: %v", worker.Err())
	}

	if worker.PID() == 0 {
		t.Error("worker PID should be set")
	}
}

func TestWorker_ExitCode(t *testing.T) {
	db := repository.SetupTestDB(t)
	_ = db

	t.Setenv("KSK8S_LOG_DIR", t.TempDir())

	logRepo := repository.NewLogRepository()
	rings := websocket.NewLogRingMap()

	logWriter, err := NewLogWriter(2, 1, logRepo, rings)
	if err != nil {
		t.Fatalf("failed to create log writer: %v", err)
	}
	defer logWriter.Close()

	// Use /bin/false which exits with code 1
	worker := NewWorkerWithCommand(2, 1, "/bin/false", []string{}, logWriter)

	if err := worker.Start(); err != nil {
		t.Fatalf("failed to start worker: %v", err)
	}

	worker.Wait()

	code := ParseExitCode(worker.Err())
	if code != 1 {
		t.Errorf("expected exit code 1, got %d", code)
	}
}

func TestWorker_Kill(t *testing.T) {
	db := repository.SetupTestDB(t)
	_ = db

	t.Setenv("KSK8S_LOG_DIR", t.TempDir())

	logRepo := repository.NewLogRepository()
	rings := websocket.NewLogRingMap()

	logWriter, err := NewLogWriter(3, 1, logRepo, rings)
	if err != nil {
		t.Fatalf("failed to create log writer: %v", err)
	}
	defer logWriter.Close()

	// Use sleep command that runs for a long time
	worker := NewWorkerWithCommand(3, 1, "/bin/sleep", []string{"30"}, logWriter)

	if err := worker.Start(); err != nil {
		t.Fatalf("failed to start worker: %v", err)
	}

	// Give it time to start
	time.Sleep(100 * time.Millisecond)

	if worker.PID() == 0 {
		t.Fatal("worker PID should be set")
	}

	// Kill the worker
	if err := worker.Kill(); err != nil {
		t.Fatalf("failed to kill worker: %v", err)
	}

	// Wait for it to die
	done := make(chan struct{})
	go func() {
		worker.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Expected
	case <-time.After(2 * time.Second):
		t.Fatal("worker did not die after kill")
	}

	// After kill, exit code should indicate signal
	code := ParseExitCode(worker.Err())
	if code == 0 {
		t.Error("expected non-zero exit code after kill")
	}
}

func TestParseExitCode_Nil(t *testing.T) {
	code := ParseExitCode(nil)
	if code != 0 {
		t.Errorf("expected 0 for nil error, got %d", code)
	}
}

func TestParseExitCode_ExitError(t *testing.T) {
	cmd := exec.Command("/bin/false")
	err := cmd.Run()
	code := ParseExitCode(err)
	if code != 1 {
		t.Errorf("expected 1 for /bin/false, got %d", code)
	}
}

func TestParseExitCode_UnknownError(t *testing.T) {
	code := ParseExitCode(os.ErrNotExist)
	if code != -1 {
		t.Errorf("expected -1 for unknown error, got %d", code)
	}
}

func TestReattach(t *testing.T) {
	// PID 0 should not be reattachable
	if Reattach(0) {
		t.Error("PID 0 should not be reattachable")
	}

	// PID -1 should not be reattachable
	if Reattach(-1) {
		t.Error("PID -1 should not be reattachable")
	}

	// PID 1 (init) should be reattachable on Linux
	// Skip this test if not on Linux or if PID 1 doesn't exist
	if Reattach(1) {
		// Expected on Linux
	}
}

func TestWorker_ProcessGroup(t *testing.T) {
	// Verify that Setpgid is set so Kill can target the process group
	worker := NewWorkerWithCommand(4, 1, "/bin/echo", []string{"test"}, nil)

	if worker.cmd.SysProcAttr == nil {
		t.Fatal("SysProcAttr should be set")
	}

	if !worker.cmd.SysProcAttr.Setpgid {
		t.Error("Setpgid should be true for process group isolation")
	}
}

func TestWorker_TaskAndClusterID(t *testing.T) {
	worker := NewWorkerWithCommand(42, 99, "/bin/echo", []string{"test"}, nil)

	if worker.TaskID() != 42 {
		t.Errorf("expected task ID 42, got %d", worker.TaskID())
	}

	if worker.ClusterID() != 99 {
		t.Errorf("expected cluster ID 99, got %d", worker.ClusterID())
	}
}

func TestLogWriter_WriteLine(t *testing.T) {
	db := repository.SetupTestDB(t)
	_ = db

	t.Setenv("KSK8S_LOG_DIR", t.TempDir())

	logRepo := repository.NewLogRepository()
	rings := websocket.NewLogRingMap()

	logWriter, err := NewLogWriter(5, 1, logRepo, rings)
	if err != nil {
		t.Fatalf("failed to create log writer: %v", err)
	}
	defer logWriter.Close()

	if err := logWriter.WriteLine("test stdout", "stdout"); err != nil {
		t.Errorf("WriteLine failed: %v", err)
	}

	if err := logWriter.WriteLine("test stderr", "stderr"); err != nil {
		t.Errorf("WriteLine failed: %v", err)
	}

	// Check ring buffer
	ring := rings.Get(5)
	if ring.Total() != 2 {
		t.Errorf("expected 2 lines in ring, got %d", ring.Total())
	}

	// Check file exists
	if _, err := os.Stat(logWriter.Path()); os.IsNotExist(err) {
		t.Error("log file should exist")
	}
}

func TestLogWriter_Flush(t *testing.T) {
	db := repository.SetupTestDB(t)
	_ = db

	t.Setenv("KSK8S_LOG_DIR", t.TempDir())

	logRepo := repository.NewLogRepository()
	rings := websocket.NewLogRingMap()

	logWriter, err := NewLogWriter(6, 1, logRepo, rings)
	if err != nil {
		t.Fatalf("failed to create log writer: %v", err)
	}
	defer logWriter.Close()

	// Write lines
	for i := 0; i < 5; i++ {
		logWriter.WriteLine("line", "stdout")
	}

	// Manual flush
	if err := logWriter.Flush(); err != nil {
		t.Errorf("Flush failed: %v", err)
	}
}

func TestLogWriter_MaxLineNumber(t *testing.T) {
	db := repository.SetupTestDB(t)
	_ = db

	t.Setenv("KSK8S_LOG_DIR", t.TempDir())

	logRepo := repository.NewLogRepository()
	rings := websocket.NewLogRingMap()

	// Pre-seed some logs
	logRepo.Create(&model.Log{TaskID: 7, LineNumber: 5, Content: "existing"})

	logWriter, err := NewLogWriter(7, 1, logRepo, rings)
	if err != nil {
		t.Fatalf("failed to create log writer: %v", err)
	}
	defer logWriter.Close()

	// First line should be 6 (5 + 1)
	logWriter.WriteLine("new line", "stdout")

	ring := rings.Get(7)
	lines, _ := ring.Snapshot()
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}
	if lines[0].LineNumber != 6 {
		t.Errorf("expected line number 6, got %d", lines[0].LineNumber)
	}
}
