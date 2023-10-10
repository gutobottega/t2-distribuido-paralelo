package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Import the Semaphore package
import "FPPDSemaforo"

var (
	data        int
	readers     int
	readMutex   sync.Mutex
	writeMutex  sync.Mutex
	readCounter sync.Mutex
	writeCount  int
	semaphore   *FPPDSemaforo.Semaphore // Semaphore for controlling access
)

func reader(id int) {
	for {
		// Wait for permission to enter the critical section
		semaphore.Wait()

		readCounter.Lock()
		readers++
		if readers == 1 {
			writeMutex.Lock() // Prevent writers when the first reader enters
		}
		readCounter.Unlock()

		// Read data
		fmt.Printf("Reader %d is reading data: %d\n", id, data)
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))

		readCounter.Lock()
		readers--
		if readers == 0 {
			writeMutex.Unlock() // Allow writers when the last reader exits
		}
		readCounter.Unlock()

		// Release the semaphore
		semaphore.Signal()

		// Simulate some work
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(200)))
	}
}

func writer(id int) {
	for {
		// Wait for permission to enter the critical section
		semaphore.Wait()

		// Write data
		writeMutex.Lock()
		writeCount++
		fmt.Printf("Writer %d is writing data: %d\n", id, writeCount)
		data = writeCount
		writeMutex.Unlock()

		// Release the semaphore
		semaphore.Signal()

		// Simulate some work
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(300)))
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	numReaders := 5
	numWriters := 2

	// Initialize the semaphore with an initial value of 1
	semaphore = FPPDSemaforo.NewSemaphore(1)

	for i := 1; i <= numReaders; i++ {
		go reader(i)
	}

	for i := 1; i <= numWriters; i++ {
		go writer(i)
	}

	// Keep the program running
	select {}
}
