package main

import (
	"fmt"
	"runtime"
	"time"
)

// Computation ...
type Computation struct {
	Value       int
	ElapsedTime time.Duration
}

// Work ...
type Work struct {
	Limit  int
	Job    chan int
	Output chan *Computation
}

// Worker ...
type Worker interface {
	CloseChannels()
	Process(cores int)
	FormatOutput() string
}

func main() {
	availableCores := getAvailableCPUCores()

	fmt.Printf("Multicore Fibonacci.\n\nDisclaimer:\n\tDoes not support thread synchronization yet. So, output\n\torder is not guaranteed whenever using more than one core.\n\n")

	cores, upTo := getUserInput(availableCores)
	work := &Work{Limit: upTo, Job: make(chan int, upTo), Output: make(chan *Computation, upTo)}

	initial := time.Now()

	work.Process(cores)
	work.FormatOutput()
	work.CloseChannels()

	fmt.Printf("\nElapsed time: %s\n", getElapsedTime(initial).Round(time.Second).String())
}

// CloseChannels ...
func (w *Work) CloseChannels() {
	close(w.Job)
	close(w.Output)
}

// Process ...
func (w *Work) Process(cores int) {
	// Creates n user specified concurrency processes.
	for i := 0; i <= cores; i++ {
		go calculateForEachWork(w.Job, w.Output)
	}

	// Injects each number into the channel.
	for n := 0; n <= w.Limit; n++ {
		w.Job <- n
	}
}

// FormatOutput ...
func (w *Work) FormatOutput() {
	for j := 0; j <= w.Limit; j++ {
		output := <-w.Output

		fmt.Printf("#%d (%s): %d\n", j, output.ElapsedTime.Round(time.Millisecond).String(), output.Value)
	}
}

func getUserInput(availableCores int) (cores int, upTo int) {
	var userCores, userNumberLimit int

	fmt.Printf("%d CPU cores available.\n\nHow many CPU core should it run with? (1 to %d): ", availableCores, availableCores)

	_, err := fmt.Scan(&userCores)

	if err != nil {
		panic(fmt.Sprintf("Can not execute program with %d cores. Must be betwen 1 and %d.", userCores, availableCores))
	}

	if userCores < 1 || userCores > 4 {
		panic(fmt.Sprintf("Can not execute program with %d cores. Must be betwen 1 and %d.", userCores, availableCores))
	}

	fmt.Printf("How many Fibonacci iterations would you like? (1 to max(int32)): ")

	_, err = fmt.Scan(&userNumberLimit)

	if err != nil {
		panic("Could not read the max number of executions.")
	}

	if upTo < 0 {
		panic(fmt.Sprintf("Can not execute program with %d value.", userNumberLimit))
	}

	fmt.Print("\n")

	return userCores - 1, userNumberLimit
}

func getElapsedTime(initial time.Time) time.Duration {
	return time.Now().Sub(initial)
}

func getAvailableCPUCores() int {
	return runtime.NumCPU()
}

func calculateForEachWork(jobs <-chan int, results chan<- *Computation) {
	// As slots become available, calculate each Fibonnaci number.
	for number := range jobs {
		initial := time.Now()
		result := &Computation{fib(number), getElapsedTime(initial)}

		results <- result
	}
}

func fib(number int) int {
	if number <= 1 {
		return number
	}

	return fib(number-1) + fib(number-2)
}
