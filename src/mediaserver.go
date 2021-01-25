package main

import (
		"flag"
		"github.com/out-of-mind/mediaserver/src/internal/ms"
)

var (
		addr string
		log_path string
		log_level string
)

func init() {
		flag.StringVar(&addr, "addr", "localhost:8002", "--addr 192.168.0.1:80 to listen 192.168.0.1:80")
		flag.StringVar(&log_path, "log-path", "mediaserver.log", "-log-path mediaserver.log to set mediaserver.log as main log file")
		flag.StringVar(&log_level, "log-level", "warning", "-log-level warning to set warning log level")
}

func main() {
		flag.Parse()
        mediaserver := ms.New(addr, log_path, log_level)
		mediaserver.Run()
}
