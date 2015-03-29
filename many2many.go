// Stefan Nilsson 2013-03-13

// This is a testbed to help you understand channels better.
package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

func main() {
	// Use different random numbers each time this program is executed.
	rand.Seed(time.Now().Unix())

	const strings = 32
	const producers = 4
	const consumers = 2

	before := time.Now()
	ch := make(chan string)
	wgp := new(sync.WaitGroup)
	wgp.Add(producers)
	for i := 0; i < producers; i++ {
		go Produce("p"+strconv.Itoa(i), strings/producers, ch, wgp)
	}
	for i := 0; i < consumers; i++ {
		go Consume("c"+strconv.Itoa(i), ch)
	}

	wgp.Wait() // Wait for all producers to finish.
	close(ch)
	fmt.Println("time:", time.Now().Sub(before))
}

// Produce sends n different strings on the channel and notifies wg when done.
func Produce(id string, n int, ch chan<- string, wg *sync.WaitGroup) {
	for i := 0; i < n; i++ {
		RandomSleep(100) // Simulate time to produce data.
		ch <- id + ":" + strconv.Itoa(i)
	}
	wg.Done()
}

// Consume prints strings received from the channel until the channel is closed.
func Consume(id string, ch <-chan string) {
	for s := range ch {
		fmt.Println(id, "received", s)
		RandomSleep(100) // Simulate time to consume data.
	}
}

// RandomSleep waits for x ms, where x is a random number, 0 â‰¤ x < n,
// and then returns.
func RandomSleep(n int) {
	time.Sleep(time.Duration(rand.Intn(n)) * time.Millisecond)
}

/*
	1) Om man byter plats på close(ch) och wgp.Wait() i slutet av main metoden så stängs kanalen innan alla gorutiner
	   hinner göra sitt arbete. Det resulterar i att Produce-rutinen försöker skriva till kanalen ch då den redan är stängd.
	   Då får man panic eftersom man kan inte skriva till stängd kanal.

	2) Om man flyttar close(ch) till slutet av Produce hamnar vi i en situation där vi startar flera Produce rutiner. En av dem kommer då
	   att stänga kanalen medan dem andra fortfarande kan vara igång och vill skriva till den här kanalen som numera är stängd.
	   Därför får man panic även den här gången.
	   Man försöker skriva till stängd kanal.

	3) Om man tar bort close(ch) fungerar programmet som det ska. Detta är för att close(ch) säger bara att inga värden kommer att skickas
	   på den kanalen mer. Eftersom programmet tar slut precis efter close(ch) spelar det ingen roll om man stänger kanalen.
	   Programmet kommer ändå sluta precis efter. Ingen kommer försöka skicka eller ta emot nåt från kanalen.
*/
