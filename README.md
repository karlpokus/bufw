# bufw
A concurrency-safe io.Writer with internal buffer and sync props for go tests

# use case
You have some function you want to test that requires an `io.Writer`, like a network connection or a file descriptor, and you want to inspect what was written at a later point. Writes to bufw are safe for concurrent use. Use the sync feature If you want to await 1+ writes.

# usage
See tests

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
- [ ] go mod
- [ ] conf
- [ ] would passing a sync chan to New be more flexible?
- [ ] trim Write input
- [x] replace bytes.Buffer with a single []byte

# license
MIT
