package repository

import (
	"sync"
	"time"

	"github.com/easzlab/ksk8s/internal/model"
)

type LogRepository struct {
	mu      sync.Mutex
	pending []model.Log
}

func NewLogRepository() *LogRepository {
	return &LogRepository{}
}

func (r *LogRepository) Create(log *model.Log) error {
	return DB.Create(log).Error
}

func (r *LogRepository) ListByTask(taskID int64, offset int, limit int) ([]model.Log, error) {
	var logs []model.Log
	err := DB.Where("task_id = ? AND line_number > ?", taskID, offset).
		Order("line_number asc").
		Limit(limit).
		Find(&logs).Error
	return logs, err
}

func (r *LogRepository) GetMaxLineNumber(taskID int64) (int, error) {
	var result struct {
		MaxLineNumber int
	}
	err := DB.Raw("SELECT COALESCE(MAX(line_number), 0) as max_line_number FROM logs WHERE task_id = ?", taskID).Scan(&result).Error
	return result.MaxLineNumber, err
}

// BatchInsert inserts multiple logs in a single transaction.
func (r *LogRepository) BatchInsert(logs []model.Log) error {
	if len(logs) == 0 {
		return nil
	}
	return DB.CreateInBatches(logs, 100).Error
}

// Enqueue adds a log to the pending buffer and flushes if threshold reached.
func (r *LogRepository) Enqueue(taskID int64, lineNumber int, content string, stream string) error {
	r.mu.Lock()
	r.pending = append(r.pending, model.Log{
		TaskID:     taskID,
		LineNumber: lineNumber,
		Content:    content,
		Stream:     stream,
		Timestamp:  time.Now(),
	})
	shouldFlush := len(r.pending) >= 100
	r.mu.Unlock()

	if shouldFlush {
		return r.Flush()
	}
	return nil
}

// Flush inserts all pending logs.
func (r *LogRepository) Flush() error {
	r.mu.Lock()
	if len(r.pending) == 0 {
		r.mu.Unlock()
		return nil
	}
	batch := make([]model.Log, len(r.pending))
	copy(batch, r.pending)
	r.pending = r.pending[:0]
	r.mu.Unlock()

	return r.BatchInsert(batch)
}
