package hooks

import (
	"context"
	"log"
	"time"
)

func CreateFixedIntervals(times uint, interval time.Duration) []time.Duration {
	intervals := make([]time.Duration, 0, times)
	var i uint
	for i = 1; i <= times; i++ {
		intervals = append(intervals, interval)
	}
	return intervals
}

func CreateLinearIntervals(times uint, interval time.Duration) []time.Duration {
	intervals := make([]time.Duration, 0, times)
	var i uint
	for i = 1; i <= times; i++ {
		intervals = append(intervals, time.Duration(i)*interval)
	}
	return intervals
}

func CreateFibonacciIntervals(times uint, interval time.Duration) []time.Duration {
	intervals := make([]time.Duration, 0, times)
	for _, val := range createFibonacciInts(times) {
		intervals = append(intervals, time.Duration(val)*interval)
	}
	return intervals
}

func createFibonacciInts(max uint) []uint {
	if max == 0 {
		panic("max should not be zero")
	}
	if max == 1 {
		return []uint{1}
	} else if max == 2 {
		return []uint{2}
	} else {
		fs := make([]uint, 0, max)
		fs = append(fs, 1, 2)
		var i uint
		for i = 3; i <= max; i++ {
			fs = append(fs, fs[i-3]+fs[i-2])
		}
		return fs
	}
}

func RunWithRetry(ctx context.Context, fn func(context.Context) error, retryIntervals []time.Duration) (err error) {
	log.Println("first run ...")
	// first run
	err = fn(ctx)
	if err == nil {
		return
	}
	// retry logic
	for _, interval := range retryIntervals {
		select {
		case <-ctx.Done():
			log.Println("cancel run ...")
			err = context.Canceled
			return
		case <-time.After(interval):
			log.Println("retry run ...")
			err = fn(ctx)
			if err == nil || err == context.Canceled || err == context.DeadlineExceeded {
				return
			}
		}
	}
	return
}
