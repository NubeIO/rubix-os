package pollqueue

import (
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"math"
	"time"
)

func (pm *NetworkPollManager) GetPollingQueueStatistics() *model.PollQueueStatistics {
	pm.pollQueueDebugMsg("GetPollingQueueStatistics()")
	stats := model.PollQueueStatistics{}
	stats.Enable = pm.Enable
	stats.MaxPollRate = pm.MaxPollRate.String()

	stats.FFNetworkUUID = pm.FFNetworkUUID
	stats.NetworkName = pm.NetworkName
	stats.FFPluginUUID = pm.FFPluginUUID
	stats.PluginName = pm.PluginName

	MaxPollExecuteTime, _ := time.ParseDuration(fmt.Sprintf("%fs", pm.MaxPollExecuteTimeSecs))
	stats.MaxPollExecuteTime = MaxPollExecuteTime.String()
	AveragePollExecuteTime, _ := time.ParseDuration(fmt.Sprintf("%fs", pm.AveragePollExecuteTimeSecs))
	stats.AveragePollExecuteTime = AveragePollExecuteTime.String()
	MinPollExecuteTime, _ := time.ParseDuration(fmt.Sprintf("%fs", pm.MinPollExecuteTimeSecs))
	stats.MinPollExecuteTime = MinPollExecuteTime.String()
	stats.TotalPollQueueLength = pm.TotalPollQueueLength
	stats.TotalStandbyPointsLength = pm.TotalStandbyPointsLength
	stats.TotalPointsOutForPolling = pm.TotalPointsOutForPolling
	stats.ASAPPriorityPollQueueLength = pm.ASAPPriorityPollQueueLength
	stats.HighPriorityPollQueueLength = pm.HighPriorityPollQueueLength
	stats.NormalPriorityPollQueueLength = pm.NormalPriorityPollQueueLength
	stats.LowPriorityPollQueueLength = pm.LowPriorityPollQueueLength
	ASAPPriorityAveragePollTime, _ := time.ParseDuration(fmt.Sprintf("%fs", pm.ASAPPriorityAveragePollTime))
	stats.ASAPPriorityAveragePollTime = ASAPPriorityAveragePollTime.String()
	HighPriorityAveragePollTime, _ := time.ParseDuration(fmt.Sprintf("%fs", pm.HighPriorityAveragePollTime))
	stats.HighPriorityAveragePollTime = HighPriorityAveragePollTime.String()
	NormalPriorityAveragePollTime, _ := time.ParseDuration(fmt.Sprintf("%fs", pm.NormalPriorityAveragePollTime))
	stats.NormalPriorityAveragePollTime = NormalPriorityAveragePollTime.String()
	LowPriorityAveragePollTime, _ := time.ParseDuration(fmt.Sprintf("%fs", pm.LowPriorityAveragePollTime))
	stats.LowPriorityAveragePollTime = LowPriorityAveragePollTime.String()
	stats.TotalPollCount = pm.TotalPollCount
	stats.ASAPPriorityPollCount = pm.ASAPPriorityPollCount
	stats.HighPriorityPollCount = pm.HighPriorityPollCount
	stats.NormalPriorityPollCount = pm.NormalPriorityPollCount
	stats.LowPriorityPollCount = pm.LowPriorityPollCount
	stats.ASAPPriorityMaxCycleTime = pm.ASAPPriorityMaxCycleTime.String()
	stats.HighPriorityMaxCycleTime = pm.HighPriorityMaxCycleTime.String()
	stats.NormalPriorityMaxCycleTime = pm.NormalPriorityMaxCycleTime.String()
	stats.LowPriorityMaxCycleTime = pm.LowPriorityMaxCycleTime.String()
	stats.ASAPPriorityLockupAlert = pm.ASAPPriorityLockupAlert
	stats.HighPriorityLockupAlert = pm.HighPriorityLockupAlert
	stats.NormalPriorityLockupAlert = pm.NormalPriorityLockupAlert
	stats.LowPriorityLockupAlert = pm.LowPriorityLockupAlert
	stats.BusyTime = fmt.Sprintf("%.1f%%", pm.BusyTime)
	EnabledTime, _ := time.ParseDuration(fmt.Sprintf("%fs", pm.EnabledTime))
	stats.EnabledTime = EnabledTime.String()
	PortUnavailableTime, _ := time.ParseDuration(fmt.Sprintf("%fs", pm.PortUnavailableTime))
	stats.PortUnavailableTime = PortUnavailableTime.String()

	return &stats
}

