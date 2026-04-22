package websocket

import (
	"sync"
)

// LogRing is a thread-safe circular buffer for log lines.
type LogRing struct {
	mu       sync.RWMutex
	lines    []LogLine
	capacity int
	head     int
	size     int
	total    int
}

type LogLine struct {
	LineNumber int    `json:"line_number"`
	Content    string `json:"content"`
	Stream     string `json:"stream"`
	Timestamp  int64  `json:"timestamp"`
}

// NewLogRing creates a ring buffer with the given capacity.
func NewLogRing(capacity int) *LogRing {
	return &LogRing{
		lines:    make([]LogLine, capacity),
		capacity: capacity,
	}
}

// Append adds a line to the ring. Overwrites oldest if full.
func (r *LogRing) Append(line LogLine) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.lines[r.head] = line
	r.head = (r.head + 1) % r.capacity
	if r.size < r.capacity {
		r.size++
	}
	r.total++
}

// Total returns the total number of lines ever appended.
func (r *LogRing) Total() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.total
}

// Snapshot returns a copy of all current lines in order, plus the total count.
func (r *LogRing) Snapshot() ([]LogLine, int) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]LogLine, r.size)
	if r.size < r.capacity {
		copy(result, r.lines[:r.size])
	} else {
		copy(result, r.lines[r.head:])
		copy(result[r.capacity-r.head:], r.lines[:r.head])
	}
	return result, r.total
}

// Since returns lines with line_number > offset.
func (r *LogRing) Since(offset int) ([]LogLine, int) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.total == 0 || offset >= r.total {
		return nil, r.total
	}

	// Determine how many lines to return
	start := offset
	if start < r.total-r.size {
		start = r.total - r.size
	}
	count := r.total - start
	result := make([]LogLine, count)

	for i := 0; i < count; i++ {
		idx := (r.head - r.size + start - r.total + r.size + i + r.capacity) % r.capacity
		if r.size < r.capacity {
			idx = start + i
		}
		result[i] = r.lines[idx]
	}
	return result, r.total
}

// LogRingMap holds a ring buffer per task ID.
type LogRingMap struct {
	mu    sync.RWMutex
	rings map[int64]*LogRing
}

func NewLogRingMap() *LogRingMap {
	return &LogRingMap{
		rings: make(map[int64]*LogRing),
	}
}

func (m *LogRingMap) Get(taskID int64) *LogRing {
	m.mu.RLock()
	r, ok := m.rings[taskID]
	m.mu.RUnlock()
	if ok {
		return r
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	r, ok = m.rings[taskID]
	if !ok {
		r = NewLogRing(10000)
		m.rings[taskID] = r
	}
	return r
}

func (m *LogRingMap) Delete(taskID int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.rings, taskID)
}

// Iterate calls fn for each task ID and its ring. Safe for read-only access.
func (m *LogRingMap) Iterate(fn func(taskID int64, ring *LogRing)) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for id, r := range m.rings {
		fn(id, r)
	}
}
