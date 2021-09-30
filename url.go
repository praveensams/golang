package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
)

type S struct {
	a int
}

func (se S) val(site string, c chan string) {

	s, e := http.Get(site)
	if e == nil {

		g, _ := io.ReadAll(s.Body)
		if v, _ := regexp.MatchString("google", string(g)); v {
			se.a = se.a + 1
			c <- "Site is there" + "  ->  " + site
		}
	}
}

func main() {
	v := &S{1}
	c := make(chan string, 3)
	lists := []string{
		"http://www.google.com",
		"http://yahoo.com",
		"http://www.gmail.com",
	}
	for i := 0; i < len(lists); i++ {
		go v.val(lists[i], c)

	}
	for i := 0; i < v.a; i++ {
		fmt.Println(<-c)
	}
}
