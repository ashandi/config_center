package application

import "sync"

type GracefulQuitFn func(fn func())

func NewAutoWaitGroup() (*sync.WaitGroup, GracefulQuitFn) {
	wg := &sync.WaitGroup{}

	return wg, func(fn func()) {
		wg.Add(1)
		fn()
		wg.Done()
	}
}
