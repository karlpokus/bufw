// Package bufw provides a concurrency-safe io.Writer with internal buffer and sync props
package bufw

import (
	"bytes"
	"sync"
)

// Bufw implements io.Writer
// safe for concurrent use
type Bufw struct {
	mu      sync.Mutex
	buf     []byte
	written chan bool
}

// Write writes whitespace trimmed bytes to an internal buffer
// also writes to the written chan if sync is enabled
func (bw *Bufw) Write(b []byte) (n int, err error) {
	bw.mu.Lock()
	defer bw.mu.Unlock()
	bw.buf = append(bw.buf, bytes.TrimSpace(b)...)
	if bw.written != nil {
		bw.written <- true
	}
	return len(b), nil
}

// Bytes resets-, and returns the buffer contents
func (bw *Bufw) Bytes() []byte {
	bw.mu.Lock()
	defer bw.mu.Unlock()
	b := bw.buf
	bw.buf = nil
	return b
}

// String returns buf as a string and resets buf
func (bw *Bufw) String() string {
	return string(bw.Bytes())
}

// Wait blocks on the written chan until a Write is performed. Used for synchronization
// A must use if sync is enabled or the Write call will block
func (bw *Bufw) Wait() {
	<-bw.written
	return
}

// WaitN blocks on the written chan until n Writes are performed
func (bw *Bufw) WaitN(n int) {
	for i := 0; i < n; i++ {
		<-bw.written
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
