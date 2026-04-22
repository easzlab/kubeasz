package service

import (
	"context"
	"fmt"
)

// Semaphore limits concurrent task execution.
type Semaphore struct {
	ch chan struct{}
}

// NewSemaphore creates a semaphore with the given capacity.
func NewSemaphore(capacity int) *Semaphore {
	return &Semaphore{
		ch: make(chan struct{}, capacity),
	}
}

// Acquire blocks until a slot is available.
func (s *Semaphore) Acquire(ctx context.Context) error {
	select {
	case s.ch <- struct{}{}:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("semaphore acquire cancelled: %w", ctx.Err())
	}
}

// Release returns a slot.
func (s *Semaphore) Release() {
	select {
	case <-s.ch:
	default:
	}
}

// Available returns the number of available slots.
func (s *Semaphore) Available() int {
	return cap(s.ch) - len(s.ch)
}
