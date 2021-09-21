module github.com/NubeDev/flow-framework

require (
	github.com/NubeDev/configor v0.0.2
	github.com/NubeDev/location v0.0.2
	github.com/NubeIO/nubeio-rubix-lib-helpers-go v0.0.7
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d
	github.com/eclipse/paho.mqtt.golang v1.3.5
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/gin v1.7.4
	github.com/go-co-op/gocron v1.7.0
	github.com/go-playground/validator/v10 v10.4.1
	github.com/go-resty/resty/v2 v2.6.0
	github.com/gobuffalo/packr/v2 v2.7.1
	github.com/gorilla/websocket v1.4.2
	github.com/h2non/filetype v1.1.1
	github.com/mustafaturan/bus/v3 v3.0.3
	github.com/mustafaturan/monoton/v2 v2.0.2
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	golang.org/x/crypto v0.0.0-20210616213533-5ff15b29337e
	gopkg.in/yaml.v2 v2.4.0
	gorm.io/datatypes v1.0.1
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.21.13
)

require github.com/NubeIO/null v4.0.1+incompatible

require (
	github.com/brocaar/lora-app-server v2.5.1+incompatible
	github.com/pkg/errors v0.9.1
	github.com/simonvetter/modbus v1.3.0
	go.bug.st/serial v1.3.2
	google.golang.org/grpc v1.40.0
)

require (
	github.com/brocaar/loraserver v2.5.0+incompatible // indirect
	github.com/brocaar/lorawan v0.0.0-20210809075358-95fc1667572e // indirect
	github.com/go-gormigrate/gormigrate/v2 v2.0.0
)

go 1.17
