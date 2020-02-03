# Gracefully Exit Go's HTTP Server

> Blog post - [gollahalli.com/blog/stopping-http-server-gracefully-context-vs-channels-vs-syncgroup/](https://www.gollahalli.com/blog/stopping-http-server-gracefully-context-vs-channels-vs-syncgroup/)

Some example to exit Go's HTTP server from a different function

## Requirements

1. Gorilla's Mux - [github.com/gorilla/mux](https://github.com/gorilla/mux)

## Usage

Create a `main.go` file and add the following:

```go
package main

import (
	"github.com/akshaybabloo/gracefully-exit-go-http-server/withchannels"
	"github.com/akshaybabloo/gracefully-exit-go-http-server/withcontext"
	"github.com/akshaybabloo/gracefully-exit-go-http-server/withsyncgroup"
)

func main() {
    // Run them one at a time

	withchannels.StartServer()
	withcontext.StartServer()
	withsyncgroup.StartServer()
}

```
