// Package collector provides functionality to collect system information
package collector

import (
	"context"
	"sync"

	"github.com/AtifChy/gofetch/internal/types"
)

type CollectorFunc func(context.Context) (*types.Info, error)

func Collect(ctx context.Context, collectors []CollectorFunc) (*types.Info, error) {
	infoChan := make(chan *types.Info, len(collectors))
	var wg sync.WaitGroup

	wg.Add(len(collectors))
	for _, collect := range collectors {
		go func(collector CollectorFunc) {
			defer wg.Done()
			info, err := collector(ctx)
			if err != nil {
				// Handle error appropriately
				return
			}
			infoChan <- info
		}(collect)
	}

	go func() {
		wg.Wait()
		close(infoChan)
	}()

	final := &types.Info{}
	for info := range infoChan {
		mergeInfo(final, info)
	}

	return final, nil
}

func mergeInfo(dst, src *types.Info) {
	if src.Host.Hostname != "" {
		dst.Host = src.Host
	}

	if len(src.Displays) > 0 {
		dst.Displays = append(dst.Displays, src.Displays...)
	}

	if src.CPU.Cores > 0 {
		dst.CPU = src.CPU
	}

	if len(src.GPUs) > 0 {
		dst.GPUs = append(dst.GPUs, src.GPUs...)
	}

	if src.Memory.Total > 0 {
		dst.Memory = src.Memory
	}

	if len(src.Disks) > 0 {
		dst.Disks = append(dst.Disks, src.Disks...)
	}
}
