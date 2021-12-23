package squel

import "sync"

var opts = options{
	Debug: false,
}

type options struct {
	Debug bool
	mu    sync.RWMutex
}

func SetDebug(value bool) {
	opts.mu.Lock()
	defer opts.mu.Unlock()
	opts.Debug = value
}
