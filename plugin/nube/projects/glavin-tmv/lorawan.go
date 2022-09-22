package main

import (
	"context"
	"crypto/tls"
	"flag"
	"log"

	"github.com/pkg/errors"
	"github.com/tealeg/xlsx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"

	"github.com/brocaar/lora-app-server/api"
)

// EXAMPLE OF CHIRPSTACK GO CLIENT: https://forum.chirpstack.io/t/creating-a-terraform-provider-for-loraserver/4081

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
	flag.StringVar(&apiHost, "api", "localhost:8080", "hostname:port to LoRa App Server API")
	flag.BoolVar(&apiInsecure, "api-insecure", false, "LoRa App Server API does not use TLS")
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

func getDeviceImportList() ([]DeviceImportRecord, error) {
	xlFile, err := xlsx.OpenFile(file)
	if err != nil {
		return nil, errors.Wrap(err, "open excel file error")
	}

	var out []DeviceImportRecord

	for _, sheet := range xlFile.Sheets {
		for i, row := range sheet.Rows {
			if i == 0 {
				continue
			}

			if len(row.Cells) != 7 {
				log.Fatalf("expected exactly 7 columns (row %d)", i+1)
			}

			devEUI := row.Cells[0].String()
			applicationID, err := row.Cells[1].Int64()
			if err != nil {
				log.Fatalf("application id parse error (row %d): %s", i+1, err)
			}
			deviceProfileID := row.Cells[2].String()
			name := row.Cells[3].String()
			description := row.Cells[4].String()
			networkKey := row.Cells[5].String()
			applicationKey := row.Cells[6].String()

			out = append(out, DeviceImportRecord{
				DevEUI:          devEUI,
				ApplicationID:   applicationID,
				DeviceProfileID: deviceProfileID,
				Name:            name,
				Description:     description,
				NetworkKey:      networkKey,
				ApplicationKey:  applicationKey,
			})
		}
	}

	return out, nil
}

func importDevices(conn *grpc.ClientConn, devices []DeviceImportRecord) error {
	deviceClient := api.NewDeviceServiceClient(conn)

	for i, dev := range devices {
		d := api.Device{
			DevEui:          dev.DevEUI,
			Name:            dev.Name,
			ApplicationId:   dev.ApplicationID,
			Description:     dev.Description,
			DeviceProfileId: dev.DeviceProfileID,
		}

		dk := api.DeviceKeys{
			DevEui: dev.DevEUI,
			NwkKey: dev.NetworkKey,
			AppKey: dev.ApplicationKey,
		}

		_, err := deviceClient.Create(context.Background(), &api.CreateDeviceRequest{
			Device: &d,
		})
		if err != nil {
			if grpc.Code(err) == codes.AlreadyExists {
				log.Printf("device %s already exists (row %d)", d.DevEui, i+2)
				continue
			}
			log.Fatalf("import error (device %s row %d): %s", d.DevEui, i+2, err)
		}

		_, err = deviceClient.CreateKeys(context.Background(), &api.CreateDeviceKeysRequest{
			DeviceKeys: &dk,
		})
		if err != nil {
			if grpc.Code(err) == codes.AlreadyExists {
				log.Printf("device-keys for device %s already exists (row %d)", d.DevEui, i+2)
				continue
			}
			log.Fatalf("import error (device %s) (row %d): %s", d.DevEui, i+2, err)
		}
	}

	return nil
}

func main() {
	conn, err := getGRPCConn()
	if err != nil {
		log.Fatal("error connecting to api", err)
	}

	if err := login(conn); err != nil {
		log.Fatal("login error", err)
	}

	rows, err := getDeviceImportList()
	if err != nil {
		log.Fatal("get device import records error", err)
	}

	if err := importDevices(conn, rows); err != nil {
		log.Fatal("import error", err)
	}
}
