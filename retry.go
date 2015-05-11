package retry

import (
	"log"
	"os"
	"time"
)

type timeoutError struct {
	error
}

func (t timeoutError) String() string {
	return "retry.timeout"
}

var logger *log.Logger = log.New(os.Stdout, "retry ", log.Ldate|log.Ltime|log.Lshortfile)

type Retry struct {
	Timeout     time.Duration
	MaxAttempts int
}

func New(timeout time.Duration, maxAttempts int) Retry {
	if maxAttempts < 1 {
		maxAttempts = 1
	}
	return Retry{Timeout: timeout, MaxAttempts: maxAttempts}
}

func (r Retry) Try(work func() error) error {
	done := make(chan struct {})
	errorChan := make(chan error)
	attempts := 0

	var expired <-chan time.Time
	if r.Timeout > 0 {
		timer := time.NewTimer(r.Timeout)
		expired = timer.C
		defer timer.Stop()
	}

	for {
		go func() {
			attempts += 1
			if err := work(); err != nil {
				errorChan <- err
			} else {
				done <- struct {}{}
			}
		}()

		select {
		case err := <-errorChan:
			if attempts == r.MaxAttempts {
				return err
			}
		case <-done:
			return nil
		case <-expired:
			logger.Println("timeout")
			return timeoutError{}
		}
	}
}

func isTimeout(err error) bool {
	if _, ok := err.(timeoutError); ok {
		return true
	}
	return false
}
