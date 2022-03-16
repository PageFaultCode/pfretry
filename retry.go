// Package retry is everything to retry a function
package retry

import (
	"fmt"
	"time"
)

// The retry function to retry
type RetryFunc func() error

// The error call back when an error occurs
type ErrorCallback func(retryName string, err error)

// The success callback when success happens
type SuccessCallback func(retryName string)

// PFRetry is a structure defining what retry takes to work
type PFRetry struct {
	retryCount       int
	retryName        string
	minimumRetryTime time.Duration
	stopFlag         chan struct{}
}

// SetRetryCount sets the retry count for the given
// retry option
func (r *PFRetry) SetRetryCount(retryCount int) {
	r.retryCount = retryCount
}

// SetRetryName sets the name of the retry to identify it
// when a callback occurs
func (r *PFRetry) SetRetryName(retryName string) {
	r.retryName = retryName
}

// Retry takes a function and retries it X times based on
// the set retryCount
func (r *PFRetry) Retry(function RetryFunc, onSuccess SuccessCallback, onError ErrorCallback) error {
	funcCount := 1
	for {
		select {
		case <-r.stopFlag:
			return fmt.Errorf("stop signal received")
		default:
			minimumTime := time.Now().Add(r.minimumRetryTime)
			err := function()
			if err != nil {
				if onError != nil {
					onError(r.retryName, err)
				}
			} else {
				if onSuccess != nil {
					onSuccess(r.retryName)
				}
				return nil
			}

			funcCount++
			if r.retryCount > 0 && funcCount > r.retryCount {
				return err
			}
			endTime := time.Now()
			if endTime.Before(minimumTime) {
				time.Sleep(minimumTime.Sub(endTime))
			}
		}
	}
}
