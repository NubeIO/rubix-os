package utils

import (
	"bytes"
	"encoding/json"
	"os"
)

func WriteDataToFileAsJSON(data interface{}, fileDIR string) (int, error) {
	//write data as buffer to json encoder
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "\t")

	err := encoder.Encode(data)
	if err != nil {
		return 0, err
	}
	file, err := os.OpenFile(fileDIR, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return 0, err
	}
	n, err := file.Write(buffer.Bytes())
	if err != nil {
		return 0, err
	}
	return n, nil
}

func CreateFile(f string) (ok bool, err error) {
	// detect if file exists
	_, err = os.Stat(f)
	if err != nil {
		return false, err
	}
	// create file if not exists
	if os.IsNotExist(err) {
		file, err := os.Create(f)
		defer file.Close()
		if err != nil {
			return false, err
		}
	}
	return true, err
}

func DeleteFile(f string) (ok bool, err error) {
	err = os.Remove(f)
	if err != nil {
		return false, err
	}
	return true, err
}
