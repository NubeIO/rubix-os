package main

import (
	"context"
	"fmt"
	"github.com/NubeDev/flow-framework/api"
	"github.com/NubeDev/flow-framework/model"
	edgerest "github.com/NubeDev/flow-framework/plugin/nube/protocals/edge28/restclient"
	"github.com/NubeDev/flow-framework/src/poller"
	"github.com/NubeDev/flow-framework/utils"
	log "github.com/sirupsen/logrus"
	"time"
)

const defaultInterval = 5000 * time.Millisecond

type polling struct {
	enable        bool
	loopDelay     time.Duration
	delayNetworks time.Duration
	delayDevices  time.Duration
	delayPoints   time.Duration
	isRunning     bool
}

var poll poller.Poller
var getUI *edgerest.UI
var getDI *edgerest.DI

func isWrite(t string) bool {
	switch t {
	case pointList.R1:
		return true
	case model.ObjectTypes.WriteHolding, model.ObjectTypes.WriteHoldings:
		return true
	case model.ObjectTypes.WriteSingleInt16, model.ObjectTypes.WriteSingleUint16:
		return true
	case model.ObjectTypes.WriteSingleFloat32, model.ObjectTypes.WriteSingleFloat64:
		return true
	}
	return false
}

func processWriteDO(pnt *model.Point, value float64, rest *edgerest.RestClient, pollCount float64) (float64, error) {
	//!utils.BoolIsNil(pnt.WriteValueOnceSync)
	_, err := rest.WriteDO(pnt.IoID, value)
	if err != nil {
		log.Errorf("edge-28: failed to write IO %s:  value:%f error:%v\n", pnt.IoID, value, err)
		return 0, err
	} else {
		log.Infof("edge-28: wrote IO %s: %v\n", pnt.IoID, value)
		return value, err
	}
}

func processWriteUO(pnt *model.Point, value float64, rest *edgerest.RestClient, pollCount float64) (float64, error) {
	//!utils.BoolIsNil(pnt.WriteValueOnceSync)
	_, err := rest.WriteUO(pnt.IoID, value)
	if err != nil {
		log.Errorf("edge-28: failed to write IO %s:  value:%f error:%v\n", pnt.IoID, value, err)
		return 0, err
	} else {
		log.Infof("edge-28: wrote IO %s: %v\n", pnt.IoID, value)
		return value, err
	}
}

func processRead(pnt *model.Point, value, pollCount float64) (float64, error) {
	cov := utils.Float64IsNil(pnt.COV)
	covEvent, _ := utils.COV(value, utils.Float64IsNil(pnt.ValueOriginal), cov)
	//write when pollCount is == 1
	//write on re-sync
	//write on covEvent

	if covEvent {
	}
	fmt.Println(getUI.Val.UI1.Val, pointList.UI1, pnt.UUID, pnt.IoID)
	//_, err := i.db.UpdatePointValue(pnt.UUID, &pnt, false)
	//log.Infof("modbus-write-cov: Addr: %s  Value: %f \n", pnt.IoID, pv)
	return value, nil

}

func (i *Instance) polling(p polling) error {
	if p.delayNetworks <= 0 {
		p.delayNetworks = defaultInterval
	}
	if p.delayDevices <= 0 {
		p.delayDevices = defaultInterval
	}
	if p.delayPoints <= 0 {
		p.delayPoints = defaultInterval
	}
	if p.enable {
		poll = poller.New()
	}
	var counter float64
	var arg api.Args
	arg.WithDevices = true
	arg.WithPoints = true
	arg.WithSerialConnection = true
	arg.WithIpConnection = true
	f := func() (bool, error) {
		nets, err := i.db.GetNetworksByPlugin(i.pluginUUID, arg)
		if err != nil {
			return false, err
		}
		if len(nets) == 0 {
			time.Sleep(15000 * time.Millisecond)
			log.Info("edge-28: NO NETWORKS FOUND")
		}
		for _, net := range nets { //NETWORKS
			if net.UUID != "" && net.PluginConfId == i.pluginUUID {
				log.Infof("edge-28: LOOP COUNT: %v\n", counter)
				counter++
				for _, dev := range net.Devices { //DEVICES
					if err != nil {
						log.Errorf("edge-28: failed to vaildate device %v %s\n", err, dev.CommonIP.Host)
					}
					rest := edgerest.NewNoAuth(dev.CommonIP.Host, dev.CommonIP.Port)
					getUI, err = rest.GetUIs()
					getDI, err = rest.GetDIs()
					dNet := p.delayNetworks
					time.Sleep(dNet)
					for _, pnt := range dev.Points { //POINTS
						var _pnt model.Point
						var pv float64
						if pnt.Priority != nil { //WRITE
							var wv float64
							if (*pnt.Priority).P16 != nil {
								wv = *pnt.Priority.P16
							}
							switch pnt.IoID {
							case pointList.R1, pointList.R2:
								if wv >= 1 {
									wv = 1
								} else {
									wv = 0
								}
								_, err = processWriteDO(pnt, wv, rest, counter)
							case pointList.DO1, pointList.DO2, pointList.DO3, pointList.DO4, pointList.DO5:
								_, err = processWriteDO(pnt, wv, rest, counter)
							case pointList.UO1, pointList.UO2, pointList.UO3, pointList.UO4, pointList.UO5, pointList.UO6, pointList.UO7:
								_, err = processWriteUO(pnt, wv, rest, counter)
							}
						} else { //READ
							_pnt.UUID = pnt.UUID
							switch pnt.IoID {
							case pointList.UI1:
								_, err = processRead(pnt, getUI.Val.UI1.Val, counter)
							case pointList.UI2:
								pv = getUI.Val.UI2.Val
								fmt.Println(getUI.Val.UI2.Val, pointList.UI2, pnt.UUID, pnt.IoID)
							case pointList.DI1:
								pv = getDI.Val.DI1.Val
								fmt.Println(getDI.Val.DI1.Val, pointList.DI1, pnt.UUID, pnt.IoID)
							}
						}

						_pnt.PresentValue = &pv
						_, err = i.db.UpdatePointValue(pnt.UUID, &_pnt, false)
						log.Infof("modbus-write-cov: Addr: %s  Value: %f \n", pnt.IoID, pv)
					}
				}
			}
		}
		if !p.enable { //TODO the disable of the polling isn't working
			return true, nil
		} else {
			return false, nil
		}

	}
	err := poll.Poll(context.Background(), f)
	if err != nil {
		return nil
	}
	return nil
}
