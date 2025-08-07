package myWaitGroup

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func worker(integer *int, integerMutex *sync.Mutex) {
	integerMutex.Lock()
	defer integerMutex.Unlock()

	*integer++
}

func runWorkers(waitGroup *MyWaitGroup, count int, integer *int, integerMutex *sync.Mutex) {
	for i := 1; i <= count; i++ {
		waitGroup.Add(1)

		go func() {
			defer waitGroup.Done()
			worker(integer, integerMutex)
		}()

	}
	waitGroup.Wait()
}

func TestMyWaitGroup(t *testing.T) {

	var myWg MyWaitGroup
	var integer int
	var integerMutex sync.Mutex

	const firstIterationsCount = 5
	const secondIterationsCount = 10

	runWorkers(&myWg, firstIterationsCount, &integer, &integerMutex)
	runWorkers(&myWg, secondIterationsCount, &integer, &integerMutex)

	if integer != (firstIterationsCount + secondIterationsCount) {
		t.Errorf("An error has occurred during passes the test: integer want to %d, but have %d", (firstIterationsCount + secondIterationsCount), integer)
	}

	fmt.Printf("Looks like the test was passed.\n")
}

func TestNegativeCount(t *testing.T) {

	defer func() {
		if err := recover(); err == nil {
			t.Errorf("An error has occurred during passes the test: want to panic during negative value, but no panic has occurred.")
		}
	}()

	var myWg MyWaitGroup
	myWg.Add(-1)

}

func TestAddDuringWaiters(t *testing.T) {

	var myWg MyWaitGroup

	fakeGoroutines := 1
	myWg.Add(1)

	const maxTimerWait = 30

	go func() {

		defer func() {
			if err := recover(); err != nil {
				for {
					myWg.Done()
					fakeGoroutines--

					if fakeGoroutines == 0 {
						break
					}
				}
			}
		}()

		start := time.Now()

		for {

			if time.Since(start).Seconds() >= maxTimerWait {
				t.Errorf("An error has occurred during passes the test: want to panic during add goroutines when waiters an existed, but no panic has occurred.")
			}

			myWg.Add(1)
			fakeGoroutines++
		}
	}()

	myWg.Wait()

}
