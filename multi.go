package main

import (
	"fmt"
	"time"
)

type S struct {
	a int
}

func (se S) ping(ch1 chan<- int, s int) {
	if s%2 == 0 {
		se.a = se.a + 1
		ch1 <- s

	}
}

func main() {
	ch1 := make(chan int)
	se := &S{0}

	for i := 0; i < 100; i++ {
		go se.ping(ch1, i)

	}
	time.Sleep(6 * time.Second)
	fmt.Println(se.a)

}
