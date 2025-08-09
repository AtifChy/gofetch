// Package cpu provides functionality to collect CPU information.
package cpu

import (
	"context"
	"strings"

	"github.com/shirou/gopsutil/v4/cpu"

	"github.com/AtifChy/gofetch/internal/types"
)

func CollectCPUInfo(ctx context.Context) (*types.Info, error) {
	infos, err := cpu.InfoWithContext(ctx)
	if err != nil {
		return nil, err
	}
	return &types.Info{
		CPU: types.CPU{
			Model: strings.TrimSpace(infos[0].ModelName),
			Cores: infos[0].Cores,
		},
	}, nil
}
