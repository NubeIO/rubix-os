package host

import (
	"github.com/NubeIO/flow-framework/src/system"
	"github.com/NubeIO/flow-framework/src/utilstime"
	"github.com/NubeIO/flow-framework/utils/array"
)

func GetCombinationData(debug bool) Combination {
	var comb Combination
	chServer := make(chan string)
	chTime := make(chan *utilstime.Time)
	chUptime := make(chan system.Details)
	chMem := make(chan *array.Array)
	chKernel := make(chan KernelInfo)
	chPro := make(chan ProgressInfo)
	chDisk := make(chan DiskInfoDetail)

	go func() { chServer <- getServerInfo(debug) }()
	go func() { chTime <- utilstime.SystemTime() }()
	go func() { chUptime <- system.Info() }()
	go func() { chMem <- GetMemory() }()
	go func() { chKernel <- getKernelData(debug) }()
	go func() { chPro <- getProgressData(debug) }()
	go func() { chDisk <- getDiskInfoDetail(debug) }()

	comb.ServerInfo = <-chServer
	comb.SystemTime = <-chTime
	comb.Uptime = <-chUptime
	comb.MemInfo = <-chMem
	comb.KernelInfo = <-chKernel
	comb.ProgressInfo = <-chPro
	comb.DiskInfo = <-chDisk

	return comb
}
