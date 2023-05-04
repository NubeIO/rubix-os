package main

import (
	"fmt"
	"testing"
	"time"
)

func sl() {
	time.Sleep(5 * time.Second)
	fmt.Println("Sleep Over.....")
}

func Test_timeout(t *testing.T) {
	err := Await(1000*time.Millisecond, 4000*time.Millisecond, func() bool {
		sl()
		return true
	})

	if err != nil {
		fmt.Println(err)
	}

}
