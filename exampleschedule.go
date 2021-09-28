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

func EventCheckLoop(file string) {
	json, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	sd, err := schedule.New(string(json))
	if err != nil {
		log.Println("Unexpected error parsing json")
	}
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	go func() {
		for dt := range ticker.C {
			if sd.Within(dt) {
				log.Println(sd.GetDescription(), dt, " is TRUE schedule.")
			} else {
				log.Println(sd.GetDescription(), dt, " is FALSE schedule.")
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
	go EventCheckLoop("./schedule/test.json")
	go EventCheckLoop("./schedule/test2.json")
	go forever()
	select {} // block forever

}
func forever() {
	for {
		//fmt.Printf("%v+\n", time.Now())
		time.Sleep(time.Second)
	}
}
