module github.com/NubeIO/flow-framework

//replace github.com/NubeIO/nubeio-rubix-lib-helpers-go => /home/aidan/code/go/nube/nubeio-rubix-lib-helpers-go
replace github.com/NubeIO/nubeio-rubix-lib-models-go => /home/aidan/code/go/nube/lib/nubeio-rubix-lib-models-go
//replace github.com/NubeIO/nubeio-rubix-lib-rest-go => /home/aidan/code/go/nube/nubeio-rubix-lib-rest-go
replace github.com/NubeDev/bacnet => /home/aidan/code/go/nube/bacnet

require (
	github.com/NubeDev/location v0.0.2
	github.com/NubeIO/configor v0.0.3
	github.com/NubeIO/nubeio-rubix-lib-helpers-go v0.2.6
	github.com/NubeIO/nubeio-rubix-lib-models-go v1.2.2
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d
	github.com/eclipse/paho.mqtt.golang v1.3.5
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/gin v1.7.4
	github.com/go-co-op/gocron v1.7.0
	github.com/go-playground/validator/v10 v10.4.1
	github.com/go-resty/resty/v2 v2.7.0
	github.com/mustafaturan/bus/v3 v3.0.3
	github.com/mustafaturan/monoton/v2 v2.0.2
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.1
	gopkg.in/yaml.v2 v2.4.0
	gorm.io/datatypes v1.0.6
	gorm.io/driver/sqlite v1.3.1
	gorm.io/gorm v1.23.2
)

require go.bug.st/serial v1.3.2

require (
	github.com/NubeDev/modbus v0.0.2
	github.com/NubeIO/nubeio-rubix-lib-rest-go v1.0.6
	github.com/PaesslerAG/gval v1.1.1
	github.com/gomarkdown/markdown v0.0.0-20210918233619-6c1113f12c4a
	github.com/grid-x/modbus v0.0.0-20220210093200-c7b3bba92b40
	github.com/influxdata/influxdb-client-go/v2 v2.5.1
	github.com/labstack/gommon v0.3.0
	github.com/martinlindhe/unit v0.0.0-20210313160520-19b60e03648d
	github.com/minio/minio-go/v7 v7.0.14
	golang.org/x/sys v0.0.0-20220412211240-33da011f77ad
	gorm.io/driver/postgres v1.3.1
)

require github.com/NubeDev/bacnet v0.0.2

require (
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/creack/goselect v0.1.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/deepmap/oapi-codegen v1.8.2 // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.13.0 // indirect
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/go-yaml/yaml v2.1.0+incompatible // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.7 // indirect
	github.com/google/uuid v1.1.2 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/grid-x/serial v0.0.0-20191104121038-e24bc9bf6f08 // indirect
	github.com/influxdata/line-protocol v0.0.0-20200327222509-2487e7298839 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.11.0 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.2.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20200714003250-2b9c44734f2b // indirect
	github.com/jackc/pgtype v1.10.0 // indirect
	github.com/jackc/pgx/v4 v4.15.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.4 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid v1.3.1 // indirect
	github.com/kr/pretty v0.2.0 // indirect
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mattn/go-sqlite3 v1.14.12 // indirect
	github.com/minio/md5-simd v1.1.0 // indirect
	github.com/minio/sha256-simd v0.1.1 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	github.com/rs/xid v1.2.1 // indirect
	github.com/ugorji/go/codec v1.1.7 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.1 // indirect
	golang.org/x/crypto v0.0.0-20220411220226-7b82a4e95df4 // indirect
	golang.org/x/net v0.0.0-20220412020605-290c469a71a5 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/xerrors v0.0.0-20220411194840-2f41105eb62f // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/ini.v1 v1.66.4 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	gorm.io/driver/mysql v1.3.2 // indirect
)

go 1.17
