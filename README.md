# ChanWitch

ChanWitch 🧙‍♀️ is a lightweight Go library for managing named channels. It allows you to store, retrieve, and close channels by their assigned names, making channel management easier and more organized.

## Features
- 🌟 Add named channels
- 🔍 Retrieve channels by name
- ❌ Close individual or all channels
- 📊 Get the number of stored channels
- 📝 String representation of stored channels

## Installation

```sh
go get github.com/SamanFekri/chanwitch
```

## Usage

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

## API

### `func NewChanWitch() *ChanWitch`
Creates a new instance of ChanWitch.

### `func (cw *ChanWitch) Add(name string, ch interface{}) error`
Adds a new channel with a name. Returns an error if the name already exists.

### `func (cw *ChanWitch) Get(name string) interface{}`
Retrieves a channel by its name. Returns `nil` if the channel does not exist.

### `func (cw *ChanWitch) Close(name string)`
Closes the specified channel and removes it from the list.

### `func (cw *ChanWitch) CloseAll()`
Closes all stored channels and clears the list.

### `func (cw *ChanWitch) String() string`
Returns a string representation of all stored channels.

### `func (cw *ChanWitch) Len() int`
Returns the number of channels in the list.

## License
This project is licensed under the MIT License. Feel free to use, modify, and distribute it.

## Contributions
Pull requests and issues are welcome! 🚀

## Author
Created by **Saman**.

