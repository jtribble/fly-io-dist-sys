package syncutils

import (
	"sync"
	"time"
)

func WaitUntilTimeout(wg *sync.WaitGroup, timeout time.Duration) (timedOut bool) {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false
	case <-time.After(timeout):
		return true
	}
}
