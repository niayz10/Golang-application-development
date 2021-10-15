package main

import (
    "fmt"
    "sync"
    "time"
)

func worker(id int, jobs <-chan int, results chan<- int) {
    for j:= range jobs {
        fmt.Printf("Worker %d starting job %d\n", id, j)
        time.Sleep(time.Second)
        fmt.Printf("Worker %d done job %d\n", id, j)
        results <- j*5
    }  
}

func main() {

    var wg sync.WaitGroup
    jobs := make(chan int, 5)
    results := make(chan int, 5)

    for i := 1; i <= 3; i++ {
        wg.Add(1)

        i := i

        go func() {
            defer wg.Done()
            worker(i, jobs, results)
        }()
    }

    for j := 1 ; j<=5; j++{
        jobs<-j
    }
    close(jobs)
    wg.Wait()

}
