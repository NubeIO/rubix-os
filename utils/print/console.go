package pprint

import (
	"encoding/json"
	"fmt"
	"os"
)

func Print(i interface{}) {
	fmt.Printf("%+v\n", i)
	return
}

func Log(i interface{}) string {
	return fmt.Sprintf("%+v\n", i)
}

func PrintJSON(x interface{}) {
	ioWriter := os.Stdout
	w := json.NewEncoder(ioWriter)
	w.SetIndent("", "    ")
	w.Encode(x)
}
