package bufw

import (
	"bytes"
	"sync"
	"strings"
)

// bufw implements io.Writer
// safe for concurrent use
type bufw struct {
	mu      sync.Mutex
	buf     bytes.Buffer
	written chan bool
}

// Write writes to an internal buffer for later inspection
// also writes to the written chan if it has been instantiated
func (bw *bufw) Write(b []byte) (n int, err error) {
	bw.mu.Lock()
	defer bw.mu.Unlock()
	n, err = bw.buf.Write(b)
	if err != nil {
		if bw.written != nil {
			bw.written <- true
		}
		return n, err
	}
	if bw.written != nil {
		bw.written <- true
	}
	return len(b), nil
}

// Bytes resets-, and returns the buffer contents
func (bw *bufw) Bytes() []byte {
	bw.mu.Lock()
	defer bw.mu.Unlock()
	b := bw.buf.Bytes()
	bw.buf.Reset()
	return b
}

// String returns buf as a string and resets buf
func (bw *bufw) String() string {
	b := bw.Bytes()
	return strings.TrimSpace(string(b))
}

// Wait blocks on the written chan until a Write is performed. Used for synchronization
// A must use if sync is enabled or the Write call will block
func (bw *bufw) Wait() bool {
	return <-bw.written
}

// WaitN blocks on the written chan until n Writes are performed
func (bw *bufw) WaitN(n int) {
	for i := 0; i < n; i++ {
		<-bw.written
	}
}

// New returns a bufw type and instantiates the written chan if enableSync is true
func New(enableSync bool) *bufw {
	w := &bufw{}
	if enableSync {
		w.written = make(chan bool)
	}
	return w
}