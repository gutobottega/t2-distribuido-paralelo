package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var (
	data        int
	readers     int
	readMutex   sync.Mutex
	writeMutex  sync.Mutex
	readCounter sync.Mutex
	writeCount  int
	semaphore   *Semaphore // Semaphore for controlling access
)

func reader(id int) {
	for {
		// Espera pela perminssão para entrar na sessão critica
		semaphore.Wait()

		readCounter.Lock()
		readers++
		if readers == 1 {
			writeMutex.Lock() // Trava os writers quando o primeiro reader entrar
		}
		readCounter.Unlock()

		fmt.Printf("Reader %d is reading data: %d\n", id, data)
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))

		readCounter.Lock()
		readers--
		if readers == 0 {
			writeMutex.Unlock() // Permite leitura quando o ultimo reader existe
		}
		readCounter.Unlock()

		// Libera o semaforo
		semaphore.Signal()

		time.Sleep(time.Millisecond * time.Duration(rand.Intn(200)))
	}
}

func writer(id int) {
	for {
		// Espera pela permissão para entrar na sessão critica
		semaphore.Wait()

		// Escreve o conteúdo
		writeMutex.Lock()
		writeCount++
		fmt.Printf("Writer %d is writing data: %d\n", id, writeCount)
		data = writeCount
		writeMutex.Unlock()

		// Libera o semaforo
		semaphore.Signal()

		time.Sleep(time.Millisecond * time.Duration(rand.Intn(300)))
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	numReaders := 5
	numWriters := 2

	// Inicializa o semaforo com 1
	semaphore = NewSemaphore(1)

	for i := 1; i <= numReaders; i++ {
		go reader(i)
	}

	for i := 1; i <= numWriters; i++ {
		go writer(i)
	}

	// Mantem o programa rodando
	select {}
}
	