package main

import (
	"fmt"
	"sync"
	"time"
)

type Sam struct {
	a, b, k int
}

var wg sync.WaitGroup

func (s *Sam) runs(c chan int, w int) {

	defer wg.Done()

	s.k = s.k + 1

	c <- s.a + s.b*w

}

func main() {
	d := &Sam{1, 20, 0}
	c := make(chan int)
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go d.runs(c, i)

	}

	time.Sleep(30 * time.Second)
	fmt.Println(d.k)
	for i := 0; i < d.k; i++ {
		fmt.Println(<-c)
	}

	wg.Wait()
}
