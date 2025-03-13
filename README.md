# ChanWitch

ChanWitch 🧙‍♀️ is a lightweight Go library for managing named channels and providing automatic timeout-based channel cleanup. It offers two powerful components:

1. `ChanWitch`: A manager for named channels with automatic cleanup
2. `PoofChan`: A generic channel that automatically closes after a period of inactivity

## Features

### ChanWitch
- 🌟 Add named channels
- 🔍 Retrieve channels by name
- ❌ Close individual or all channels
- 📊 Get the number of stored channels
- 📝 String representation of stored channels
- 🔒 Thread-safe operations

### PoofChan
- ⏰ Automatic timeout-based closing
- 🔄 Activity-based timer reset
- 🎯 Generic type support
- 🔔 Optional close and reset callbacks
- 🔒 Thread-safe operations

## Installation

```sh
go get github.com/SamanFekri/chanwitch
```

## Usage

### ChanWitch Basic Usage

```go
package main

import (
    "fmt"
    "github.com/SamanFekri/chanwitch"
)

func main() {
    cw := chanwitch.NewChanWitch()
    ch := make(chan interface{})

    err := cw.Add("myChannel", ch)
    if err != nil {
        fmt.Println("Error:", err)
    }

    retrieved := cw.Get("myChannel")
    if retrieved != nil {
        fmt.Println("Channel found!")
    }

    fmt.Println("Number of channels:", cw.Len())

    cw.Close("myChannel")
}
```

### PoofChan Usage

```go
package main

import (
    "fmt"
    "time"
    "github.com/SamanFekri/chanwitch"
)

func main() {
    // Create a PoofChan with 1 buffer size and 5 second timeout
    pc := chanwitch.NewPoofChan[int](1, 5*time.Second,
        func() {
            fmt.Println("Channel closed due to inactivity")
        },
        func() {
            fmt.Println("Channel activity detected, timer reset")
        },
    )

    // Send a value
    pc.Send(42)

    // Receive a value
    value, ok := pc.Receive()
    if ok {
        fmt.Println("Received:", value)
    }
}
```

### Using ChanWitch with PoofChan

```go
package main

import (
    "time"
    "github.com/SamanFekri/chanwitch"
)

func main() {
    cw := chanwitch.NewChanWitch()

    // Create a PoofChan for temporary data processing
    pc := chanwitch.NewPoofChan[int](10, 30*time.Second,
        func() {
            // Clean up when the channel closes
            cw.Remove("tempProcessing")
        },
        nil,
    )

    // Add the PoofChan to ChanWitch
    err := cw.Add("tempProcessing", pc)
    if err != nil {
        // Handle error
    }

    // Use the channel
    go func() {
        if ch, ok := cw.Get("tempProcessing").(*chanwitch.PoofChan[int]); ok {
            ch.Send(42)
        }
    }()
}
```

## API

### ChanWitch API

#### `func NewChanWitch() *ChanWitch`
Creates a new instance of ChanWitch.

#### `func (cw *ChanWitch) Add(name string, ch interface{}) error`
Adds a new channel with a name. Returns an error if the name already exists.

#### `func (cw *ChanWitch) Get(name string) interface{}`
Retrieves a channel by its name. Returns `nil` if the channel does not exist.

#### `func (cw *ChanWitch) Close(name string) error`
Closes the specified channel and removes it from the list. Returns an error if the channel doesn't exist.

#### `func (cw *ChanWitch) CloseAll()`
Closes all stored channels and clears the list.

#### `func (cw *ChanWitch) Remove(name string)`
Removes a channel from the list without closing it.

#### `func (cw *ChanWitch) String() string`
Returns a string representation of all stored channels.

#### `func (cw *ChanWitch) Len() int`
Returns the number of channels in the list.

### PoofChan API

#### `func NewPoofChan[T any](size int, duration time.Duration, closeCb func(), resetCb func()) *PoofChan[T]`
Creates a new PoofChan with the specified buffer size, timeout duration, and optional callbacks.

#### `func (pc *PoofChan[T]) Send(value T) bool`
Sends a value to the channel. Returns false if the channel is closed.

#### `func (pc *PoofChan[T]) Receive() (T, bool)`
Receives a value from the channel. Returns false if the channel is closed.

## Use Cases

1. **Temporary Data Processing** 🚀
   - Use PoofChan for processing temporary data streams
   - Automatically clean up when processing is complete

2. **WebSocket Connections** 🌐
   - Manage multiple WebSocket connections with ChanWitch
   - Use PoofChan for individual connection timeouts

3. **Background Tasks** ⚙️
   - Track long-running tasks with ChanWitch
   - Use PoofChan for task-specific timeouts

4. **Resource Management** 🔧
   - Manage multiple resources with different lifetimes
   - Automatic cleanup of inactive resources

## License
This project is licensed under the MIT License. Feel free to use, modify, and distribute it.

## Contributions
Pull requests and issues are welcome! 🚀

## Author
Created by **Saman**.

