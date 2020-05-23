package main

import (
	"errors"
	"fmt"
	"runtime"
	"time"
)

// Worker ...
type Worker interface {
	CloseChannels()
	Process(cores int)
}

// Result ...
type Result struct {
	Value       interface{}
	ElapsedTime time.Duration
}

// Work ...
type Work struct {
	Input   interface{}
	Jobs    chan int
	Outputs chan *Result
}

// CloseChannels ...
func (w *Work) CloseChannels() {
	close(w.Jobs)
	close(w.Outputs)
}

// Process ...
func (w *Work) Process(cores int) error {
	switch w.Input.(type) {
	case int:
	default:
		return errors.New("Can not process strings")
	}

	// Creates n user specified concurrency processes.
	for i := 0; i <= cores; i++ {
		go calculateForEachWork(w.Jobs, w.Outputs)
	}

	// Injects each number into the channel.
	for n := 0; n <= w.Input.(int); n++ {
		w.Jobs <- n
	}

	// Retrieves results as they are available.
	for j := 0; j <= w.Input.(int); j++ {
		<-w.Outputs
	}

	return nil
}

func main() {
	availableCores := getAvailableCPUCores()

	fmt.Printf("Multicore Fibonacci.\n\nDisclaimer:\n\tDoes not support thread synchronization yet. So, output order\n\tis not guaranteed whenever using more than one core.\n")

	cores, upTo := getUserInput(availableCores)
	work := &Work{Input: upTo, Jobs: make(chan int, upTo), Outputs: make(chan *Result, upTo)}

	diff := createMeasurementTimer()

	if err := work.Process(cores); err == nil {
		work.CloseChannels()

		fmt.Printf("\nElapsed time: %s\n", diff().Round(time.Second).String())
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

	return userCores, userNumberLimit
}

func createMeasurementTimer() func() time.Duration {
	start := time.Now()

	return func() time.Duration {
		return time.Now().Sub(start)
	}
}

func getAvailableCPUCores() int {
	return runtime.NumCPU()
}

func calculateForEachWork(jobs <-chan int, results chan<- *Result) {
	// As slots become available, calculate each Fibonnaci number.
	for number := range jobs {
		diff := createMeasurementTimer()
		fib, _ := fib(number)
		result := &Result{fib, diff()}

		results <- result
	}
}

func fib(number int) (int, error) {
	if number < 0 {
		return 0, errors.New("Can not calculate Fibonacci with negative numbers")
	}

	if number <= 1 {
		return number, nil
	}

	result, first, second := 0, 0, 1

	for i := 2; i <= number; i++ {
		result = first

		first, second = second, first+second
	}

	return result, nil
}
