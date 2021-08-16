package main

import (
	"fmt"
	"github.com/NubeDev/flow-framework/helpers"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"log"
)

func main() {

	//MAKE CLIENT
	client := resty.New()
	//MAKE TOKEN
	getToken, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(`{"name":"admin"}`).
		SetBasicAuth("admin", "admin").
		Post("http://0.0.0.0:8888/client")
	if err != nil {
		log.Println("getToken err:", err, getToken.Status())
	}
	log.Println("getToken:", getToken, "status", getToken.Status())

	r := gjson.Get(string(getToken.Body()), "token")
	token := r.Str

	//GET TOKEN
	user, err := client.NewRequest().
		SetHeader("Authorization", token).
		Get("http://0.0.0.0:8888/user")
	if err != nil {
		log.Println("user err:", err, user.Status())
	}
	log.Println("user:", user, "status", user.Status())

	name, _ := helpers.MakeUUID()
	name = fmt.Sprintf("name_%s", name)

	//ADD NETWORK
	addNetwork, err := client.NewRequest().
		SetHeader("Authorization", token).
		SetBody(map[string]interface{}{"name": name, "description": "description"}).
		Post("http://0.0.0.0:8888/api/networks")
	if err != nil {
		log.Println("addNetwork err:", err, addNetwork.Status())
	}
	log.Println("addNetwork:", addNetwork, "status", addNetwork.Status())

	//ADD NETWORK
	r = gjson.Get(string(addNetwork.Body()), "uuid")
	getNetworkUUID := r.Str
	log.Println("getNetworkUUID:", getNetworkUUID)
	getNetwork, err := client.NewRequest().
		SetHeader("Authorization", token).
		SetPathParams(map[string]string{
			"uuid": getNetworkUUID,
		}).
		Get("http://0.0.0.0:8888/api/networks/{uuid}")
	if err != nil {
		log.Println("addNetwork err:", err, getNetwork.Status())
	}
	log.Println("addNetwork:", getNetwork, "status", getNetwork.Status())

	//EDIT NETWORK
	log.Println("getNetworkUUID:", getNetworkUUID)
	editNetwork, err := client.NewRequest().
		SetHeader("Authorization", token).
		SetBody(map[string]interface{}{"name": "new_name_" + name}).
		SetPathParams(map[string]string{
			"uuid": getNetworkUUID,
		}).
		Patch("http://0.0.0.0:8888/api/networks/{uuid}")
	if err != nil {
		log.Println("editNetwork err:", err, editNetwork.Status())
	}
	log.Println("editNetwork:", editNetwork, "status", editNetwork.Status())

	//DELETE NETWORK
	log.Println("getNetworkUUID:", getNetworkUUID)
	deleteNetwork, err := client.NewRequest().
		SetHeader("Authorization", token).
		SetPathParams(map[string]string{
			"uuid": getNetworkUUID,
		}).
		Delete("http://0.0.0.0:8888/api/networks/{uuid}")
	if err != nil {
		log.Println("deleteNetwork err:", err, deleteNetwork.Status())
	}
	log.Println("deleteNetwork:", deleteNetwork, "status", deleteNetwork.Status())

	if getToken.Status() == "200 OK" {
		fmt.Println("getToken", "PASS")
	} else {
		fmt.Println("getToken", "FAIL")
	}
	if user.Status() == "200 OK" {
		fmt.Println("user", "PASS")
	} else {
		fmt.Println("user", "FAIL")
	}
	if addNetwork.Status() == "200 OK" {
		fmt.Println("addNetwork", "PASS")
	} else {
		fmt.Println("addNetwork", "FAIL")
	}
	if getNetwork.Status() == "200 OK" {
		fmt.Println("getNetwork", "PASS")
	} else {
		fmt.Println("getNetwork", "FAIL")
	}
	if editNetwork.Status() == "200 OK" {
		fmt.Println("editNetwork", "PASS")
	} else {
		fmt.Println("editNetwork", "FAIL")
	}
	if deleteNetwork.Status() == "200 OK" {
		fmt.Println("deleteNetwork", "PASS")
	} else {
		fmt.Println("deleteNetwork", "FAIL")
	}

}
