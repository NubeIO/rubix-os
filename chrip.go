package main

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/brocaar/lora-app-server/api"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
)

// JWTCredentials provides JWT credentials for gRPC
type JWTCredentials struct {
	token string
}

// GetRequestMetadata returns the meta-data for a request.
func (j *JWTCredentials) GetRequestMetadata(ctx context.Context, url ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": j.token,
	}, nil
}

// RequireTransportSecurity ...
func (j *JWTCredentials) RequireTransportSecurity() bool {
	return false
}

// SetToken sets the JWT token.
func (j *JWTCredentials) SetToken(token string) {
	j.token = token
}

// DeviceImportRecord defines a record for a device to import.
type DeviceImportRecord struct {
	DevEUI          string
	ApplicationID   int64
	DeviceProfileID string
	Name            string
	Description     string
	NetworkKey      string
	ApplicationKey  string
}

var (
	username       string
	password       string
	file           string
	apiHost        string
	apiInsecure    bool
	jwtCredentials *JWTCredentials
)

func init() {
	jwtCredentials = &JWTCredentials{}

	flag.StringVar(&username, "username", "admin", "LoRa App Server username")
	flag.StringVar(&password, "password", "admin", "LoRa App Server password")
	flag.StringVar(&file, "file", "", "Path to Excel file")
	flag.StringVar(&apiHost, "api", "123.209.246.84:8080", "hostname:port to LoRa App Server API")
	flag.BoolVar(&apiInsecure, "api-insecure", true, "LoRa App Server API does not use TLS")
	flag.Parse()
}

func getGRPCConn() (*grpc.ClientConn, error) {
	dialOpts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithPerRPCCredentials(jwtCredentials),
	}

	if apiInsecure {
		log.Println("using insecure api")
		dialOpts = append(dialOpts, grpc.WithInsecure())
	} else {
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: true,
		})))
	}

	conn, err := grpc.Dial(apiHost, dialOpts...)
	if err != nil {
		return nil, errors.Wrap(err, "grpc dial error")
	}

	return conn, nil
}

func login(conn *grpc.ClientConn) error {
	internalClient := api.NewInternalServiceClient(conn)
	fmt.Println("login")
	resp, err := internalClient.Login(context.Background(), &api.LoginRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		return errors.Wrap(err, "login error")
	}

	jwtCredentials.SetToken(resp.Jwt)

	return nil
}
func GenerateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func main() {

	fmt.Println(GenerateSecureToken(8))

	conn, err := getGRPCConn()
	if err != nil {
		log.Fatal("error connecting to api", err)
	}
	if err := login(conn); err != nil {
		log.Fatal("login error", err)
	}
	fmt.Println(conn.GetState())

	app := api.NewApplicationServiceClient(conn)

	//get, err := app.Get(context.Background())
	//if err != nil {
	//	return
	//}

	aa, err := app.Get(context.Background(), &api.GetApplicationRequest{
		//Device: &dev,
	})

	fmt.Println(aa.GetApplication(), 99999)

	deviceClient := api.NewDeviceServiceClient(conn)

	//add a device
	dev := api.Device{
		DevEui:          "c9014a013d89fa5d",
		Name:            "test name",
		ApplicationId:   2,
		Description:     "what up",
		DeviceProfileId: "b42d727b-1b10-46b8-972c-991046a4f952",
	}

	//activation
	dk := api.DeviceKeys{
		DevEui: "c9014a013d89fa5d",
		NwkKey: "01020304050607080807060504030201",
		AppKey: "00000000000000000000000000000000",
	}

	//get create a device
	_, err = deviceClient.Create(context.Background(), &api.CreateDeviceRequest{
		Device: &dev,
	})
	if err != nil {
		fmt.Println("ERROR: CreateDeviceRequest")
		fmt.Println(err)
	}
	_, err = deviceClient.CreateKeys(context.Background(), &api.CreateDeviceKeysRequest{
		DeviceKeys: &dk,
	})
	if err != nil {
		fmt.Println("ERROR: CreateDeviceKeysRequest")
		fmt.Println(err)
	}

	d := api.ListDeviceRequest{
		Limit:  10,
		Offset: 0,
	}
	//get all devices
	list, err := deviceClient.List(context.Background(), &d)
	if err != nil {
		fmt.Println("ERROR: get all devices")
		fmt.Println(err)
	}

	fmt.Println(list)

}
