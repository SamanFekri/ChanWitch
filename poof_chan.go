package chanwitch

import (
	"time"
)

// PoofChan provides a generic channel that automatically closes if it remains inactive for a specified duration.
type PoofChan[T any] struct {
	ch        chan T        // The main channel for data communication
	timer     *time.Timer   // Timer used to track inactivity duration
	duration  time.Duration // The configured timeout duration
	closeCb   func()        // Callback function invoked when the channel closes
	resetCb   func()        // Callback function invoked when the timer resets due to non-empty channel
	resetChan chan struct{} // Internal signal channel to reset the timer
}

// NewPoofChan initializes a PoofChan with a given timeout duration and an optional close callback.
func NewPoofChan[T any](size int, duration time.Duration, closeCb func(), resetCb func()) *PoofChan[T] {
	pc := &PoofChan[T]{
		ch:        make(chan T, size),
		timer:     time.NewTimer(duration),
		duration:  duration,
		closeCb:   closeCb,
		resetCb:   resetCb,
		resetChan: make(chan struct{}, 1),
	}

	// Goroutine that monitors inactivity and closes the channel when the timer expires.
	go func() {
		for {
			select {
			case <-pc.timer.C: // Timer expired, check if the channel is empty
				if len(pc.ch) > 0 { // The channel is not empty, reset the timer
					pc.timer.Reset(pc.duration)
					if pc.resetCb != nil {
						pc.resetCb()
					}
					continue
				}
				pc.Close()
				return
			case <-pc.resetChan: // Reset the timer on activity
				if !pc.timer.Stop() {
					<-pc.timer.C // Drain the timer if it has already expired
				}
				pc.timer.Reset(pc.duration)
			}
		}
	}()

	return pc
}

// Send transmits a value to the channel, resetting the inactivity timer upon success.
// Returns false if the channel is already closed.
func (pc *PoofChan[T]) Send(value T) bool {
	select {
	case <-pc.timer.C:
		return false // Channel is already closed
	default:
		// Signal timer reset if possible
		select {
		case pc.resetChan <- struct{}{}:
		default:
		}
		pc.ch <- value
		return true
	}
}

// Receive retrieves a value from the channel, resetting the inactivity timer if successful.
// Returns false if the channel is closed.
func (pc *PoofChan[T]) Receive() (T, bool) {
	val, ok := <-pc.ch
	if ok {
		// Reset the timer on a successful receive
		select {
		case pc.resetChan <- struct{}{}:
		default:
		}
	}
	return val, ok
}

// Close closes the channel and stops the timer.
func (pc *PoofChan[T]) Close() {
	pc.safeClose(pc.ch)
	pc.timer.Stop()
}

// safeClose closes a channel and prevents a panic if it's already closed.
func (pc *PoofChan[T]) safeClose(ch chan T) {
	defer func() {
		if r := recover(); r != nil {
			// Recover from panic if the channel is already closed
		}
	}()
	// close the ch
	close(ch)
	if pc.closeCb != nil {
		pc.closeCb()
	}
}
