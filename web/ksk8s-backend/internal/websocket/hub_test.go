package websocket

import (
	"testing"
	"time"
)

func TestHub_RegisterAndUnregister(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	client := &Client{
		Hub:    hub,
		Send:   make(chan []byte, 256),
		TaskID: 1,
	}

	hub.Register(client)
	time.Sleep(50 * time.Millisecond)

	if hub.ClientCount(1) != 1 {
		t.Errorf("expected 1 client for task 1, got %d", hub.ClientCount(1))
	}

	hub.Unregister(client)
	time.Sleep(50 * time.Millisecond)

	if hub.ClientCount(1) != 0 {
		t.Errorf("expected 0 clients for task 1 after unregister, got %d", hub.ClientCount(1))
	}
}

func TestHub_Broadcast(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	client := &Client{
		Hub:    hub,
		Send:   make(chan []byte, 256),
		TaskID: 1,
	}

	hub.Register(client)
	time.Sleep(50 * time.Millisecond)

	msg := []byte("test message")
	hub.Broadcast(1, msg)

	select {
	case received := <-client.Send:
		if string(received) != string(msg) {
			t.Errorf("expected %s, got %s", msg, received)
		}
	case <-time.After(500 * time.Millisecond):
		t.Error("timeout waiting for broadcast message")
	}

	hub.Unregister(client)
}

func TestHub_Broadcast_MultipleClients(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	client1 := &Client{
		Hub:    hub,
		Send:   make(chan []byte, 256),
		TaskID: 1,
	}
	client2 := &Client{
		Hub:    hub,
		Send:   make(chan []byte, 256),
		TaskID: 1,
	}
	client3 := &Client{
		Hub:    hub,
		Send:   make(chan []byte, 256),
		TaskID: 2, // different task
	}

	hub.Register(client1)
	hub.Register(client2)
	hub.Register(client3)
	time.Sleep(50 * time.Millisecond)

	msg := []byte("task 1 message")
	hub.Broadcast(1, msg)

	// Both clients for task 1 should receive
	for i, client := range []*Client{client1, client2} {
		select {
		case received := <-client.Send:
			if string(received) != string(msg) {
				t.Errorf("client %d: expected %s, got %s", i, msg, received)
			}
		case <-time.After(500 * time.Millisecond):
			t.Errorf("client %d: timeout waiting for broadcast", i)
		}
	}

	// Client 3 should not receive
	select {
	case <-client3.Send:
		t.Error("client 3 should not receive message for task 1")
	case <-time.After(100 * time.Millisecond):
		// Expected - no message
	}

	hub.Unregister(client1)
	hub.Unregister(client2)
	hub.Unregister(client3)
}

func TestHub_Broadcast_SlowClientDrop(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	// Client with tiny buffer
	client := &Client{
		Hub:    hub,
		Send:   make(chan []byte, 1),
		TaskID: 1,
	}

	hub.Register(client)
	time.Sleep(50 * time.Millisecond)

	// Fill the buffer
	hub.Broadcast(1, []byte("msg1"))
	time.Sleep(10 * time.Millisecond)
	// Don't read from client.Send - simulate slow consumer

	// This broadcast should be dropped, not block
	hub.Broadcast(1, []byte("msg2"))

	hub.Unregister(client)
}

func TestHub_ClientCount_EmptyTask(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	if hub.ClientCount(999) != 0 {
		t.Error("expected 0 clients for non-existent task")
	}
}

func TestLogRing_Basic(t *testing.T) {
	ring := NewLogRing(5)

	ring.Append(LogLine{LineNumber: 1, Content: "line 1"})
	ring.Append(LogLine{LineNumber: 2, Content: "line 2"})

	if ring.Total() != 2 {
		t.Errorf("expected total 2, got %d", ring.Total())
	}

	lines, total := ring.Snapshot()
	if total != 2 {
		t.Errorf("expected snapshot total 2, got %d", total)
	}
	if len(lines) != 2 {
		t.Errorf("expected 2 lines, got %d", len(lines))
	}
	if lines[0].Content != "line 1" {
		t.Errorf("expected line 1, got %s", lines[0].Content)
	}
}

func TestLogRing_Overflow(t *testing.T) {
	ring := NewLogRing(3)

	for i := 1; i <= 5; i++ {
		ring.Append(LogLine{LineNumber: i, Content: "line"})
	}

	if ring.Total() != 5 {
		t.Errorf("expected total 5, got %d", ring.Total())
	}

	lines, _ := ring.Snapshot()
	if len(lines) != 3 {
		t.Errorf("expected 3 lines (capacity), got %d", len(lines))
	}

	// Should contain lines 3, 4, 5
	if lines[0].LineNumber != 3 {
		t.Errorf("expected first line number 3, got %d", lines[0].LineNumber)
	}
}

func TestLogRing_Since(t *testing.T) {
	ring := NewLogRing(10)

	for i := 1; i <= 5; i++ {
		ring.Append(LogLine{LineNumber: i, Content: "line"})
	}

	lines, total := ring.Since(2)
	if total != 5 {
		t.Errorf("expected total 5, got %d", total)
	}
	if len(lines) != 3 {
		t.Errorf("expected 3 lines since 2, got %d", len(lines))
	}
}

func TestLogRing_Since_BeyondCapacity(t *testing.T) {
	ring := NewLogRing(3)

	for i := 1; i <= 10; i++ {
		ring.Append(LogLine{LineNumber: i, Content: "line"})
	}

	// Request lines since 5, but only 8, 9, 10 are in the ring
	lines, total := ring.Since(5)
	if total != 10 {
		t.Errorf("expected total 10, got %d", total)
	}
	if len(lines) != 3 {
		t.Errorf("expected 3 lines (capacity), got %d", len(lines))
	}
}

func TestLogRingMap_Get(t *testing.T) {
	m := NewLogRingMap()

	ring1 := m.Get(1)
	ring2 := m.Get(1)

	if ring1 != ring2 {
		t.Error("same task ID should return same ring")
	}

	ring3 := m.Get(2)
	if ring1 == ring3 {
		t.Error("different task IDs should return different rings")
	}
}

func TestLogRingMap_Delete(t *testing.T) {
	m := NewLogRingMap()

	ring := m.Get(1)
	ring.Append(LogLine{LineNumber: 1, Content: "test"})

	m.Delete(1)

	// After delete, Get should create a new ring
	ring2 := m.Get(1)
	if ring2.Total() != 0 {
		t.Error("new ring after delete should be empty")
	}
}

func TestLogRingMap_Iterate(t *testing.T) {
	m := NewLogRingMap()

	m.Get(1).Append(LogLine{LineNumber: 1})
	m.Get(2).Append(LogLine{LineNumber: 1})

	count := 0
	m.Iterate(func(taskID int64, ring *LogRing) {
		count++
	})

	if count != 2 {
		t.Errorf("expected 2 rings, got %d", count)
	}
}
