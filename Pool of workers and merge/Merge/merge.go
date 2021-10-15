package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func Merge(cs ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup
	wg.Add(len(cs))
	for _, c := range cs {
		go func(c <-chan int) {
			for v := range c {
				out <- v
			}
			wg.Done()
		}(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func InitializationChan(vs ...int) <-chan int {
	c := make(chan int)
	go func() {
		for _, v := range vs {
			c <- v
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
		close(c)
	}()
	return c
}

func main(){
	a := InitializationChan(0, 1, 2, 3, 4)
	b := InitializationChan(5, 6, 7, 8, 9)
	c := InitializationChan(10, 11, 12, 13, 14, 15)
	for v := range Merge(a, b, c) {
		fmt.Println(v)
	}
}