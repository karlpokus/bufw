package bufw

import (
	"bytes"
	"testing"
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
	w.Wait()
	output := w.Bytes()
	if !bytes.Equal(output, input) {
		t.Errorf("%s and %s are not equal", output, input)
	}
}

func TestWaitN(t *testing.T) {
	w := New(true)
	input := []byte(" hello ")
	go func() { w.Write(input) }()
	go func() { w.Write(input) }()
	go func() { w.Write(input) }()
	w.WaitN(3)
	expected := []byte("hellohellohello")
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
