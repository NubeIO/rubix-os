package host

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-framework/utils"
	"os/exec"
	"strconv"
	"strings"
)

func ByteCountSI(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}

const (
	memTotal     = "MemTotal"
	memFree      = "MemFree"
	memAvailable = "MemAvailable"
	memInfo      = "/proc/meminfo"
)

func getMem(memType string) (out string, err error) {
	switch memType {
	case memTotal:
	case memFree:
	case memAvailable:
	default:
		return "", errors.New("invalid memory type try: MemTotal, MemFree or MemAvailable")
	}
	cmd := exec.Command("grep", memType, memInfo)
	c, _ := cmd.CombinedOutput()
	strOut := string(c)
	mt := fmt.Sprintf("%s:", memType)
	strOut = strings.ReplaceAll(strOut, mt, "")
	strOut = strings.ReplaceAll(strOut, "kB", "")
	strOut = strings.ReplaceAll(strOut, " ", "")
	strOut = strings.ReplaceAll(strOut, "\n", "")
	ramSize, _ := strconv.ParseInt(strOut, 10, 64)
	ramSizeInt := ramSize * 1000
	return ByteCountSI(ramSizeInt), nil
}

type Memory struct {
	MemoryTotal     string `json:"memory_total"`
	MemoryFree      string `json:"memory_free"`
	MemoryAvailable string `json:"memory_available"`
}

func GetMemory() *utils.Array {
	var mem Memory
	memT, err := getMem(memTotal)
	if err != nil {
		return nil
	}
	memF, err := getMem(memFree)
	if err != nil {
		return nil
	}
	memA, err := getMem(memAvailable)
	if err != nil {
		return nil
	}
	mem.MemoryTotal = memT
	mem.MemoryFree = memF
	mem.MemoryAvailable = memA
	out := utils.NewArray()
	out.Add(mem)
	return out

}
