package system

import (
	"fmt"
	systats "github.com/NubeIO/lib-system"
	"strings"
)

func (inst *System) GetSystem() (systats.System, error) {
	return inst.syStats.GetSystem()
}

func (inst *System) DiscUsage() ([]systats.Disk, error) {
	return inst.syStats.GetDisks()
}

type DiskUsage struct {
	Size      string `json:"size"`
	Used      string `json:"used"`
	Available string `json:"available"`
	Usage     string `json:"usage"`
}

type Disk struct {
	FileSystem string    `json:"file_system"`
	Type       string    `json:"type"`
	MountedOn  string    `json:"mounted_on"`
	Usage      DiskUsage `json:"usage"`
}

func (inst *System) DiscUsagePretty() ([]Disk, error) {
	var out []Disk
	disks, err := inst.syStats.GetDisks()
	if err != nil {
		return nil, err
	}
	for _, disk := range disks {
		newDisk := Disk{
			FileSystem: disk.FileSystem,
			Type:       disk.Type,
			MountedOn:  disk.MountedOn,
			Usage: DiskUsage{
				Size:      bytePretty(disk.Usage.Size),
				Used:      bytePretty(disk.Usage.Used),
				Available: bytePretty(disk.Usage.Available),
				Usage:     disk.Usage.Usage,
			},
		}
		out = append(out, newDisk)
	}
	return out, nil
}

type Memory struct {
	Stats     systats.Memory `json:"stats"`
	Available string
	Free      string
	Used      string
	Total     string
	Unit      string
}

type TopProcesses struct {
	Count int    `json:"count"`
	Sort  string `json:"sort"`
}

func (inst *System) GetTopProcesses(body TopProcesses) ([]systats.Process, error) {
	count := body.Count
	sort := body.Sort
	return inst.syStats.GetTopProcesses(count, sort)
}

type MemoryUsage struct {
	MemoryPercentageUsed float64
	MemoryPercentage     string
	MemoryAvailable      string
	MemoryFree           string
	MemoryUsed           string
	MemoryTotal          string
	SwapPercentageUsed   float64
	SwapPercentage       string
	SwapFree             string
	SwapUsed             string
	SwapTotal            string
}

func (inst *System) GetMemoryUsage() (*MemoryUsage, error) {

	m, err := inst.syStats.GetMemory(systats.Kilobyte)
	if err != nil {
		return nil, err
	}
	s, err := inst.syStats.GetSwap(systats.Kilobyte)
	if err != nil {
		return nil, err
	}

	return &MemoryUsage{
		MemoryPercentageUsed: m.PercentageUsed,
		MemoryPercentage:     fmt.Sprintf("%s", format(float32(m.PercentageUsed))) + "%",
		MemoryAvailable:      bytePretty(kbToByte(m.Available)),
		MemoryFree:           bytePretty(kbToByte(m.Free)),
		MemoryUsed:           bytePretty(kbToByte(m.Used)),
		MemoryTotal:          bytePretty(kbToByte(m.Total)),
		SwapPercentageUsed:   s.PercentageUsed,
		SwapPercentage:       fmt.Sprintf("%s", format(float32(s.PercentageUsed))) + "%",
		SwapFree:             bytePretty(kbToByte(s.Free)),
		SwapUsed:             bytePretty(kbToByte(s.Used)),
		SwapTotal:            bytePretty(kbToByte(s.Total)),
	}, nil
}

func (inst *System) GetMemory() (systats.Memory, error) {
	return inst.syStats.GetMemory(systats.Megabyte)
}

func (inst *System) GetSwap() (systats.Swap, error) {
	return inst.syStats.GetSwap(systats.Megabyte)
}

func kbToByte(input uint64) uint64 {
	return uint64(float64(input) * 1024)
}

func format(num float32) string {
	s := fmt.Sprintf("%.2f", num)
	return strings.TrimRight(strings.TrimRight(s, "0"), ".")
}

const (
	BYTE     = 1.0
	KILOBYTE = 1024 * BYTE
	MEGABYTE = 1024 * KILOBYTE
	GIGABYTE = 1024 * MEGABYTE
	TERABYTE = 1024 * GIGABYTE
)

func bytePretty(size uint64) string {
	unit := ""
	value := float32(size)

	switch {
	case size >= TERABYTE:
		unit = "T"
		value = value / TERABYTE
	case size >= GIGABYTE:
		unit = "G"
		value = value / GIGABYTE
	case size >= MEGABYTE:
		unit = "M"
		value = value / MEGABYTE
	case size >= KILOBYTE:
		unit = "K"
		value = value / KILOBYTE
	case size >= BYTE:
		unit = "B"
	case size == 0:
		return "0"
	}
	stringValue := fmt.Sprintf("%.2f", value)
	stringValue = strings.TrimSuffix(stringValue, ".00")
	return fmt.Sprintf("%s%s", stringValue, unit)
}
