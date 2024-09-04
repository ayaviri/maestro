package internal

import "sync"

// Represents a channel that be closed more than once safely
type SafeClosureChannel struct {
	C    chan any
	once sync.Once
}

func NewSafeClosureChannel() *SafeClosureChannel {
	return &SafeClosureChannel{C: make(chan any)}
}

func (c *SafeClosureChannel) SafeClose() {
	c.once.Do(func() {
		close(c.C)
	})
}
