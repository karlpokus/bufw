# bufw
A concurrency-safe io.Writer with internal buffer and sync props for go tests

[![GoDoc](https://godoc.org/github.com/karlpokus/bufw?status.svg)](https://godoc.org/github.com/karlpokus/bufw)

# use case
You have some function you want to test that requires an `io.Writer`, like a network connection or a file, and you want to inspect what was written at a later point. Writes to bufw are safe for concurrent use. Use the sync feature to await n writes.

# example
```go
func Run() {
	w := New()
	thing := NewThing(w)
	go func() {
		thing.Do() // writes to w
	}()
	err := w.Wait() // proceed only after w has been written to
	if err != nil {
		log.Fatal("Wait timeout err", err)
	}
	b := w.Bytes() // read bytes from w
}
```

# install
```bash
$ go get github.com/karlpokus/bufw
```

# test
```bash
$ go test -v -cover -race
```

# todo
- [x] test coverage
- [ ] would passing a sync chan to New be more flexible?
- [x] trim Write input
- [x] replace bytes.Buffer with a single []byte
- [x] godoc
- [x] Wait timeout
- [x] return an error if Wait is used but written chan is nil
- [ ] let wait return buffer contents
- [x] remove sync opt

# license
MIT
