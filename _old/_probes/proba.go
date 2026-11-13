package main

import (
	"log"
	"time"
)

func main() {
	t1 := time.Now()
	t2 := time.Now().Add(time.Second)

	log.Print(t1.Sub(t2))
}
