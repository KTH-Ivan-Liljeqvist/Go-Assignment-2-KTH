# Go-Assignment-2-KTH

This assignment focused very much on synchronization, deadlocks and data races.

The first part of the assignment was fixing two bugs. The first bug had to do with channels and the second bug had to do with WaitGroups.
The bugs are fixed in the files bug1.go and bug2.go.

In the second part of the assignment I had to experiment with channels and WaitGroups. This part of the assignment required thorough understanding of unbuffered and buffered channels and their nature. 
This part is in the many2many.go.

The third part of the assignment was to build an "oracle"-program like ELIZA (http://en.wikipedia.org/wiki/ELIZA). The user can ask different questions and the Oracle will answer the questions. The oracle will also make different prophecies.

This had to be accomplished using three main goroutines.
First goroutine for capturing a question and starting another goroutine that will answer this question.
Second goroutine for making random prophecies.
Third goroutine for capturing and printing answers and prophecies to the console.

You could create other methods and goroutines if you wanted to.
