package main

import (
	"fmt"
	"time"
)

type Sam struct {
	a int
}

func (s *Sam) add() {
	s.a = s.a + 1
}

func main() {
	k := &Sam{}
	go k.add()
	go k.add()
	fmt.Println(k)
	t0 := time.Now().String()[:25]
	fmt.Println(t0)
}
