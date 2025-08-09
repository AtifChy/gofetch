// Package memory provides functionality to collect memory information.
package memory

import (
	"context"

	"github.com/AtifChy/gofetch/internal/types"
	"github.com/shirou/gopsutil/v4/mem"
)

func CollectMemoryInfo(ctx context.Context) (*types.Info, error) {
	m, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		return nil, err
	}
	return &types.Info{
		Memory: types.Memory{
			Total:       m.Total,
			Used:        m.Used,
			Free:        m.Free,
			UsedPercent: uint8(m.UsedPercent),
		},
	}, nil
}
