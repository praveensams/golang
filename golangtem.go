package main

import (
	"fmt"
	"log"
	"os"
	"text/template"
)

type S struct {
	A, B int
}

func main() {
	d := S{1, 2}
	s, e := template.New("sam").Parse("My name is {{.B}}\n")
	if e == nil {
		fmt.Println("bye")
	}
	err := s.Execute(os.Stdout, d)
	if err == nil {
		log.Print("testing")
	}
}
