package main

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

var wg sync.Mutex
var vge sync.WaitGroup

type S struct {
	x int
}

type SlackRequestBody struct {
	Text string `json:"text"`
}

func (count *S) validate(i int, se string) {
	h := se
	addr := fmt.Sprintf("%s:%d", h, i)
	s, err := net.DialTimeout("tcp", addr, time.Second*2)
	if err == nil {
		count.x = count.x + 1
		ads := fmt.Sprintf("%s  -> %d", se, i)
		fmt.Println(ads)
		defer s.Close()

	}

	defer vge.Done()
}

func main() {

	count := &S{
		x: 0,
	}
	count.x = 0

	for i := 0; i < 1024; i++ {
		vge.Add(1)
		go count.validate(i, os.Args[1])

	}

	vge.Wait()

}
