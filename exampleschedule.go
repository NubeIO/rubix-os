package main

import (
	"fmt"
	"github.com/NubeDev/flow-framework/schedule"
	"io/ioutil"
	"log"
	"time"
)

func TestOverlappingIntervals(weekday time.Weekday) (bool, []schedule.IntervalValue) {
	json, err := ioutil.ReadFile("./schedule/test_overlap.json")
	if err != nil {
		fmt.Println(err)
	}
	sd, err := schedule.New(string(json))
	if err != nil {
		fmt.Println("Unexpected error parsing json")
	}
	dt := time.Now()
	dt = time.Date(dt.Year(), dt.Month(), dt.Day(), 11, 00, 0, 0, time.Local)
	for {
		if dt.Weekday() == weekday {
			break
		}
		dt = dt.AddDate(0, 0, -1)
	}
	shifts := sd.MatchingIntervals(dt)
	if len(shifts) != 2 {
		return true, shifts
	}
	return false, shifts
}

func EventCheckLoop() {
	json, err := ioutil.ReadFile("./schedule/test.json")
	if err != nil {
		log.Fatal(err)
	}
	sd, err := schedule.New(string(json))
	if err != nil {
		log.Println("Unexpected error parsing json")
	}
	log.Println(sd.String())
	log.Println(sd.Intervals)
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	go func() {
		for dt := range ticker.C {
			if sd.Within(dt) {
				log.Println(dt, " is WITHIN schedule.")
			} else {
				log.Println(dt, " is OUTSIDE schedule.")
			}
		}
	}()
	time.Sleep(time.Hour * 1)
	ticker.Stop()
}

func main() {
	o, oo := TestOverlappingIntervals(time.Monday)
	fmt.Println(oo)
	fmt.Println(o, "is an overlap if its true")
	EventCheckLoop()

}