func (pm *NetworkPollManager) StartPollingStatistics() {
	pm.pollQueueDebugMsg("StartPollingStatistics()")
	pm.PollingStartTimeUnix = time.Now().Unix()
	pm.AveragePollExecuteTimeSecs = 0
	pm.MaxPollExecuteTimeSecs = 0
	pm.MinPollExecuteTimeSecs = 0
	pm.ASAPPriorityAveragePollTime = 0
	pm.HighPriorityAveragePollTime = 0
	pm.NormalPriorityAveragePollTime = 0
	pm.LowPriorityAveragePollTime = 0
	pm.TotalPollCount = 0
	pm.ASAPPriorityPollCount = 0
	pm.HighPriorityPollCount = 0
	pm.NormalPriorityPollCount = 0
	pm.LowPriorityPollCount = 0
	pm.ASAPPriorityPollCountForAvg = 0
	pm.HighPriorityPollCountForAvg = 0
	pm.NormalPriorityPollCountForAvg = 0
	pm.LowPriorityPollCountForAvg = 0
	pm.ASAPPriorityLockupAlert = false
	pm.HighPriorityLockupAlert = false
	pm.NormalPriorityLockupAlert = false
	pm.LowPriorityLockupAlert = false
	pm.PortUnavailableTime = 0
	pm.PortUnavailableStartTime = 0
}

