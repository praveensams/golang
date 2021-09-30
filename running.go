package main

import (
	"html/template"
	"io"
	"os"
)

type S struct {
	A, B int
}

var f io.Writer

func (s S) runs() {
	d := template.Must(template.ParseFiles("sam.html"))
	d.Execute(os.Stdout, s)
}

func main() {
	//c := make(chan io.Writer)
	sd := S{1, 2}
	for i := 0; i < 20; i++ {
		sd.A = i

		sd.runs()
	}

}
