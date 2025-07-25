package myWaitGroup

import (
	"fmt"
	"sync"
	"testing"
)

func worker(integer *int, integerMutex *sync.Mutex) {
	integerMutex.Lock()
	defer integerMutex.Unlock()

	*integer++
}

func TestMyWaitGroup(t *testing.T) {

	var myWg MyWaitGroup
	var integer int
	var integerMutex sync.Mutex

	const firstIterationsCount = 5
	const secondIterationsCount = 10

	for i := 1; i <= firstIterationsCount; i++ {
		myWg.Add(1)

		go func() {
			defer myWg.Done()
			worker(&integer, &integerMutex)
		}()

	}

	myWg.Wait()

	for i := 1; i <= secondIterationsCount; i++ {
		myWg.Add(1)

		go func() {
			defer myWg.Done()
			worker(&integer, &integerMutex)
		}()

	}

	myWg.Wait()

	if integer != (firstIterationsCount + secondIterationsCount) {
		t.Errorf("An error has occurred during passes the test: integer want to %d, but have %d", (firstIterationsCount + secondIterationsCount), integer)
	}

	fmt.Printf("Looks like the test was passed.\n")
}
