package main

import (
	"errors"
	"fmt"
	"github.com/NubeDev/bacnet"
	"github.com/NubeDev/bacnet/btypes"
	"github.com/NubeDev/bacnet/btypes/segmentation"
	"github.com/NubeDev/bacnet/network"
	"github.com/NubeIO/lib-networking/networking"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/uuid"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/utils/boolean"
	"github.com/NubeIO/rubix-os/utils/integer"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

type WhoIsOpts struct {
	WhoIs           bool   `json:"who_is"`
	InterfacePort   string `json:"interface_port"`
	LocalDeviceIP   string `json:"local_device_ip"`
	LocalDevicePort int    `json:"local_device_port"`
	LocalDeviceId   int    `json:"local_device_id"`
	Low             int    `json:"low"`
	High            int    `json:"high"`
	GlobalBroadcast bool   `json:"global_broadcast"`
	NetworkNumber   uint16 `json:"network_number"`
}

func (inst *Instance) masterWhoIs(opts *WhoIsOpts) (resp []*model.Device, err error) {
	// get network
	interfaces, err := getInterface(opts.InterfacePort)
	if err != nil {
		return nil, err
	}
	localDeviceIP := interfaces.IP
	if opts.LocalDeviceIP != "" {
		localDeviceIP = opts.LocalDeviceIP
	}
	localDevicePort := opts.LocalDevicePort
	if localDevicePort == 0 {
		localDevicePort = 47808
	}
	localDeviceId := opts.LocalDeviceId
	if localDeviceId == 0 {
		v := randInt(1, 100000)
		localDeviceId = v
	}
	interfaceName := interfaces.Interface
	log.Infof("run bacnet whoIs interface: %s ip: %s port: %d deviceID %d", interfaceName, localDeviceIP, localDevicePort, localDeviceId)
	localDevice, err := network.New(&network.Network{Interface: interfaceName, Port: localDevicePort, Store: inst.BacStore})
	if err != nil {
		log.Error(err)
		return
	}
	defer func(localDevice *network.Network, closeLogs bool) {
		err := localDevice.NetworkClose(closeLogs)
		if err != nil {
			log.Error(err)
		}
	}(localDevice, false)
	go localDevice.NetworkRun()
	var devices []btypes.Device
	if opts.WhoIs {
		devices, err = localDevice.Whois(&bacnet.WhoIsOpts{
			Low:             opts.Low,
			High:            opts.Low,
			GlobalBroadcast: true,
			NetworkNumber:   opts.NetworkNumber,
		})
		if err != nil {
			log.Error(err)
			return
		}
	} else {
		devices, err = localDevice.NetworkDiscover(&bacnet.WhoIsOpts{
			Low:             opts.Low,
			High:            opts.Low,
			GlobalBroadcast: true,
			NetworkNumber:   opts.NetworkNumber,
		})
		if err != nil {
			log.Error(err)
			return
		}
	}
	var devicesList []*model.Device
	for _, device := range devices {
		if device.DeviceName == "" {
			device.DeviceName = fmt.Sprintf("deviceId_%d_networkNum_%d", device.DeviceID, device.NetworkNumber)
		}
		newDevice := &model.Device{
			CommonUUID: model.CommonUUID{
				UUID: uuid.SmallUUID(),
			},
			CommonEnable: model.CommonEnable{
				Enable: boolean.NewTrue(),
			},
			Name: device.DeviceName,
			CommonDevice: model.CommonDevice{
				CommonIP: model.CommonIP{
					Host: device.Ip,
					Port: device.Port,
				},
				Manufacture: device.VendorName,
			},

			DeviceMac:      integer.New(device.MacMSTP),
			DeviceObjectId: integer.New(device.DeviceID),
			NetworkNumber:  integer.New(device.NetworkNumber),
			MaxADPU:        integer.New(int(device.MaxApdu)),
			Segmentation:   string(convertSegmentation(segmentation.SegmentedType(device.Segmentation))),
		}
		devicesList = append(devicesList, newDevice)
	}
	return devicesList, nil
}

var nets = networking.New()

func getInterface(networkInterface string) (networking.NetworkInterfaces, error) {
	if networkInterface == "" {
		interfaces, err := nets.GetInterfacesNames()
		if err != nil {
			return networking.NetworkInterfaces{}, err
		}
		for _, name := range interfaces.Names {
			if name != "lo" {
				iface, err := nets.GetNetworkByIface(name)
				if iface.IP != "" {
					return iface, err
				}
			}
		}
		return networking.NetworkInterfaces{}, errors.New("network interface can not be empty try, eth0")
	} else {
		return nets.GetNetworkByIface(networkInterface)
	}
}

// RandInt returns a random int within the specified range.
func randInt(range1, range2 int) int {
	if range1 == range2 {
		return range1
	}
	var min, max int
	if range1 > range2 {
		max = range1
		min = range2
	} else {
		max = range2
		min = range1
	}
	rand.Seed(time.Now().UnixNano())
	return min + rand.Intn(max-min)
}

func bodyMasterWhoIs(ctx *gin.Context) (dto *WhoIsOpts, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}
