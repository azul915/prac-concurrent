package main

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"
)

const (
	Green = 0
	Red   = 1
)

func main() {

	lgr, _ := zap.NewDevelopment()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	listInList := [][]bool{{true, true, true}, {true, false, true}, {false, false, true}}
	signal := make(chan int, 8)
	for _, v := range listInList {
		AllOK(lgr, v, signal)

		select {
		case sig := <-signal:
			switch {
			case sig == Green:
				lgr.Debug("Green")
			case sig == Red:
				lgr.Debug("Red")
			}
		case <-ctx.Done():
			lgr.Debug("Timeout happend.")
		}
	}

}

func AllOK(lgr *zap.Logger, li []bool, signal chan<- int) {

	var status int

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wg.Done()
		status = Status(lgr, li)
	}()
	wg.Wait()

	if status == 0 {
		signal <- Green
	} else {
		signal <- Red
	}
}

func Status(lgr *zap.Logger, li []bool) int {

	status := 0

	wg := new(sync.WaitGroup)
	for _, v := range li {
		wg.Add(1)
		go func(v bool) {
			defer wg.Done()
			if !v {
				status = 1
			}
		}(v)
	}
	wg.Wait()
	return status
}
