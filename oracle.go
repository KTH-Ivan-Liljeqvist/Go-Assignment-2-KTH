// Stefan Nilsson 2013-03-13
// Modified and further developed by Ivan Liljeqvist 2015-03-29

// This program implements an ELIZA-like oracle (en.wikipedia.org/wiki/ELIZA).
package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

const (
	star   = "Pythia"
	venue  = "Delphi"
	prompt = "> "
)

func main() {
	fmt.Printf("Welcome to %s, the oracle at %s.\n", star, venue)
	fmt.Println("Your questions will be answered in due time.")

	oracle := Oracle()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(prompt)
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fmt.Printf("%s heard: %s\n", star, line)
		oracle <- line // The channel doesn't block.
	}
}

// Oracle returns a channel on which you can send your questions to the oracle.
// You may send as many questions as you like on this channel, it never blocks.
// The answers arrive on stdout, but only when the oracle so decides.
// The oracle also prints sporadic prophecies to stdout even without being asked.
func Oracle() chan<- string {
	//channel for questions
	questions := make(chan string)
	//channel for answers
	answers := make(chan string)

	//routine for handling the questions and answering them
	go handleTheQuestions(questions, answers)
	//routine for making random prophecies
	go makeRandomProphecies(answers)

	//routine for printing answers and prophecies to the console.
	go handleOutput(answers)

	return questions
}

/*
	This function writes random prophecies to the channel 'answers' passed in  as parameter.
	It does so forever with a random pause in between the prophecies.
*/

func makeRandomProphecies(answers chan<- string) {

	for {
		time.Sleep(time.Duration(rand.Intn(20)+5) * time.Second)
		prophecy("", answers)
	}
}

/*
	This function reads values from channel 'answers' passed in as parameter and prints them to the console.
*/

func handleOutput(answers <-chan string) {
	for message := range answers {
		fmt.Println("The Oracle has something to say: ", message)
		fmt.Print(prompt)
	}
}

/*
	This function takes two channels as parameters - 'questions' to recieve questions from and 'answers' to write the answers.
	As long as 'questions' is open the method will answer the questions and write the answers to the channel 'answers'
*/
func handleTheQuestions(questions <-chan string, answers chan<- string) {

	//as long as the question channel is open, answer questions
	for q := range questions {
		go prophecy(q, answers)
	}
}

// This is the oracle's secret algorithm.
// It waits for a while and then sends a message on the answer channel.
// TODO: make it better.
func prophecy(question string, answer chan<- string) {
	// Keep them waiting. Pythia, the original oracle at Delphi,
	// only gave prophecies on the seventh day of each month.
	time.Sleep(time.Duration(20-(15)+rand.Intn(10)) * time.Second)

	// Find the longest word.
	longestWord := ""
	words := strings.Fields(question) // Fields extracts the words into a slice.
	for _, w := range words {
		if len(w) > len(longestWord) {
			longestWord = w
		}
	}

	// Cook up some pointless nonsense.
	nonsense := []string{
		"The moon is dark.",
		"The sun is bright.",
	}
	answer <- longestWord + "... " + nonsense[rand.Intn(len(nonsense))]
}

func init() { // Functions called "init" are executed before the main function.
	// Use new pseudo random numbers every time.
	rand.Seed(time.Now().Unix())
}
