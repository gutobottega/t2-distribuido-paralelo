package main

import (
	"fmt"
	"math/rand"
	"time"

	"FPPDSemaforo/FPPDSemaforo"
)

var (
	elves       = 0
	reindeer    = 0
	santaSem    = FPPDSemaforo.NewSemaphore(0)
	reindeerSem = FPPDSemaforo.NewSemaphore(0)
	elfTex      = FPPDSemaforo.NewSemaphore(1)
	mutex       = FPPDSemaforo.NewSemaphore(1)
)

func main() {
	go santaClaus()
	for i := 0; i < 3; i++ {
		go elf(i)
	}
	for i := 0; i < 9; i++ {
		go reindeerFunc(i)
	}

	time.Sleep(10 * time.Second) // Espera um tempo para que as goroutines terminem
}

func elf(id int) {
	elfTex.Wait()
	mutex.Wait()
	elves += 1
	if elves == 3 {
		santaSem.Signal()
	} else {
		elfTex.Signal()
	}
	mutex.Signal()

	getHelp(id)

	mutex.Wait()
	elves -= 1
	if elves == 0 {
		elfTex.Signal()
	}
	mutex.Signal()
}

func reindeerFunc(id int) {
	mutex.Wait()
	reindeer += 1
	if reindeer == 9 {
		santaSem.Signal()
	}
	mutex.Signal()
	reindeerSem.Wait()
	getHitched(id)
}

func santaClaus() {
	// Tirar uma soneca
	fmt.Println("Papai Noel está tirando uma soneca")
	time.Sleep(time.Duration(rand.Intn(5)) * time.Second)

	for {
		santaSem.Wait()
		mutex.Wait()
		if reindeer == 9 {
			prepareSleigh()
			for i := 0; i < 9; i++ {
				reindeerSem.Signal()
			}
		} else if elves == 3 {
			helpElves()
		}
		mutex.Signal()
	}
}

func prepareSleigh() {
	fmt.Println("Papai Noel está preparando o trenó")
	time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
}

func helpElves() {
	fmt.Println("Papai Noel está ajudando os elfos")
	time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
}

func getHitched(id int) {
	fmt.Printf("Rena %d está sendo engatada\n", id)
	time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
}

func getHelp(id int) {
	fmt.Printf("Elfo %d está recebendo ajuda\n", id)
	time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
}
