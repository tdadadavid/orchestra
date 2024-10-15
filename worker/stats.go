package worker

import (
	"log"

	"github.com/c9s/goprocinfo/linux"
)


type Stats struct {
	MemStats  *linux.MemInfo
	DiskStats *linux.Disk
	CpuStats  *linux.CPUStat
	LoadStats *linux.LoadAvg
	TaskCount int
}


//////
// Memory Related Metrics of Workers
//////

//TotalMemKb returns the total memory.
func (s *Stats) TotalMemKb() uint64 {
	return s.MemStats.MemTotal
}

//TotalMemKb returns the available memory.
func (s *Stats) AvailableMemKb() uint64 {
	return s.MemStats.MemAvailable
}

//UsedMemKb returns the used memory.
func (s *Stats) UsedMemKb() uint64 {
	return s.MemStats.MemTotal - s.MemStats.MemAvailable
}

//UsedMemPercent returns the used memory in percentage 
func (s *Stats) UsedMemPercent() uint64 {
	return s.MemStats.MemAvailable / s.MemStats.MemTotal
}


//////
// Disk Related Metrics of Workers
//////

//TotalDisk returns the total disk size 
func (s *Stats) TotalDisk() uint64 {
	return s.DiskStats.All
}

//FreeSpaceInDisk returns the free space in disk. 
func (s *Stats) FreeSpaceInDisk() uint64 {
	return s.DiskStats.Free
}

//UsedDisk returns the used space in disk. 
func (s *Stats) UsedDisk() uint64 {
	return s.DiskStats.Used
}

//////
// CPU Related Metrics of Workers
//////

//UsedDisk returns the used space in disk. 
//ReadMore: https://claude.ai/chat/1eb22b8b-f8f1-48dd-b676-1df9e227c867
func (s *Stats) CpuUsage() float64 {
	// sum the idle states
	idle := s.CpuStats.Idle + s.CpuStats.IOWait

	// sum the non-idle states
	nonIdle := s.CpuStats.Guest + s.CpuStats.IRQ + s.CpuStats.System + s.CpuStats.SoftIRQ + s.CpuStats.User + s.CpuStats.Nice + s.CpuStats.GuestNice;

	total := idle + nonIdle

	if total == 0 {
		return 0.00;
	}

	return (float64(total) - float64(total)) / float64(total)
}


func GetStats() *Stats {
	return &Stats{
		MemStats: GetMemoryInfo(),
		DiskStats: GetDiskInfo(),
		CpuStats: GetCpuInfo(),
		LoadStats: GetLoadInfo(),
	}
}

func GetMemoryInfo() *linux.MemInfo {
	memstats, err :=  linux.ReadMemInfo("/proc/meminfo")
	if err != nil {
		log.Println("Error reading from linux `/proc/meminfo` file")
		return &linux.MemInfo{}
	}

	return memstats
}

func GetDiskInfo() *linux.Disk {
	diskstats, err := linux.ReadDisk("/proc/diskstats")
	if err != nil {
		log.Println("Error reading from linux `/proc/diskstats` file")
		return &linux.Disk{}
	}

	return diskstats
}

func GetCpuInfo() *linux.CPUStat {
	cpuStats, err := linux.ReadStat("/proc/stat")
	if err != nil {
		log.Println("Error reading from linux `/proc/stat` file")
		return &linux.CPUStat{}
	}

	return &cpuStats.CPUStatAll
}

func GetLoadInfo() *linux.LoadAvg {
	loadAvg, err := linux.ReadLoadAvg("/proc/loadavg")
	if err != nil {
		log.Println("Error reading from linux `/proc/loadavg` file")
		return &linux.LoadAvg{}
	}

	return loadAvg
}
