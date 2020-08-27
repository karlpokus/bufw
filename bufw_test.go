package bufw

import (
	"bytes"
	"testing"
)

func TestWait(t *testing.T) {
	w := New()
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
	w := New()
	w.SyncTimeout("100ms")
	err := w.Wait()
	if err != ErrTimeout {
		t.Errorf("Expected %s, got %s", ErrTimeout, err)
	}
}

func TestWaitN(t *testing.T) {
	w := New()
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
	w := New()
	w.SyncTimeout("100ms")
	input := []byte(" hello ")
	go func() { w.Write(input) }()
	go func() { w.Write(input) }()
	n, err := w.WaitN(3)
	if err != ErrTimeout {
		t.Errorf("Expected %s, got %s", ErrTimeout, err)
	}
	if n != 2 {
		t.Errorf("Expected %d, got %d", 2, n)
	}
	expected := []byte("hellohello")
	output := w.Bytes()
	if !bytes.Equal(output, expected) {
		t.Errorf("%s and %s are not equal", output, expected)
	}
}

func TestReset(t *testing.T) {
	w := New()
	input := []byte("hello")
	go func() { w.Write(input) }()
	w.Wait()
	output := w.Bytes()
	if !bytes.Equal(output, input) {
		t.Errorf("%s and %s are not equal", output, input)
	}
	input = []byte("bye")
	go func() { w.Write(input) }()
	w.Wait()
	output = w.Bytes()
	if !bytes.Equal(output, input) {
		t.Errorf("%s and %s are not equal", output, input)
	}
}
