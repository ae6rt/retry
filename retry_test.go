package retry

import (
	"errors"
	"testing"
	"time"
)

func TestHappyPath(t *testing.T) {
	r := New(0*time.Second, 3)
	err := r.Try(func() error {
		return nil
	})
	if err != nil {
		t.Fatalf("Unexpected error: %v\n", err)
	}
}

func TestRetryExceeded(t *testing.T) {
	r := New(0*time.Second, 3)
	tries := 0
	err := r.Try(func() error {
		tries += 1
		return errors.New("")
	})
	if err == nil {
		t.Fatalf("Expecting error\n")
	}
	if tries != 3 {
		t.Fatalf("Expecting 3 but got %d\n", tries)
	}
}

func TestTimeout(t *testing.T) {
	r := New(2*time.Second, 1)
	err := r.Try(func() error {
		time.Sleep(4 * time.Second)
		return nil
	})
	if err == nil {
		t.Fatalf("Expected error\n")
	}
	if _, ok := err.(timeoutError); !ok {
		t.Fatalf("Expected retry.timeoutError\n")
	}
}

func TestIsTimeout(t *testing.T) {
	if !isTimeout(timeoutError{}) {
		t.Fatalf("Expected timeout error\n")
	}
	if isTimeout(errors.New("")) {
		t.Fatalf("Unexpected error\n")
	}
}
