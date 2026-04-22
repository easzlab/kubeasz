package service

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/easzlab/ksk8s/internal/repository"
	"github.com/easzlab/ksk8s/internal/websocket"
)

// LogWriter writes log lines to both a file and an in-memory ring buffer,
// with async batch inserts to MySQL.
type LogWriter struct {
	repo           *repository.LogRepository
	rings          *websocket.LogRingMap
	filePath       string
	file           *os.File
	writer         *bufio.Writer
	mu             sync.Mutex
	ticker         *time.Ticker
	stop           chan struct{}
	lineNum        int
	taskID         int64
	ansibleFailed  int
	ansibleHosts   int
	inPlayRecap    bool
}

var playRecapFailedRe = regexp.MustCompile(`failed=(\d+)`)
var playRecapLineRe = regexp.MustCompile(`^\s*\S+\s*:\s*ok=\d+`)

// NewLogWriter creates a LogWriter for a task.
func NewLogWriter(taskID int64, clusterID int64, repo *repository.LogRepository, rings *websocket.LogRingMap) (*LogWriter, error) {
	baseDir := os.Getenv("KSK8S_LOG_DIR")
	if baseDir == "" {
		baseDir = "/var/log/ksk8s"
	}
	logDir := filepath.Join(baseDir, fmt.Sprintf("%d", clusterID))
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("create log dir: %w", err)
	}
	logPath := filepath.Join(logDir, fmt.Sprintf("task_%d.log", taskID))
	f, err := os.Create(logPath)
	if err != nil {
		return nil, fmt.Errorf("create log file: %w", err)
	}

	lw := &LogWriter{
		repo:     repo,
		rings:    rings,
		filePath: logPath,
		file:     f,
		writer:   bufio.NewWriter(f),
		ticker:   time.NewTicker(1 * time.Second),
		stop:     make(chan struct{}),
		taskID:   taskID,
	}

	// Resume line number from DB if this is a re-attachment
	maxLine, _ := repo.GetMaxLineNumber(taskID)
	lw.lineNum = maxLine

	go lw.flushLoop()
	return lw, nil
}

// WriteLine writes a single line to file, ring buffer, and MySQL batch.
func (lw *LogWriter) WriteLine(content string, stream string) error {
	lw.mu.Lock()
	lw.lineNum++
	lineNum := lw.lineNum

	// Parse Ansible PLAY RECAP lines
	if strings.Contains(content, "PLAY RECAP") {
		lw.inPlayRecap = true
	} else if lw.inPlayRecap {
		if playRecapLineRe.MatchString(content) {
			lw.ansibleHosts++
			if m := playRecapFailedRe.FindStringSubmatch(content); len(m) > 1 {
				if n, _ := strconv.Atoi(m[1]); n > 0 {
					lw.ansibleFailed += n
				}
			}
		} else if strings.TrimSpace(content) == "" {
			lw.inPlayRecap = false
		}
	}
	lw.mu.Unlock()

	ts := time.Now().UnixMilli()

	// File write
	lw.mu.Lock()
	fmt.Fprintf(lw.writer, "[%s] %s\n", stream, content)
	lw.mu.Unlock()

	// Ring buffer
	line := websocket.LogLine{
		LineNumber: lineNum,
		Content:    content,
		Stream:     stream,
		Timestamp:  ts,
	}
	lw.rings.Get(lw.taskID).Append(line)

	// MySQL batch
	return lw.repo.Enqueue(lw.taskID, lineNum, content, stream)
}

// AnsibleResult returns the number of hosts and total failures seen in PLAY RECAP.
func (lw *LogWriter) AnsibleResult() (hosts int, failed int) {
	lw.mu.Lock()
	defer lw.mu.Unlock()
	return lw.ansibleHosts, lw.ansibleFailed
}

// flushLoop periodically flushes file and DB batch.
func (lw *LogWriter) flushLoop() {
	for {
		select {
		case <-lw.ticker.C:
			lw.Flush()
		case <-lw.stop:
			return
		}
	}
}

// Flush forces file and DB flush.
func (lw *LogWriter) Flush() error {
	lw.mu.Lock()
	err1 := lw.writer.Flush()
	err2 := lw.file.Sync()
	lw.mu.Unlock()

	err3 := lw.repo.Flush()

	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}
	return err3
}

// Close flushes everything and closes the file.
func (lw *LogWriter) Close() error {
	close(lw.stop)
	lw.ticker.Stop()
	_ = lw.Flush()
	return lw.file.Close()
}

// Path returns the log file path.
func (lw *LogWriter) Path() string {
	return lw.filePath
}
