package retry

import (
	"log"
	"os"
	"time"
)

var logger *log.Logger = log.New(os.Stdout, "retry ", log.Ldate|log.Ltime|log.Lshortfile)

type timeoutError struct {
	error
}

func (t timeoutError) String() string {
	return "retry.timeout"
}

type Retry struct {
	timeout     time.Duration
	maxAttempts int
}

func New(timeout time.Duration, maxAttempts int) Retry {
	if maxAttempts < 1 {
		maxAttempts = 1
	}
	return Retry{timeout: timeout, maxAttempts: maxAttempts}
}

func (r Retry) Try(work func() error) error {
	done := make(chan struct{})
	errorChan := make(chan error)
	attempts := 0

	var expired <-chan time.Time
	if r.timeout > 0 {
		timer := time.NewTimer(r.timeout)
		expired = timer.C
		defer timer.Stop()
	}

	for {
		go func() {
			attempts += 1
			if err := work(); err != nil {
				errorChan <- err
			} else {
				done <- struct{}{}
			}
		}()

		select {
		case err := <-errorChan:
			if attempts == r.maxAttempts {
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
