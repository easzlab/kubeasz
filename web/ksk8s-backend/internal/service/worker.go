package service

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"syscall"
)

// Worker represents a running subprocess that executes ezctl commands.
type Worker struct {
	cmd       *exec.Cmd
	taskID    int64
	clusterID int64
	logWriter *LogWriter
	pid       int
	done      chan struct{}
	err       error
}

// ezctlPath returns the path to the ezctl executable, defaulting to /etc/kubeasz/ezctl.
func ezctlPath() string {
	if p := os.Getenv("KSK8S_EZCTL_PATH"); p != "" {
		return p
	}
	return "/etc/kubeasz/ezctl"
}

// NewWorker creates a Worker that will run the given command.
func NewWorker(taskID int64, clusterID int64, clusterName string, step string, logWriter *LogWriter, extraArgs ...string) *Worker {
	// Build ezctl command
	var cmd *exec.Cmd
	path := ezctlPath()
	args := []string{}
	if step == "90" || step == "full" {
		args = []string{"setup", clusterName, "90"}
	} else if step == "start" {
		args = []string{"start", clusterName}
	} else if step == "stop" {
		args = []string{"stop", clusterName}
	} else if step == "upgrade" {
		args = []string{"upgrade", clusterName}
	} else if step == "backup" {
		args = []string{"backup", clusterName}
	} else if step == "restore" {
		args = []string{"restore", clusterName}
	} else if step == "destroy" {
		args = []string{"destroy", clusterName}
	} else if step == "kca-renew" {
		args = []string{"kca-renew", clusterName}
	} else if step == "add-node" {
		args = []string{"add-node", clusterName}
	} else if step == "del-node" {
		args = []string{"del-node", clusterName}
	} else if step == "add-master" {
		args = []string{"add-master", clusterName}
	} else if step == "del-master" {
		args = []string{"del-master", clusterName}
	} else if step == "add-etcd" {
		args = []string{"add-etcd", clusterName}
	} else if step == "del-etcd" {
		args = []string{"del-etcd", clusterName}
	} else {
		args = []string{"setup", clusterName, step}
	}
	args = append(args, extraArgs...)
	cmd = exec.Command(path, args...)

	env := os.Environ()
	// Inject locale so Ansible modules on target nodes use UTF-8 encoding
	env = append(env, "LC_ALL=C.UTF-8", "LANG=C.UTF-8")
	cmd.Env = env
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	return &Worker{
		cmd:       cmd,
		taskID:    taskID,
		clusterID: clusterID,
		logWriter: logWriter,
		done:      make(chan struct{}),
	}
}

// NewWorkerWithCommand creates a Worker with an arbitrary command (for testing/mock).
func NewWorkerWithCommand(taskID int64, clusterID int64, name string, args []string, logWriter *LogWriter) *Worker {
	cmd := exec.Command(name, args...)
	cmd.Env = os.Environ()
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	return &Worker{
		cmd:       cmd,
		taskID:    taskID,
		clusterID: clusterID,
		logWriter: logWriter,
		done:      make(chan struct{}),
	}
	}

// Start begins execution and captures stdout/stderr.
func (w *Worker) Start() error {
	stdout, err := w.cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("stdout pipe: %w", err)
	}
	stderr, err := w.cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("stderr pipe: %w", err)
	}

	if err := w.cmd.Start(); err != nil {
		return fmt.Errorf("start worker: %w", err)
	}

	w.pid = w.cmd.Process.Pid

	go w.scan(stdout, "stdout")
	go w.scan(stderr, "stderr")

	go func() {
		w.err = w.cmd.Wait()
		_ = w.logWriter.Flush()
		close(w.done)
	}()

	return nil
}

func (w *Worker) scan(pipe io.Reader, stream string) {
	scanner := bufio.NewScanner(pipe)
	scanner.Buffer(make([]byte, 4096), 1024*1024)
	for scanner.Scan() {
		_ = w.logWriter.WriteLine(scanner.Text(), stream)
	}
}

// Done returns a channel that closes when the worker exits.
func (w *Worker) Done() <-chan struct{} {
	return w.done
}

// Err returns the exit error, if any. Only valid after Done is closed.
func (w *Worker) Err() error {
	return w.err
}

// PID returns the process ID.
func (w *Worker) PID() int {
	return w.pid
}

// Kill sends SIGKILL to the entire process group.
func (w *Worker) Kill() error {
	if w.cmd.Process == nil {
		return nil
	}
	return syscall.Kill(-w.cmd.Process.Pid, syscall.SIGKILL)
}

// TaskID returns the associated task ID.
func (w *Worker) TaskID() int64 {
	return w.taskID
}

// ClusterID returns the associated cluster ID.
func (w *Worker) ClusterID() int64 {
	return w.clusterID
}

// Wait blocks until the worker exits.
func (w *Worker) Wait() error {
	<-w.done
	return w.err
}

// Reattach tries to re-attach to an existing process by PID.
// Returns true if the process exists and we successfully signaled it.
func Reattach(pid int) bool {
	if pid <= 0 {
		return false
	}
	err := syscall.Kill(pid, 0)
	return err == nil
}

// ParseExitCode extracts the exit code from an *exec.ExitError.
func ParseExitCode(err error) int {
	if err == nil {
		return 0
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		if ws, ok := exitErr.Sys().(syscall.WaitStatus); ok {
			return ws.ExitStatus()
		}
	}
	return -1
}

// ExitCodeStr returns the exit code as a string for DB storage.
func ExitCodeStr(err error) *string {
	code := strconv.Itoa(ParseExitCode(err))
	return &code
}
