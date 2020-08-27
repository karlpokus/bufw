// Package bufw provides a concurrency-safe io.Writer with internal buffer and sync props
package bufw

import (
	"bytes"
	"errors"
	"sync"
	"time"
)

var ErrTimeout = errors.New("timeout")

// Bufw implements io.Writer. Safe for concurrent use
type Bufw struct {
	sync.Mutex
	buf     []byte
	written chan bool
	ttl     time.Duration
}

// Write writes whitespace trimmed bytes to an internal buffer
// and writes to the written chan
func (w *Bufw) Write(b []byte) (n int, err error) {
	w.Lock()
	defer w.Unlock()
	w.buf = append(w.buf, bytes.TrimSpace(b)...)
	w.written <- true
	return len(b), nil
}

// Bytes resets-, and returns the buffer contents
func (w *Bufw) Bytes() []byte {
	w.Lock()
	defer w.Unlock()
	b := w.buf
	w.buf = nil
	return b
}

// String returns buf as a string and resets buf
func (w *Bufw) String() string {
	return string(w.Bytes())
}

// Wait blocks on the written chan until a Write is performed or a timeout occurs.
// A timeout will return an error.
func (w *Bufw) Wait() error {
	timer := time.NewTimer(w.ttl)
	select {
	case <-timer.C:
		return ErrTimeout
	case <-w.written:
		timer.Stop()
		return nil
	}
}

// WaitN blocks on the written chan until n Writes are performed or a timeout occurs.
// Returns an error and number of successful writes performed.
func (w *Bufw) WaitN(n int) (int, error) {
	for i := 0; i < n; i++ {
		err := w.Wait()
		if err != nil {
			return i, err
		}
	}
	return n, nil
}

/// SyncTimeout sets the timeout for the Wait and WaitN funcs
func (w *Bufw) SyncTimeout(ttl string) error {
	d, err := time.ParseDuration(ttl)
	w.ttl = d
	return err
}

// New returns a Bufw type with a default wait timeout of 10s
func New() *Bufw {
	d, _ := time.ParseDuration("10s")
	return &Bufw{
		ttl: d,
		written: make(chan bool),
	}
}
