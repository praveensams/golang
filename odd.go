package main

import (
	"fmt"
	"time"
)

type Num struct {
	a int
	b int
}

func (n *Num) find(ch chan<- int, i int) {
	if i%2 == 0 {
		n.a = n.a + 1
		ch <- i

	}
}

func recovers() {
	if r := recover(); r != nil {
		fmt.Println("hiii")
	}
}

func main() {
	ch := make(chan int)
	d := &Num{0, 0}
	for i := 0; i < 100; i++ {
		go d.find(ch, i)
	}

	time.Sleep(3 * time.Second)
	for i := 0; i < d.a; i++ {
		je, ok := <-ch
		if !ok {
			fmt.Println(je)
		}

	}

}
