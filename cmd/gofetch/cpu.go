package main

import (
	"context"
	"strings"

	"github.com/shirou/gopsutil/v4/cpu"
)

type CPU struct {
	Model string
	Cores int32
}

func collectCPUInfo(ctx context.Context) (*Info, error) {
	infos, err := cpu.InfoWithContext(ctx)
	if err != nil {
		return nil, err
	}
	return &Info{
		CPU: CPU{
			strings.TrimSpace(infos[0].ModelName),
			infos[0].Cores,
		},
	}, nil
}
