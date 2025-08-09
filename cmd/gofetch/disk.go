package main

import (
	"context"

	"github.com/shirou/gopsutil/v4/disk"
)

type Disk struct {
	FsType      string
	Mountpoint  string
	Total       uint64
	Free        uint64
	Used        uint64
	UsedPercent uint8
}

func collectDiskInfo(ctx context.Context) (*Info, error) {
	parts, err := disk.PartitionsWithContext(ctx, true)
	if err != nil {
		return nil, err
	}
	var out []Disk
	for _, part := range parts {
		usage, err := disk.UsageWithContext(ctx, part.Mountpoint)
		if err != nil {
			return nil, err
		}
		out = append(out, Disk{
			FsType:      part.Fstype,
			Mountpoint:  part.Mountpoint,
			Total:       usage.Total,
			Free:        usage.Free,
			Used:        usage.Used,
			UsedPercent: uint8(usage.UsedPercent),
		})
	}
	return &Info{Disks: out}, nil
}
