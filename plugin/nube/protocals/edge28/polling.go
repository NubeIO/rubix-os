package main

import (
	"context"
	"github.com/NubeDev/flow-framework/api"
	"github.com/NubeDev/flow-framework/model"
	edgerest "github.com/NubeDev/flow-framework/plugin/nube/protocals/edge28/restclient"
	"github.com/NubeDev/flow-framework/src/poller"
	"github.com/NubeDev/flow-framework/utils"
	log "github.com/sirupsen/logrus"
	"time"
)

const defaultInterval = 2500 * time.Millisecond //default polling is 2.5 sec

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

//TODO add COV WriteValueOnceSync and InSync
func (i *Instance) processWrite(pnt *model.Point, value float64, rest *edgerest.RestClient, pollCount float64, isUO bool) (float64, error) {
	//!utils.BoolIsNil(pnt.WriteValueOnceSync)
	var err error
	if isUO {
		_, err = rest.WriteUO(pnt.IoID, value)
	} else {
		_, err = rest.WriteDO(pnt.IoID, value)
	}
	if err != nil {
		log.Errorf("edge-28: failed to write IO %s:  value:%f error:%v\n", pnt.IoID, value, err)
		return 0, err
	} else {
		log.Infof("edge-28: wrote IO %s: %v\n", pnt.IoID, value)
		return value, err
	}
}

//GPIOValueToDigital
//TODO remove this and get from helpers
func GPIOValueToDigital(value float64) float64 {
	if value < 0.2 {
		return 1 //ON / Closed Circuit
	} else { //previous functions used > 0.6 as an OFF threshold.
		return 0 //OFF / Open Circuit
	}
}

func (i *Instance) processRead(pnt *model.Point, value float64, pollCount float64) (float64, error) {
	cov := utils.Float64IsNil(pnt.COV) //TODO add in point scaling to get COV to work correct (as in scale temp or 0-10)
	covEvent, _ := utils.COV(value, utils.Float64IsNil(pnt.PresentValue), cov)
	if pollCount == 1 || !utils.BoolIsNil(pnt.InSync) {
		pnt.InSync = utils.NewTrue()
		_, err := i.db.UpdatePointValue(pnt.UUID, pnt, false)
		if err != nil {
			log.Errorf("edge-28: READ UPDATE POINT %s: %v\n", pnt.IoID, value)
			return value, err
		}
		if utils.BoolIsNil(pnt.InSync) {
			log.Infof("edge-28: READ POINT SYNC %s: %v\n", pnt.IoID, value)
		} else {
			log.Infof("edge-28: READ ON START %s: %v\n", pnt.IoID, value)
		}
	} else if covEvent {
		pnt.InSync = utils.NewTrue()
		_, err := i.db.UpdatePointValue(pnt.UUID, pnt, false)
		if err != nil {
			log.Errorf("edge-28: READ UPDATE POINT %s: %v\n", pnt.IoID, value)
			return value, err
		} else {
			log.Infof("edge-28: READ ON START %s: %v\n", pnt.IoID, value)
		}

	}
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
						if pnt.Priority != nil { //WRITE TODO add in all the extra GPIO
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
								_, err = i.processWrite(pnt, wv, rest, counter, false)
							case pointList.DO1, pointList.DO2, pointList.DO3, pointList.DO4, pointList.DO5:
								_, err = i.processWrite(pnt, wv, rest, counter, false)
							case pointList.UO1, pointList.UO2, pointList.UO3, pointList.UO4, pointList.UO5, pointList.UO6, pointList.UO7:
								_, err = i.processWrite(pnt, wv, rest, counter, true)
							}
						} //READ TODO add in all the extra GPIO
						_pnt.UUID = pnt.UUID
						switch pnt.IoID {
						case pointList.UI1:
							pv = getUI.Val.UI1.Val
							_, err = i.processRead(pnt, pv, counter)
						case pointList.UI2:
							pv = getUI.Val.UI2.Val
							_, err = i.processRead(pnt, pv, counter)
						case pointList.DI1:
							pv = GPIOValueToDigital(getDI.Val.DI1.Val)
							_, err = i.processRead(pnt, pv, counter)
						}
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