func (pm *NetworkPollManager) PollCompleteStatsUpdate(pp *PollingPoint, pollTimeSecs float64) {
	pm.pollQueueDebugMsg("PollCompleteStatsUpdate()")

	if pm.MaxPollExecuteTimeSecs == 0 || pollTimeSecs > pm.MaxPollExecuteTimeSecs {
		pm.MaxPollExecuteTimeSecs = pollTimeSecs
	}
	if pm.MinPollExecuteTimeSecs == 0 || pollTimeSecs < pm.MinPollExecuteTimeSecs {
		pm.MinPollExecuteTimeSecs = pollTimeSecs
	}
	pm.AveragePollExecuteTimeSecs = ((pm.AveragePollExecuteTimeSecs * float64(pm.TotalPollCount)) + pollTimeSecs) / (float64(pm.TotalPollCount) + 1)
	pm.TotalPollCount++
	pm.EnabledTime = time.Since(time.Unix(pm.PollingStartTimeUnix, 0)).Seconds()
	pm.BusyTime = math.Round((((pm.AveragePollExecuteTimeSecs*float64(pm.TotalPollCount))/pm.EnabledTime)*100)*1000) / 1000 // percentage rounded to 3 decimal places

	pm.TotalPollQueueLength = int64(pm.PollQueue.PriorityQueue.Len())
	pm.TotalStandbyPointsLength = int64(pm.PollQueue.StandbyPollingPoints.Len())
	pm.TotalPointsOutForPolling = int64(pm.PollQueue.OutstandingPollingPoints.Len())

	pm.ASAPPriorityPollQueueLength = 0
	pm.HighPriorityPollQueueLength = 0
	pm.NormalPriorityPollQueueLength = 0
	pm.LowPriorityPollQueueLength = 0

	for _, pp := range pm.PollQueue.PriorityQueue.PriorityQueue {
		if pp != nil {
			switch pp.PollPriority {
			case model.PRIORITY_ASAP:
				pm.ASAPPriorityPollQueueLength++
			case model.PRIORITY_HIGH:
				pm.HighPriorityPollQueueLength++
			case model.PRIORITY_NORMAL:
				pm.NormalPriorityPollQueueLength++
			case model.PRIORITY_LOW:
				pm.LowPriorityPollQueueLength++
			}
		}
	}
	pm.TotalPollQueueLength = int64(pm.PollQueue.PriorityQueue.Len())

	switch pp.PollPriority {
	case model.PRIORITY_ASAP:
		pm.ASAPPriorityPollCount++
		if pp.QueueEntryTime <= 0 {
			return
		}
		pollTime := float64(time.Now().Unix() - pp.QueueEntryTime)
		pm.ASAPPriorityAveragePollTime = ((pm.ASAPPriorityAveragePollTime * float64(pm.ASAPPriorityPollCountForAvg)) + pollTime) / (float64(pm.ASAPPriorityPollCountForAvg) + 1)
		pm.ASAPPriorityPollCountForAvg++

	case model.PRIORITY_HIGH:
		pm.HighPriorityPollCount++
		if pp.QueueEntryTime <= 0 {
			return
		}
		pollTime := float64(time.Now().Unix() - pp.QueueEntryTime)
		pm.HighPriorityAveragePollTime = ((pm.HighPriorityAveragePollTime * float64(pm.HighPriorityPollCountForAvg)) + pollTime) / (float64(pm.HighPriorityPollCountForAvg) + 1)
		pm.HighPriorityPollCountForAvg++

	case model.PRIORITY_NORMAL:
		pm.NormalPriorityPollCount++
		if pp.QueueEntryTime <= 0 {
			return
		}
		pollTime := float64(time.Now().Unix() - pp.QueueEntryTime)
		pm.NormalPriorityAveragePollTime = ((pm.NormalPriorityAveragePollTime * float64(pm.NormalPriorityPollCountForAvg)) + pollTime) / (float64(pm.NormalPriorityPollCountForAvg) + 1)
		pm.NormalPriorityPollCountForAvg++

	case model.PRIORITY_LOW:
		pm.LowPriorityPollCount++
		if pp.QueueEntryTime <= 0 {
			return
		}
		pollTime := float64(time.Now().Unix() - pp.QueueEntryTime)
		pm.LowPriorityAveragePollTime = ((pm.LowPriorityAveragePollTime * float64(pm.LowPriorityPollCountForAvg)) + pollTime) / (float64(pm.LowPriorityPollCountForAvg) + 1)
		pm.LowPriorityPollCountForAvg++

	}

}

func (pm *NetworkPollManager) PartialPollStatsUpdate() {
	pm.pollQueueDebugMsg("PartialPollStatsUpdate()")
	pm.TotalPollQueueLength = int64(pm.PollQueue.PriorityQueue.Len())
	pm.TotalStandbyPointsLength = int64(pm.PollQueue.StandbyPollingPoints.Len())
	pm.TotalPointsOutForPolling = int64(pm.PollQueue.OutstandingPollingPoints.Len())

	pm.EnabledTime = time.Since(time.Unix(pm.PollingStartTimeUnix, 0)).Seconds()

	if pm.PortUnavailableTimeout != nil {
		pm.PortUnavailableTime += time.Since(time.Unix(pm.PortUnavailableStartTime, 0)).Seconds()
		pm.PortUnavailableStartTime = time.Now().Unix()
	}

	pm.ASAPPriorityPollQueueLength = 0
	pm.HighPriorityPollQueueLength = 0
	pm.NormalPriorityPollQueueLength = 0
	pm.LowPriorityPollQueueLength = 0

	for _, pp := range pm.PollQueue.PriorityQueue.PriorityQueue {
		if pp != nil {
			switch pp.PollPriority {
			case model.PRIORITY_ASAP:
				pm.ASAPPriorityPollQueueLength++
			case model.PRIORITY_HIGH:
				pm.HighPriorityPollQueueLength++
			case model.PRIORITY_NORMAL:
				pm.NormalPriorityPollQueueLength++
			case model.PRIORITY_LOW:
				pm.LowPriorityPollQueueLength++
			}
		}
	}
	pm.TotalPollQueueLength = int64(pm.PollQueue.PriorityQueue.Len())
}
