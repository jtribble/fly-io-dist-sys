package maelstrom

import (
	"sync"
	"time"
)

type Callback func(response *Message)

type Callbacks interface {
	Register(msg *Message, callback Callback)

	Retrieve(inReplyTo *int) Callback

	RetryableMessages(threshold time.Duration) []*Message
}

func NewCallbacks() Callbacks {
	return &defaultCallbacks{
		callbacks: make(map[int]registration),
		mu:        sync.Mutex{},
	}
}

type defaultCallbacks struct {
	callbacks map[int]registration
	mu        sync.Mutex
}

type registration struct {
	msg      *Message
	callback Callback
	lastSent time.Time
}

func (c *defaultCallbacks) Register(msg *Message, callback Callback) {
	if callback == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.callbacks[*msg.Body.MsgId] = registration{
		msg:      msg,
		callback: callback,
		lastSent: time.Now(),
	}
}

func (c *defaultCallbacks) Retrieve(inReplyTo *int) Callback {
	if inReplyTo == nil {
		return nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if entry, ok := c.callbacks[*inReplyTo]; ok {
		delete(c.callbacks, *inReplyTo)
		return entry.callback
	}
	return nil
}

func (c *defaultCallbacks) RetryableMessages(threshold time.Duration) (messages []*Message) {
	now := time.Now()
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, c := range c.callbacks {
		if time.Since(c.lastSent) > threshold {
			messages = append(messages, c.msg)
			c.lastSent = now
		}
	}
	return
}
