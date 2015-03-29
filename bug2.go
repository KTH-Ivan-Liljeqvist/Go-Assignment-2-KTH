package main

import "fmt"
import "sync"

/*
The problem was that the program ended before all of the numbers were printed to the console.
I fixed the bug by adding a wait group. Whenever a number is printed to the console, Done is called
on the wait group. This way you know that all of the numbers will be printed before the program ends.
*/

func main() {

	const NUMBER_OF_ITERATIONS = 11

	//create the waitgroup
	waitGroup := new(sync.WaitGroup)
	waitGroup.Add(NUMBER_OF_ITERATIONS)

	//create the channel
	ch := make(chan int)

	//start the routine which will print all numbers in the channel
	go Print(ch, waitGroup)

	//add numbers to the channel, these numbers will be printed by the go routine Print we started above
	for i := 1; i <= NUMBER_OF_ITERATIONS; i++ {
		ch <- i

	}

	//wait for all the numbers to be printed to the console
	waitGroup.Wait()
	//close the channel
	close(ch)

}

// Print prints all numbers sent on the channel.
// The function returns when the channel is closed.
func Print(ch <-chan int, waitGroup *sync.WaitGroup) {

	for n := range ch { // reads from channel until it's closed
		fmt.Println(n)
		waitGroup.Done()
	}
}
