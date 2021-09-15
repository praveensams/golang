package main

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

type S struct {
	x int
}

var wg sync.Mutex
var vge sync.WaitGroup

func (count *S) validate(i int, se string, c chan string) {
	h := se
	addr := fmt.Sprintf("%s:%d", h, i)
	s, err := net.DialTimeout("tcp", addr, time.Second*2)
	if err == nil {
		count.x = count.x + 1
		ads := fmt.Sprintf("%s  -> %d", se, i)
		c <- ads
		defer s.Close()

	}

	defer vge.Done()
}

func main() {
	count := &S{
		x: 0,
	}
	c := make(chan string, 10)
	count.x = 0
	for j := 1; j < len(os.Args); j++ {
		for i := 0; i < 1024; i++ {
			vge.Add(1)
			go count.validate(i, os.Args[j], c)
		}
	}

	vge.Wait()
	time.Sleep(4 * time.Second)
	for j := 0; j < count.x; j++ {
		fmt.Println(<-c)
	}
}
