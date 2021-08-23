package main

import (
	valid "github.com/asaskevich/govalidator"
)

// Model validate
type User struct {
	Email string `valid:"required,email"`
	FullName string `valid:"required,alpha"`
	Password string `valid:"required"`
	Phone string `valid:"numeric"`
}

// validate tùy biến thông báo lỗi với dấu "~"
type CustomUser struct {
	Email string `json:"email"  valid:"required,email~aaa 66666666666666666666666666666666666"`
	FullName string `valid:"required~Yêu cầu nhập tên"`
	Password string `valid:"required~Yêu cầu nhập mật khẩu"`
	Phone string `valid:"numeric"`
}

func main() {
	ValidStruct()
	ValidStructCustom()
	DemoValidMap()
}


// validate struct
func ValidStruct(){
	Mock := User{
		Email: "linhtrinhvietgmail.com",
		FullName: "linhtrinhviet",
		Password: "",
		Phone: "0369232329",
	}

	println("------validate struct-------")
	result ,err := valid.ValidateStruct(Mock)
	if err != nil {
		println("error: " + err.Error())
	}
	println(result)
}

func ValidStructCustom() {
	Mock := CustomUser{
		Email: "linhtrinhvietgmail.com",
		FullName: "linhtrinhviet",
		// Password: "123456",
		Phone: "0369232329",
	}

	println("------validate struct custom error msg-------")
	result ,err := valid.ValidateStruct(Mock)
	if err != nil {
		println("error: " + err.Error())
	}
	println(result)
}



// Validate Map

func DemoValidMap() {
	var mapTemplate = map[string]interface{}{
		"name":"required,alpha",
		"email":"required,email",
		"phone":"numeric",
		"address":map[string]interface{}{
			"line1":"required,alphanum",
			"line2":"alphanum",
			"postal-code":"numeric",
		},
	}

	var inputMap = map[string]interface{}{
		"name":"Linh",
		"email":"linhgtx@gmail.com",
		"phone" : "0369232329",
		"address":map[string]interface{}{
			"line1":"",
			"line2":"",
			"postal-code":"",
		},
	}

	println("------validate Map-------")
	result, err := valid.ValidateMap(inputMap, mapTemplate)
	if err != nil {
		println("error: " + err.Error())
	}
	println(result)
}
