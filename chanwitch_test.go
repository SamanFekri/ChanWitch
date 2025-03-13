package chanwitch

import (
	"testing"
	"time"
)

func TestNewChanWitch(t *testing.T) {
	cw := NewChanWitch()
	if cw == nil {
		t.Error("NewChanWitch returned nil")
	}
}

func TestAdd(t *testing.T) {
	cw := NewChanWitch()
	ch := make(chan interface{})

	// Test adding a new channel
	err := cw.Add("test1", ch)
	if err != nil {
		t.Errorf("Add failed: %v", err)
	}

	// Test adding a channel with existing name
	err = cw.Add("test1", ch)
	if err == nil {
		t.Error("Add should return error for existing channel name")
	}
}

func TestGet(t *testing.T) {
	cw := NewChanWitch()
	ch := make(chan interface{})
	cw.Add("test1", ch)

	// Test getting existing channel
	got := cw.Get("test1")
	if got != ch {
		t.Error("Get returned wrong channel")
	}

	// Test getting non-existing channel
	got = cw.Get("nonexistent")
	if got != nil {
		t.Error("Get should return nil for non-existing channel")
	}
}

func TestClose(t *testing.T) {
	cw := NewChanWitch()
	ch := make(chan interface{})
	cw.Add("test1", ch)

	// Test closing existing channel
	cw.Close("test1")

	// Verify channel is removed
	if cw.Get("test1") != nil {
		t.Error("Channel should be removed after closing")
	}

	// Test closing non-existing channel (should not panic)
	cw.Close("nonexistent")
}

func TestCloseAll(t *testing.T) {
	cw := NewChanWitch()
	ch1 := make(chan interface{})
	ch2 := make(chan interface{})
	cw.Add("test1", ch1)
	cw.Add("test2", ch2)

	cw.CloseAll()

	// Verify all channels are removed
	if cw.Get("test1") != nil || cw.Get("test2") != nil {
		t.Error("All channels should be removed after CloseAll")
	}
}

func TestString(t *testing.T) {
	cw := NewChanWitch()
	ch1 := make(chan interface{})
	ch2 := make(chan interface{})
	cw.Add("test1", ch1)
	cw.Add("test2", ch2)

	time.Sleep(300 * time.Millisecond)

	expected := " <- test1\n <- test2"
	got := cw.String()
	if got != expected {
		t.Errorf("String() = %v, want %v", got, expected)
	}

	// Test empty ChanWitch
	cw = NewChanWitch()
	got = cw.String()
	if got != "" {
		t.Error("String() should return empty string for empty ChanWitch")
	}
}

func TestLen(t *testing.T) {
	cw := NewChanWitch()
	ch1 := make(chan interface{})
	ch2 := make(chan interface{})
	cw.Add("test1", ch1)
	cw.Add("test2", ch2)

	if cw.Len() != 2 {
		t.Error("Len() should return 2")
	}

	// Test after Close a channel
	cw.Close("test1")
	if cw.Len() != 1 {
		t.Error("Len() should return 1 after Close a channel")
	}

	// Add a new channel
	ch3 := make(chan interface{})
	cw.Add("test3", ch3)
	if cw.Len() != 2 {
		t.Error("Len() should return 2 after Add a new channel")
	}

	// Close all channels
	cw.CloseAll()
	if cw.Len() != 0 {
		t.Error("Len() should return 0 after CloseAll")
	}

	// Test empty ChanWitch
	cw = NewChanWitch()
	if cw.Len() != 0 {
		t.Error("Len() should return 0 for empty ChanWitch")
	}

}

// Test ChanWitch with PoofChan
func TestPoofChan(t *testing.T) {
	cw := NewChanWitch()
	ch1 := NewPoofChan[int](2, 1*time.Second, func() { cw.Close("test1") }, nil)

	cw.Add("test1", ch1)

	go func() {
		for i := range 5 {
			ch1.Send(i)
		}
	}()

	// listen to the channel
	for i := range 5 {
		v, ok := ch1.Receive()
		if !ok {
			t.Errorf("Recvive() failed")
		}
		if v != i {
			t.Errorf("Recvive() = %v, want %v", v, i)
		}
	}

	time.Sleep(2 * time.Second)
	// Verify channel is removed
	if cw.Get("test1") != nil {
		t.Error("Channel should be removed after closing")
	}
}
