package main

import "fmt"

func thread(ch chan int) {
	ch <- 1
	ch <- 2
	close(ch)
}

func main() {
	ch := make(chan int)
	go thread(ch)

	for i := 0; i < 3; i++ {
		select {
		case x, ok := <-ch:
			if ok {
				fmt.Println("ch", ":", x)
			} else {
				fmt.Println("canal cerrado")
			}
		}
	}
}
