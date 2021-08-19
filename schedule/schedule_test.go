package schedule

import (
	"io/ioutil"
	"testing"
	"time"
)

func TestSchedule(t *testing.T) {
	json, err := ioutil.ReadFile("./test.json")
	if err != nil {
		t.Fatal(err)
	}
	sd, err := New(string(json))
	if err != nil {
		t.Fatal(err)
	}
	t.Log("TestSchedule: ", sd)
	if len(sd.Intervals) != 22 {
		t.Error("Expected 22 intervals")
	}
}

func TestLocateInSchedule(t *testing.T) {
	json, err := ioutil.ReadFile("./test.json")
	if err != nil {
		t.Fatal(err)
	}
	sd, err := New(string(json))
	if err != nil {
		t.Error("Unexpected error parsing json")
	}
	// rewind the current time to Monday
	dt := time.Now()
	dt = time.Date(dt.Year(), dt.Month(), dt.Day(), 8, 0, 0, 0, time.Local)
	for {
		if dt.Weekday() == time.Monday {
			break
		}
		dt = dt.AddDate(0, 0, -1)
	}
	if !sd.Within(dt) {
		t.Errorf("Expected to find %s in Schedule", dt)
	}
}

func TestLocateNotInSchedule(t *testing.T) {
	json, err := ioutil.ReadFile("./test.json")
	if err != nil {
		t.Fatal(err)
	}
	sd, err := New(string(json))
	if err != nil {
		t.Error("Unexpected error parsing json")
	}
	// rewind the current time to Monday
	dt := time.Now()
	dt = time.Date(dt.Year(), dt.Month(), dt.Day(), 7, 59, 0, 0, time.Local)
	for {
		if dt.Weekday() == time.Monday {
			break
		}
		dt = dt.AddDate(0, 0, -1)
	}
	if sd.Within(dt) {
		t.Errorf("Expected to NOT find %s in Schedule", dt)
	}
}

func TestSearchYear(t *testing.T) {
	json, err := ioutil.ReadFile("./test.json")
	if err != nil {
		t.Fatal(err)
	}
	sd, err := New(string(json))
	if err != nil {
		t.Error("Unexpected error parsing json")
	}
	// rewind the current time to Monday
	dt := time.Now()
	dt = time.Date(dt.Year(), dt.Month(), dt.Day(), 10, 15, 0, 0, time.Local)
	for {
		if dt.Weekday() == time.Monday {
			break
		}
		dt = dt.AddDate(0, 0, -1)
	}
	for i := 0; i < 366; i++ {
		if dt.Weekday() == time.Sunday {
			if sd.Within(dt) {
				t.Errorf("Expected to NOT find %s in Schedule", dt)
			}
		} else if !sd.Within(dt) {
			t.Errorf("Expected to find %s in schedule", dt)
		}
		dt = dt.AddDate(0, 0, 1)
	}
}

func TestTick(t *testing.T) {
	json, err := ioutil.ReadFile("./test.json")
	if err != nil {
		t.Fatal(err)
	}
	sd, err := New(string(json))
	if err != nil {
		t.Error("Unexpected error parsing json")
	}
	ticker := time.NewTicker(time.Minute)
	go func() {
		for dt := range ticker.C {
			if sd.Within(dt) {
				t.Log(dt, " is WITHIN schedule.")
			} else {
				t.Log(dt, " is OUTSIDE schedule.")
			}
		}
	}()

	time.Sleep(time.Second * 10)
	ticker.Stop()
}