package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func counter(counterCh chan int, init int) {
	counter := init

	rand := rand.Intn(10)

	for i := init; i < init+rand; i++ {
		counterCh <- counter
		fmt.Println("Counter: ", counter)
		counter += 1
		time.Sleep(1 * time.Second)
	}
}

func backupCounter(counterCh chan int, wg *sync.WaitGroup) {
	var count int
	var init int
	lastMessage := time.Now()

LOOP:
	for {
		select {
		case receive := <-counterCh:
			count = receive
			lastMessage = time.Now()
		default:
			if time.Since(lastMessage) > 2*time.Second {
				break LOOP
			}
		}
	}

	fmt.Println("server error - backup encounters")
	init = count + 1
	wg.Add(2)

	go func() {
		counter(counterCh, init)
		wg.Done()
	}()

	go backupCounter(counterCh, wg)

	wg.Done()
}

func main() {
	counterCh := make(chan int)
	var wg sync.WaitGroup

	wg.Add(2)
	counterInit := 0

	go func() {
		counter(counterCh, counterInit)
		wg.Done()
	}()

	go backupCounter(counterCh, &wg)

	wg.Wait()
}
