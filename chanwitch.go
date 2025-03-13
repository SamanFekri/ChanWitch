package chanwitch

import (
    "sync"
    "errors"
)

// ChanWitch manages named channels with an inactivity timeout.
type ChanWitch struct {
    channels map[string]interface{}
    mutex    sync.Mutex
}

// NewChanWitch creates a new ChanWitch instance.
func NewChanWitch() *ChanWitch {
    return &ChanWitch{
        channels: make(map[string]interface{}),
    }
}

// Add adds an existing channel to ChanWitch.
// if the name exists, return an error.
func (cw *ChanWitch) Add(name string, ch interface{}) error {
    cw.mutex.Lock()
    defer cw.mutex.Unlock()

		if _, exists := cw.channels[name]; exists {
			return errors.New("Channel already exists")
		}

    cw.channels[name] = ch
		return nil
}

// Remove removes a channel from the list.
func (cw *ChanWitch) Remove(name string) {
    cw.mutex.Lock()
    defer cw.mutex.Unlock()
    // check the channel exists
    if _, exists := cw.channels[name]; !exists {
        return
    }
    // Remove the channel from the list
    delete(cw.channels, name)
}

// Get returns an existing channel if it exists, otherwise nil.
func (cw *ChanWitch) Get(name string) interface{} {
    cw.mutex.Lock()
    defer cw.mutex.Unlock()

    if ch, exists := cw.channels[name]; exists {
        return ch
    }

    return nil
}

// Close closes all the channels in the list.
func (cw *ChanWitch) CloseAll() {
    cw.mutex.Lock()
    defer cw.mutex.Unlock()

    for name, ch := range cw.channels {
        safeClose(ch.(chan interface{}))
        // Remove the channel from the list
        delete(cw.channels, name)
    }
}

// Close closes the specified channel. If the channel does not exist, return an error.
func (cw *ChanWitch) Close(name string) error {
    cw.mutex.Lock()
    defer cw.mutex.Unlock()

    if ch, exists := cw.channels[name]; exists {
        safeClose(ch.(chan interface{}))
        // Remove the channel from the list
        delete(cw.channels, name)
        return nil
    }
    return errors.New("Channel does not exist")
}

// String returns a string representation of the ChanWitch.
func (cw *ChanWitch) String() string {
    str := ""
    for name := range cw.channels {
        str += " <- " + name + "\n"
    }
    // Remove the last newline
    if len(str) > 0 {
        str = str[:len(str)-1]
    }
    return str
}

// Len returns the number of channels in the list.
func (cw *ChanWitch) Len() int {
    return len(cw.channels)
}

// safeClose closes a channel and prevents a panic if it's already closed.
func safeClose(ch chan interface{}) {
    defer func() {
        if r := recover(); r != nil {
            // Recover from panic if the channel is already closed
        }
    }()
    close(ch)
}
