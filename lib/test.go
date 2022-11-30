package main

import (
	"fmt"
	"time"
)

func amo() {

a:
	time.Sleep(time.Second * 1)
	fmt.Println("sadas")
	goto a
}

func main() {
	amo()
}
