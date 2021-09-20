package main

import (
	"fmt"
	"math/rand"
	"time"
)

func Poll() {
	r := rand.New(rand.NewSource(99))
	c := time.Tick(1 * time.Second)
	for range c {
		//Download the current contents of the URL and do something with it
		fmt.Printf("Grab at %s\n", time.Now())
		// add a bit of jitter
		jitter := time.Duration(r.Int31n(100)) * time.Millisecond
		time.Sleep(jitter)
	}
}

func main() {
	//go obj.Poll()
	Poll()
}
