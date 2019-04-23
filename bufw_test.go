package bufw

import (
	"bytes"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	w := New(false)
	if w.written != nil {
		t.Errorf("expected %v to be nil", w.written)
	}
	w = New(true)
	if w.written == nil {
		t.Errorf("expected %v to be a chan", w.written)
	}
}

func TestSequential(t *testing.T) {
	w := New(false)
	input := []byte("hello")
	w.Write(input)
	output := w.Bytes()
	if !bytes.Equal(output, input) {
		t.Errorf("%s and %s are not equal", output, input)
	}
}

func TestString(t *testing.T) {
	w := New(false)
	input := []byte(" hello ")
	w.Write(input)
	output := w.String()
	inputTrimmed := string(bytes.TrimSpace(input))
	if output != inputTrimmed {
		t.Errorf("%s and %s are not equal", output, inputTrimmed)
	}
}

func TestWait(t *testing.T) {
	w := New(true)
	input := []byte("hello")
	go func() { w.Write(input) }()
	err := w.Wait()
	if err != nil {
		t.Errorf("Expected nil, got %s", err)
	}
	output := w.Bytes()
	if !bytes.Equal(output, input) {
		t.Errorf("%s and %s are not equal", output, input)
	}
}

func TestWaitTimeout(t *testing.T) {
	w := New(true)
	w.SyncTimeout("100ms")
	input := []byte("hello")
	go func() {
		time.Sleep(1 * time.Second)
		w.Write(input)
	}()
	err := w.Wait()
	if err != ErrTimeout {
		t.Errorf("Expected %s, got %s", ErrTimeout, err)
	}
	output := w.Bytes()
	var expected []byte
	if !bytes.Equal(output, expected) {
		t.Errorf("%s and %s are not equal", output, expected)
	}
}

func TestWaitN(t *testing.T) {
	w := New(true)
	input := []byte(" hello ")
	go func() { w.Write(input) }()
	go func() { w.Write(input) }()
	go func() { w.Write(input) }()
	n, err := w.WaitN(3)
	if err != nil {
		t.Errorf("Expected nil, got %s", err)
	}
	if n != 3 {
		t.Errorf("Expected %d, got %d", 3, n)
	}
	expected := []byte("hellohellohello")
	output := w.Bytes()
	if !bytes.Equal(output, expected) {
		t.Errorf("%s and %s are not equal", output, expected)
	}
}

func TestWaitNTimeout(t *testing.T) {
	w := New(true)
	w.SyncTimeout("100ms")
	input := []byte(" hello ")
	go func() { w.Write(input) }()
	go func() { w.Write(input) }()
	go func() {
		time.Sleep(1 * time.Second)
		w.Write(input)
	}()
	n, err := w.WaitN(3)
	if err != ErrTimeout {
		t.Errorf("Expected %s, got %s", ErrTimeout, err)
	}
	if n != 2 {
		t.Errorf("Expected %d, got %d", 0, n)
	}
	expected := []byte("hellohello")
	output := w.Bytes()
	if !bytes.Equal(output, expected) {
		t.Errorf("%s and %s are not equal", output, expected)
	}
}

func TestReset(t *testing.T) {
	w := New(false)
	input := []byte("hello")
	w.Write(input)
	output := w.Bytes()
	if !bytes.Equal(output, input) {
		t.Errorf("%s and %s are not equal", output, input)
	}
	input = []byte("bye")
	w.Write(input)
	output = w.Bytes()
	if !bytes.Equal(output, input) {
		t.Errorf("%s and %s are not equal", output, input)
	}
}
