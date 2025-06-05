package myWaitGroup

import (
	"fmt"
	"testing"
	"time"
)

func worker(id int) {
	fmt.Printf("Worker %d is starting\n", id)

	time.Sleep(time.Second)
	fmt.Printf("Worker %d is done\n", id)
}

func TestMyWaitGroup(t *testing.T) {

	var myWg MyWaitGroup

	for i := 1; i <= 5; i++ {
		myWg.Add(1)

		go func() {
			defer myWg.Done()
			worker(i)
		}()

	}

	myWg.Wait()

	for i := 5; i <= 10; i++ {
		myWg.Add(1)

		go func() {
			defer myWg.Done()
			worker(i)
		}()

	}

	myWg.Wait()

	fmt.Printf("Looks like the test was passed.")
}
