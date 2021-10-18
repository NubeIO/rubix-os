package host

import (
	"github.com/NubeDev/flow-framework/src/system"
	"github.com/NubeDev/flow-framework/src/utilstime"
	"github.com/NubeDev/flow-framework/utils"
)

func GetCombinationData(debug bool) Combination {
	var comb Combination
	chServer := make(chan string)
	chTime := make(chan *utilstime.Time)
	chUptime := make(chan system.Details)
	chCPU := make(chan CPUInfo)
	chMem := make(chan *utils.Array)
	chKernel := make(chan KernelInfo)
	chPro := make(chan ProgressInfo)
	chDisk := make(chan DiskInfoDetail)

	go func() { chServer <- getServerInfo(debug) }()
	go func() { chTime <- utilstime.SystemTime() }()
	go func() { chUptime <- system.Info() }()
	go func() { chCPU <- getCPUInfo(debug) }()
	go func() { chMem <- GetMemory() }()
	go func() { chKernel <- getKernelData(debug) }()
	go func() { chPro <- getProgressData(debug) }()
	go func() { chDisk <- getDiskInfoDetail(debug) }()

	comb.ServerInfo = <-chServer
	comb.SystemTime = <-chTime
	comb.Uptime = <-chUptime
	comb.CPUInfo = <-chCPU
	comb.MemInfo = <-chMem
	comb.KernelInfo = <-chKernel
	comb.ProgressInfo = <-chPro
	comb.DiskInfo = <-chDisk

	return comb
}
