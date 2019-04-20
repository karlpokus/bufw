// Package bufw provides a concurrency-safe io.Writer with internal buffer and sync props
package bufw

import (
	"bytes"
	"sync"
)

// Bufw implements io.Writer
// safe for concurrent use
type Bufw struct {
	sync.Mutex
	buf     []byte
	written chan bool
}

// Write writes whitespace trimmed bytes to an internal buffer
// also writes to the written chan if sync is enabled
func (w *Bufw) Write(b []byte) (n int, err error) {
	w.Lock()
	defer w.Unlock()
	w.buf = append(w.buf, bytes.TrimSpace(b)...)
	if w.written != nil {
		w.written <- true
	}
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

// Wait blocks on the written chan until a Write is performed. Used for synchronization
// A must use if sync is enabled or the Write call will block
func (w *Bufw) Wait() {
	<-w.written
	return
}

// WaitN blocks on the written chan until n Writes are performed
func (w *Bufw) WaitN(n int) {
	for i := 0; i < n; i++ {
		w.Wait()
	}
}

// New returns a Bufw type and instantiates the written chan if enableSync is true
func New(enableSync bool) *Bufw {
	w := &Bufw{}
	if enableSync {
		w.written = make(chan bool)
	}
	return w
}
