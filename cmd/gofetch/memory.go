package main

import (
	"context"

	"github.com/shirou/gopsutil/v4/mem"
)

type Memory struct {
	Total       uint64
	Free        uint64
	Used        uint64
	UsedPercent uint8
}

func collectMemoryInfo(ctx context.Context) (*Info, error) {
	m, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		return nil, err
	}
	return &Info{
		Memory: Memory{
			Total:       m.Total,
			Used:        m.Used,
			Free:        m.Free,
			UsedPercent: uint8(m.UsedPercent),
		},
	}, nil
}
