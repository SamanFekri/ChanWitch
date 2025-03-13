package chanwitch

import (
	"testing"
	"time"
)

func TestNewPoofChan(t *testing.T) {
	// Test with minimal configuration
	pc := NewPoofChan[int](1, 100*time.Millisecond, nil, nil)
	if pc == nil {
		t.Error("NewPoofChan returned nil")
	}

	// Test with callbacks
	pc = NewPoofChan[int](1, 100*time.Millisecond,
		func() { t.Log("Channel closed") },
		func() { t.Log("Channel reset") },
	)
	if pc == nil {
		t.Error("NewPoofChan returned nil with callbacks")
	}
}

func TestPoofChanSendReceive(t *testing.T) {
	pc := NewPoofChan[int](1, 100*time.Millisecond, nil, nil)

	// Test sending and receiving
	value := 42
	if !pc.Send(value) {
		t.Error("Send failed")
	}

	received, ok := pc.Receive()
	if !ok {
		t.Error("Receive failed")
	}
	if received != value {
		t.Errorf("Received %v, expected %v", received, value)
	}
}

func TestPoofChanTimeout(t *testing.T) {
	closeCalled := false
	pc := NewPoofChan[int](1, 50*time.Millisecond, func() { closeCalled = true }, nil)

	// Send a value
	pc.Send(42)
	pc.Receive()

	// Wait for timeout
	time.Sleep(100 * time.Millisecond)

	if !closeCalled {
		t.Error("Close callback should have been called")
	}
}

func TestPoofChanResetOnActivity(t *testing.T) {
	resetCalled := false
	pc := NewPoofChan[int](2, 100*time.Millisecond, nil, func() { resetCalled = true })

	// Send initial value
	pc.Send(42)

	// Wait for half the timeout duration
	time.Sleep(50 * time.Millisecond)

	// Send another value to reset timer
	pc.Send(43)

	// Wait for half the timeout duration again
	time.Sleep(150 * time.Millisecond)

	// Channel should still be open
	_, ok := pc.Receive()
	if !ok {
		t.Error("Channel should still be open after reset")
	}

	if !resetCalled {
		t.Error("Reset callback should have been called")
	}
}

func TestPoofChanStringType(t *testing.T) {
	pc := NewPoofChan[string](1, 100*time.Millisecond, nil, nil)

	// Test sending and receiving string
	value := "test"
	if !pc.Send(value) {
		t.Error("Send failed")
	}

	received, ok := pc.Receive()
	if !ok {
		t.Error("Receive failed")
	}
	if received != value {
		t.Errorf("Received %v, expected %v", received, value)
	}
} 