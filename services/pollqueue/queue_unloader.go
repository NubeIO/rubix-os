package pollqueue

import (
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/utils/float"
	"time"
	// log "github.com/sirupsen/logrus"
)

// REFS:
//  - GOLANG HEAP https://pkg.go.dev/container/heap
//  - Worker Queue tutorial: https://www.opsdash.com/blog/job-queues-in-go.html

type QueueUnloader struct {
	NextPollPoint *PollingPoint
	// NextUnloadTimer 		*time.Timer
	NextUnloadTimer *time.Ticker
	CancelChannel   chan bool
}

func (pm *NetworkPollManager) StartQueueUnloader() {
	pm.StopQueueUnloader()
	pm.pollQueueDebugMsg("StartQueueUnloader()")
	ql := &QueueUnloader{nil, nil, nil}
	pm.PluginQueueUnloader = ql
	if pm.PluginQueueUnloader.NextPollPoint == nil {
		pm.pollQueueDebugMsg("StartQueueUnloader() pm.PluginQueueUnloader.NextPollPoint == nil")
		pp, err := pm.PollQueue.GetNextPollingPoint()
		if pp != nil && err == nil {
			pm.PluginQueueUnloader.NextPollPoint = pp
		}
	}

	refreshRate := 100 * time.Millisecond // Default MaxPollRate
	if pm.DBHandlerRef != nil {
		var netArg api.Args
		net, err := pm.DBHandlerRef.GetNetwork(pm.FFNetworkUUID, netArg)
		if err != nil {
			pm.pollQueueDebugMsg(fmt.Sprintf("NetworkPollManager.StartQueueUnloader(): couldn't find network %s", pm.FFNetworkUUID))
			return
		}
		if float.NonNil(net.MaxPollRate) > 0 {
			refreshRate, _ = time.ParseDuration(fmt.Sprintf("%fs", *net.MaxPollRate))
			pm.pollQueueDebugMsg(fmt.Sprintf("NetworkPollManager.StartQueueUnloader(): net.MaxPollRate %d ", net.MaxPollRate))
		}
	} else {
		pm.pollQueueErrorMsg("StartQueueUnloader(): NetworkPollManager DBHandlerRef is undefined")
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
				// pm.pollQueueDebugMsg("RELOAD QUEUE TICKER")
				pm.postNextPointCallback()
			}
		}
	}()
}

func (pm *NetworkPollManager) StopQueueUnloader() {
	pm.pollQueueDebugMsg("StopQueueUnloader()")
	if pm.PluginQueueUnloader != nil && pm.PluginQueueUnloader.NextUnloadTimer != nil && pm.PluginQueueUnloader.CancelChannel != nil {
		pm.PluginQueueUnloader.NextUnloadTimer.Stop()
		pm.PluginQueueUnloader.NextUnloadTimer = nil
		pm.PluginQueueUnloader.CancelChannel <- true
		pm.PluginQueueUnloader.CancelChannel = nil
		pm.pollQueueDebugMsg("StopQueueUnloader() NextUnloadTimer stopped and CancelChannel closed")
	}
	pm.PluginQueueUnloader = nil
}

// This function should be called from the Polling service.
func (pm *NetworkPollManager) GetNextPollingPoint() (pp *PollingPoint, callback func(pp *PollingPoint, writeSuccess, readSuccess bool, pollTimeSecs float64, pointUpdate, resetToConfiguredPriority bool, retryType PollRetryType)) {
	if pm.PluginQueueUnloader != nil && pm.PluginQueueUnloader.NextPollPoint != nil {
		pp := pm.PluginQueueUnloader.NextPollPoint
		pm.PluginQueueUnloader.NextPollPoint = nil
		// Moving the line below to a reoccurring timer instead.
		// pm.PluginQueueUnloader.NextUnloadTimer = time.AfterFunc(pm.MaxPollRate, pm.postNextPointCallback)
		return pp, pm.PollingPointCompleteNotification
	}
	// pm.pollQueueDebugMsg("GetNextPollingPoint(): No pollingPoint available")
	return nil, nil
}

// This is the callback function that is called by the reoccurring timer (seperate go routine) made in StartQueueUnloader().
func (pm *NetworkPollManager) postNextPointCallback() {
	if pm.PluginQueueUnloader != nil && pm.PluginQueueUnloader.NextPollPoint == nil {
		pp, err := pm.PollQueue.GetNextPollingPoint()
		if pp != nil && err == nil {
			pm.PluginQueueUnloader.NextPollPoint = pp
			addSuccess := pm.PollQueue.OutstandingPollingPoints.AddPollingPoint(pp)
			if !addSuccess {
				pm.pollQueueErrorMsg(fmt.Sprintf("Modbus postNextPointCallback(): polling point could not be added to OutstandingPollingPoints slice.  (%s)", pp.FFPointUUID))
			}
		}
	}
}
