package main

import (
		"github.com/out-of-mind/mediaserver/src/internal/ms"
)

func main() {
        mediaserver := ms.New(":8002", "localhost")
		mediaserver.Run()
}
