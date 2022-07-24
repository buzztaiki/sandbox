package main

import (
	"flag"
	"log"
	"os"
	"time"
)

var duration = flag.Int("duration", 5, "")
var count = flag.Int("count", 5, "")
var noRecover = flag.Bool("no-recover", false, "dont recover from unhealthy state")

func main() {
	flag.Parse()

	c := time.Tick(time.Duration(*duration) * time.Second)
	ok := true
	n := 0
	for next := range c {
		log.Println(next, ok)
		if ok {
			if _, err := os.Create("/tmp/healthy"); err != nil {
				log.Fatal("failed to create:", err)
			}
		} else {
			if err := os.RemoveAll("/tmp/healthy"); err != nil {
				log.Fatal("failed to remove:", err)
			}
		}

		if ok || !*noRecover {
			n++
			if n == *count {
				n = 0
				ok = !ok
			}
		}
	}
}
