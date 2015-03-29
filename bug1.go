package main

import "fmt"

/*
	The problem was that the unbuffered channel we created waited for some go routine to take the data, but because no other
	routine was active it just paused the program and waited.

	I solved the problem by creating a go routine that inserts the value into the channel.
	The value is then extracted from the channel by the main method.

	Another solution could be to make the channel buffered.
*/

func main() {

	//create the channel
	ch := make(chan string)

	//start a goroutine and place the value into the channel
	go func() {
		ch <- "Hello world!"
	}()

	//extract the value and print to the console.
	fmt.Println(<-ch)
}
