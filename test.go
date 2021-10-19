package main

import (
	log "github.com/sirupsen/logrus"
	"time"
)

func main() {
	//var wg sync.WaitGroup
	log.Println("Starting Weekly Checks")

}

func forever() {
	for {
		//fmt.Printf("%v+\n", time.Now())
		time.Sleep(time.Second)
	}
}
