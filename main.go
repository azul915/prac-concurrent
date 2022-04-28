package main

import (
	"sync"

	"go.uber.org/zap"
)

// 特定の値の場合他のgoroutineをやめるor捨てるs

const (
	Green = 0
	Red   = 1
)

func main() {

	lgr, _ := zap.NewDevelopment()

	for _, v := range [][]bool{{true, true, true}, {true, false, true}, {false, false, true}} {

		go func(v []bool) {
			if AllOK(lgr, v) {
				lgr.Info("AllOK")
			} else {
				lgr.Info("Something Error happend")
			}
		}(v)
	}

	// time.Sleep(3 * time.Second)

}

func AllOK(logger *zap.Logger, bl []bool) bool {
	return Status(logger, bl) == Green
}

func Status(lgr *zap.Logger, list []bool) int {
	signal := Green

	wg := new(sync.WaitGroup)
	for _, v := range list {

		wg.Add(1)

		go func(v bool) {
			defer wg.Done()
			if !v {
				signal = Red
			}
		}(v)
	}
	wg.Wait()

	return signal
}
