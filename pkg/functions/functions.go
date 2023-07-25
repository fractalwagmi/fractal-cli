package functions

import "time"

func RunWithRetryPolicy(backoffSchedule []time.Duration, fn func() error) error {
	var err error

	for _, backoff := range backoffSchedule {
		err = fn()
		if err == nil {
			break
		}

		time.Sleep(backoff)
	}

	return err
}
