package main

import (
	"fmt"
	"regexp"
	"sync"
)

var wg sync.WaitGroup

type S struct {
	d int
	s string
}

func (se *S) fin(ch chan string, i string) {
	defer wg.Done()
	if s, e := regexp.MatchString(`sa\w{3,4}\b`, i); s && e == nil {

		se.d = se.d + 1

		ch <- i
	}

}

func main() {
	c := make(chan string, 100)

	f := S{}

	d := []string{"samer", "damer", "jamer", "aamer", "laamer", "sameee", "samwwww"}

	for i := 0; i < len(d); i++ {
		wg.Add(1)
		go f.fin(c, d[i])
	}
	wg.Wait()
	//time.Sleep(5 * time.Second)
	for j := 0; j < f.d; j++ {
		fmt.Println(<-c)
	}

}
