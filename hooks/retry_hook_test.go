package hooks

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestCreateFibonacciInts(t *testing.T) {
	ints := createFibonacciInts(10)
	t.Logf("%v", ints)
}

func TestCreateFinonacciIntervals(t *testing.T) {
	ints := CreateFibonacciIntervals(10, time.Second)
	t.Logf("%v", ints)
}

func TestCreateFixedIntervals(t *testing.T) {
	ints := CreateFixedIntervals(10, time.Second)
	t.Logf("%v", ints)
}

func TestCreateLinearIntervals(t *testing.T) {
	ints := CreateLinearIntervals(10, time.Second)
	t.Logf("%v", ints)
}

func print(ctx context.Context) (err error) {
	return errors.New("must fail")
}

func TestRetryFail(t *testing.T) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	go func() {
		<-time.After(time.Second * 20)
		cancelFunc()
	}()
	err := RunWithRetry(ctx, print, CreateFibonacciIntervals(10, time.Second*3))
	t.Logf("run error= %v", err)
}
