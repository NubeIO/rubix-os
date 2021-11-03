package main

import (
	"encoding/json"
	"fmt"
	"github.com/NubeDev/flow-framework/model"
	"io"
	"io/ioutil"
	"net/http"
)

func getProducerHistories(url string) ([]*model.History, error) {
	var histories []*model.History
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("error: getProducerHistories")
		return histories, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			fmt.Println("error: getProducerHistories")
		}
	}(res.Body)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("error: getProducerHistories")
		return histories, err
	}
	if err = json.Unmarshal(body, &histories); err != nil {
		fmt.Println("error: getProducerHistories")
		return histories, err
	}
	return histories, err
}
