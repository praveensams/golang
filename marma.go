package main

import "sync"

type S struct {
	a, b int
	c    int
}

func (s S) adds(c chan<- int) {
	wge.Lock()
	c <- s.a + s.b
	wge.Unlock()
	wg.Done()
}

var wg sync.WaitGroup
var wge sync.Mutex

func main() {
	d := S{1, 2, 0}
	c := make(chan int, 100)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		//	d.a = i
		d.adds(c)
	}
	wg.Wait()
}
