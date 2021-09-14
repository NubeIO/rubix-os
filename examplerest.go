package main

import (
	"fmt"
	"github.com/NubeDev/flow-framework/client"
	"github.com/NubeDev/flow-framework/model"
)

func main() {

	c := client.NewSessionWithToken("fakeToken123", "0.0.0.0", "1660")
	var writerCloneModel model.WriterClone
	writerCloneModel.ThingClass = "test"
	getPlat, err := c.EditWriterClone("wrc_c9736617dd2046d7", writerCloneModel, false)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(getPlat.ThingClass)

}
