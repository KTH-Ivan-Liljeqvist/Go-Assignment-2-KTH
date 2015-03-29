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

	wgc := new(sync.WaitGroup)
	wgc.Add(consumers)

	for i := 0; i < producers; i++ {
		go Produce("p"+strconv.Itoa(i), strings/producers, ch, wgp)
	}
	for i := 0; i < consumers; i++ {
		go Consume("c"+strconv.Itoa(i), ch, wgc)
	}

	wgp.Wait() // Wait for all producers to finish.
	close(ch)
	wgc.Wait() // Wait for all consumers to finish.

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
func Consume(id string, ch <-chan string, wg *sync.WaitGroup) {
	for s := range ch {
		fmt.Println(id, "received", s)
		RandomSleep(100) // Simulate time to consume data.
	}

	wg.Done()
}

// RandomSleep waits for x ms, where x is a random number, 0 â‰¤ x < n,
// and then returns.
func RandomSleep(n int) {
	time.Sleep(time.Duration(rand.Intn(n)) * time.Millisecond)
}

/*
	1) Vad händer om man byter plats på satserna wgp.Wait() och close(ch) i slutet av main-funktionen?

	   Om man byter plats på close(ch) och wgp.Wait() i slutet av main metoden så stängs kanalen innan alla gorutiner
	   hinner göra sitt arbete. Det resulterar i att Produce-rutinen försöker skriva till kanalen ch då den redan är stängd.
	   Då får man panic eftersom man kan inte skriva till stängd kanal.

	2) Vad händer om man flyttar close(ch) från main-funktionen och i stället stänger kanalen i slutet av funktionen Produce?

	   Om man flyttar close(ch) till slutet av Produce hamnar vi i en situation där vi startar flera Produce rutiner. En av dem kommer då
	   att stänga kanalen medan dem andra fortfarande kan vara igång och vill skriva till den här kanalen som numera är stängd.
	   Därför får man panic även den här gången.
	   Man försöker skriva till stängd kanal.

	3) Vad händer om man tar bort satsen close(ch) helt och hållet?

	   Om man tar bort close(ch) fungerar programmet som det ska. Detta är för att close(ch) säger bara att inga värden kommer att skickas
	   på den kanalen mer. Eftersom programmet tar slut precis efter close(ch) spelar det ingen roll om man stänger kanalen.
	   Programmet kommer ändå sluta precis efter. Ingen kommer försöka skicka eller ta emot nåt från kanalen.

	4) Vad händer om man ökar antalet konsumenter från 2 till 4?

	   Om man ökar antalet konsumenter från 2 till 4 körs programmet snabbare. Detta är för att konsumenterna körs parallelt och läser av kanalen samtidigt. Dem gör arbetet
	   samtidigt och blir därför klara snabbare.

	5) Kan man vara säker på att alla strängar blir utskrivna innan programmet stannar?

	   Man kan vara säker på att programmet printar alla strängar eftersom vi använder en obuffrad kanal. Det innebär att varje producent väntar på att
	   en konsument ska printa det värdet som den producenten skrev till kanalen innan samma producent fortsätter att skriva nästa värde till kanalen.
	   Och programmet fortsätter köras och tar inte slut förrän alla producenter är klara eftersom vi har en wait-group som tar hand om det.
	   Producenterna kallar inte på wait-group:ens Done metod innan alla värden som den producenten skickade till kanalen blir utskrivna av någon konsument.

	6) Ändra koden genom att lägga till en ny WaitGroup som väntar tills alla konsumenter blivit klara. Lämna in koden för det modifierade programmet.

		Jag la till en WaitGroup (wgc) som väntar på att alla konsumenter ska bli klara innan programmet stannar.


*/
