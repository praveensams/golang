package main

import (
	"fmt"
	"regexp"
	"time"
)

type S struct {
	a int
	b int
}

func (s *S) finds(c chan string, i string) {
	if se, _ := regexp.MatchString(`\w{3}`, i); se {
		fmt.Println(10)
		s.a = s.a + 1
		c <- i
	}
}

func main() {
	d := &S{1, 2}
	c := make(chan string)
	go d.finds(c, "same")
	time.Sleep(2 * time.Second)
	fmt.Println(d)

}
