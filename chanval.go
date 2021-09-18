package main

import (
	"log"
)

//var wg sync.Mutex

type S struct {
	b int
	a int
}

func (s S) adding(ch chan S, i int) {
	//wg.Lock()
	s.b = s.b + i
	ch <- s
	//wg.Unlock()
}

func main() {
	che := make(chan S)
	lst := []int{1, 2, 3, 4, 5, 6}
	st := S{b: 20}

	for _, i := range lst {
		go st.adding(che, i)
	}
	for range lst {
		df := <-che
		log.Printf("%+v\n", df.b)
	}
}
