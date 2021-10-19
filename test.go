package main

import (
	"fmt"
	"github.com/NubeDev/flow-framework/src/schedule"
	log "github.com/sirupsen/logrus"
	"time"
)

func main() {
	//var wg sync.WaitGroup
	log.Println("Starting Weekly Checks")
	check, err := schedule.WeeklyCheck("./src/schedule/weekly_schedule.json", "TEST")
	if err != nil {
		return
	}
	fmt.Println(check.IsActive)
}

func forever() {
	for {
		//fmt.Printf("%v+\n", time.Now())
		time.Sleep(time.Second)
	}
}
