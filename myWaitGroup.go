package myWaitGroup

import (
	"sync/atomic"
	"unsafe"
)

type waiter struct {
	semaphore chan struct{}
	next      unsafe.Pointer
}

type waitersList struct {
	head  unsafe.Pointer
	count uint64
}

type MyWaitGroup struct {
	counter int64
	waiters waitersList
}

func (w *MyWaitGroup) Add(count int) {

	waitersCount := atomic.LoadUint64(&w.waiters.count)

	if count > 0 && waitersCount > 0 {
		panic("myWaitGroup already contains waiting goroutines!")
	}

	newValue := atomic.AddInt64(&w.counter, int64(count))
	if newValue < 0 {
		panic("myWaitGroup counter has become negative!")
	}

	if newValue == 0 {

		for {

			head := atomic.LoadPointer(&w.waiters.head)
			if head == nil {
				return
			}

			next := atomic.LoadPointer(&(*waiter)(head).next)

			if atomic.CompareAndSwapPointer(&w.waiters.head, head, next) {
				(*waiter)(head).semaphore <- struct{}{}
				close((*waiter)(head).semaphore)
			}

		}

	}

}

func (w *MyWaitGroup) Done() {
	w.Add(-1)
}

func (w *MyWaitGroup) Wait() {

	newWaiter := &waiter{semaphore: make(chan struct{}, 1)}

	for {

		head := atomic.LoadPointer(&w.waiters.head)
		newWaiter.next = head

		if atomic.CompareAndSwapPointer(&w.waiters.head, head, unsafe.Pointer(newWaiter)) {
			break
		}

	}

	atomic.AddUint64(&w.waiters.count, 1)

	<-newWaiter.semaphore

	for {
		waitersCount := atomic.LoadUint64(&w.waiters.count)

		if atomic.CompareAndSwapUint64(&w.waiters.count, waitersCount, waitersCount-1) {
			return
		}

	}

}
