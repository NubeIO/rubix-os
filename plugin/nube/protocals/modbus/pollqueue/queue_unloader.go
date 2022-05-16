package pollqueue

import (
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"time"
	//log "github.com/sirupsen/logrus"
)

// REFS:
//  - GOLANG HEAP https://pkg.go.dev/container/heap
//  - Worker Queue tutorial: https://www.opsdash.com/blog/job-queues-in-go.html

type QueueUnloader struct {
	NextPollPoint *PollingPoint
	//NextUnloadTimer 		*time.Timer
	NextUnloadTimer *time.Ticker
	CancelChannel   chan bool
}

func (pm *NetworkPollManager) StartQueueUnloader() {
	pollQueueDebugMsg("StartQueueUnloader() 1")
	pm.StopQueueUnloader()
	pollQueueDebugMsg("StartQueueUnloader() 2")
	ql := &QueueUnloader{nil, nil, nil}
	pm.PluginQueueUnloader = ql
	if pm.PluginQueueUnloader.NextPollPoint == nil {
		pollQueueDebugMsg("StartQueueUnloader() pm.PluginQueueUnloader.NextPollPoint == nil")
		pp, err := pm.PollQueue.GetNextPollingPoint()
		if pp != nil && err == nil {
			pm.PluginQueueUnloader.NextPollPoint = pp
		}
	}
	var netArg api.Args
	net, err := pm.DBHandlerRef.GetNetwork(pm.FFNetworkUUID, netArg)
	if err != nil {
		pollQueueDebugMsg(fmt.Sprintf("NetworkPollManager.StartQueueUnloader(): couldn't find network %s", pm.FFNetworkUUID))
		return
	}
	pollQueueDebugMsg(fmt.Sprintf("NetworkPollManager.StartQueueUnloader(): net.MaxPollRate %d ", net.MaxPollRate))
	refreshRate := 100 * time.Millisecond
	if *net.MaxPollRate > 0 {
		refreshRate, _ = time.ParseDuration(fmt.Sprintf("%fs", *net.MaxPollRate))
	}
	pm.MaxPollRate = refreshRate
	ticker := time.NewTicker(refreshRate)
	pm.PluginQueueUnloader.NextUnloadTimer = ticker
	done := make(chan bool)
	pm.PluginQueueUnloader.CancelChannel = done

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				pollQueueDebugMsg("RELOAD QUEUE TICKER")
				pm.postNextPointCallback()
			}
		}
	}()
}

func (pm *NetworkPollManager) StopQueueUnloader() {
	pollQueueDebugMsg("StopQueueUnloader()")
	if pm.PluginQueueUnloader != nil && pm.PluginQueueUnloader.NextUnloadTimer != nil && pm.PluginQueueUnloader.CancelChannel != nil {
		pm.PluginQueueUnloader.NextUnloadTimer.Stop()
		pm.PluginQueueUnloader.CancelChannel <- true
		pollQueueDebugMsg("StopQueueUnloader() NextUnloadTimer stopped and CancelChannel closed")
	}
	pm.PluginQueueUnloader = nil
}

//This function should be called from the Polling service. It will start a timer that posts the next polling point.
func (pm *NetworkPollManager) GetNextPollingPoint() (pp *PollingPoint, callback func(pp *PollingPoint, writeSuccess, readSuccess bool, pollTimeSecs float64, pointUpdate bool)) {
	pollQueueDebugMsg("GetNextPollingPoint()")
	if pm.PluginQueueUnloader != nil && pm.PluginQueueUnloader.NextPollPoint != nil {
		pp := pm.PluginQueueUnloader.NextPollPoint
		pm.PluginQueueUnloader.NextPollPoint = nil
		//Moving the line below to a reoccurring timer instead.
		//pm.PluginQueueUnloader.NextUnloadTimer = time.AfterFunc(pm.MaxPollRate, pm.postNextPointCallback)
		return pp, pm.PollingPointCompleteNotification
	}
	pollQueueDebugMsg("GetNextPollingPoint(): No pollingPoint available")
	return nil, nil
}

//This is the callback function that is called by the timer made in (pm *NetworkPollManager) GetNextPollingPoint().
func (pm *NetworkPollManager) postNextPointCallback() {
	pollQueueDebugMsg("postNextPointCallback()")
	if pm.PluginQueueUnloader != nil && pm.PluginQueueUnloader.NextPollPoint == nil {
		pp, err := pm.PollQueue.GetNextPollingPoint()
		if pp != nil && err == nil {
			pm.PluginQueueUnloader.NextPollPoint = pp
		}
	}
}
