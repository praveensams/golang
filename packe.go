package main

import (
	"fmt"
	"net"
)

func main() {
	d, err := net.DialTimeout("tcp", "www.google.com:80", 10)
	if err == nil {
		fmt.Sprintf("hiii")
	}
	if 3 > 7 {
		fmt.Println(d)
	}
}
