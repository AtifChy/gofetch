// Package disk provides functionality to collect disk information
package disk

import (
	"context"

	"github.com/shirou/gopsutil/v4/disk"

	"github.com/AtifChy/gofetch/internal/types"
)

func CollectDiskInfo(ctx context.Context) (*types.Info, error) {
	parts, err := disk.PartitionsWithContext(ctx, true)
	if err != nil {
		return nil, err
	}
	var out []types.Disk
	for _, part := range parts {
		usage, err := disk.UsageWithContext(ctx, part.Mountpoint)
		if err != nil {
			return nil, err
		}
		out = append(out, types.Disk{
			FsType:      part.Fstype,
			Mountpoint:  part.Mountpoint,
			Total:       usage.Total,
			Free:        usage.Free,
			Used:        usage.Used,
			UsedPercent: uint8(usage.UsedPercent),
		})
	}
	return &types.Info{Disks: out}, nil
}
