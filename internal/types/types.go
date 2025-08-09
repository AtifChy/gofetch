// Package types
package types

import "time"

type Info struct {
	Host     Host
	Displays []Display
	CPU      CPU
	GPUs     []GPU
	Memory   Memory
	Disks    []Disk
}

type Host struct {
	Username        string
	Hostname        string
	OS              string
	Kernel          string
	PlatformVersion string
	Uptime          time.Duration
}

type CPU struct {
	Model string
	Cores int32
}

type GPU struct {
	Name string
	VRAM uint64
}

type Memory struct {
	Total       uint64
	Free        uint64
	Used        uint64
	UsedPercent uint8
}

type Display struct {
	Width       int32
	Height      int32
	RefreshRate int32
	IsPrimary   bool
}

type Disk struct {
	FsType      string
	Mountpoint  string
	Total       uint64
	Free        uint64
	Used        uint64
	UsedPercent uint8
}
