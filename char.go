package main

import (
	"fmt"
	"regexp"
	"sync"
	"time"
)

var wg sync.WaitGroup

type Sam struct {
	b int
	c int
}

func (a *Sam) find(c chan string, i string) {
	defer wg.Done()
	if d, _ := regexp.MatchString(`e\w{2}\b`, i); d {
		a.b = a.b + 1
		c <- i
	}
}

func main() {
	ch := make(chan string, 100)
	d := &Sam{}
	g := []string{"sam", "bam", "cam", "jam", "jem", "ext", "ext", "ext", "ext", "mmv", "mmv", "ext", "ext", "ext", "ext", "ext", "ext", "ext", "uls", "ext", "uls", "ext", "ext", "ext", "ext", "ext", "set", "rel", "php", "amp", "ext", "tmh", "uls", "amp", "amp", "src", "php", "amp", "amp", "amp", "raw", "amp", "rel", "php", "amp", "amp", "amp", "wmf", "org", "png", "rel", "org", "rel", "and", "max", "org", "rel", "png", "rel", "ico", "rel", "xml", "php", "rel", "rsd", "xml", "org", "api", "php", "rsd", "rel", "org", "rel", "org", "rel", "dns", "org", "rel", "dns", "org", "ltr", "ltr", "elt", "div", "div", "div", "div", "div", "top", "div", "div", "div", "div", "img", "alt", "src", "org", "svg", "svg", "png", "org", "svg", "svg", "png", "org", "svg", "svg", "png", "div", "div", "ogg", "img", "alt", "src", "org", "svg", "svg", "png", "org", "svg", "svg", "png", "org", "svg", "svg", "png", "div", "div", "div", "div", "the", "div", "div", "div", "div", "div", "div", "nav", "div", "div", "ltr", "dir", "ltr", "div", "div", "can", "div", "div", "top", "div", "not", "the", "see", "div", "elt", "div", "div", "png", "img", "alt", "src", "org", "png", "png", "org", "png", "png", "org", "png", "png", "div", "div", "png", "div", "div", "div", "div", "div", "div", "div", "div", "org", "jpg", "src", "org", "amp", "amp", "src", "org", "amp", "amp", "src", "org", "amp", "amp", "src", "org", "amp", "amp", "src", "org", "amp", "amp", "src", "org", "amp", "amp", "src", "org", "amp", "amp", "src", "org", "amp", "amp", "src", "org", "amp", "amp", "src", "org", "amp", "amp", "src", "org", "amp", "amp", "src", "org", "amp", "amp", "src", "org", "amp", "amp", "src", "org", "amp", "amp", "src", "org", "api", "php", "amp", "amp", "amp", "amp", "amp", "amp", "srt", "amp", "amp", "srt", "dir", "ltr", "src", "org", "api", "php", "amp", "amp", "amp", "amp", "amp", "amp", "srt", "amp", "amp", "srt", "dir", "ltr", "src", "org", "api", "php", "amp", "amp", "amp", "amp", "amp", "amp", "srt", "amp", "amp", "srt", "dir"}
	for i := 0; i < len(g); i++ {
		wg.Add(1)
		go d.find(ch, g[i])
	}
	wg.Wait()
	time.Sleep(2 * time.Second)
	for j := 0; j < d.b; j++ {
		fmt.Println(<-ch)
	}
}
