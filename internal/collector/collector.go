// Package collector provides functionality to collect system information
package collector

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/AtifChy/gofetch/internal/types"
)

type CollectorFunc func(context.Context) (*types.Info, error)

func Collect(ctx context.Context, collectors []CollectorFunc) (*types.Info, error) {
	type result struct {
		info *types.Info
		err  error
	}

	resultChan := make(chan result, len(collectors))
	var wg sync.WaitGroup

	wg.Add(len(collectors))
	for _, collect := range collectors {
		go func(collector CollectorFunc) {
			defer wg.Done()
			info, err := collector(ctx)
			resultChan <- result{info: info, err: err}
		}(collect)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	final := &types.Info{}
	successCount := 0

	for result := range resultChan {
		if result.err != nil {
			log.Printf("Error collecting info: %s", result.err)
			continue
		}
		merge(final, result.info)
		successCount++
	}

	if successCount == 0 {
		return nil, fmt.Errorf("all collectors failed to retrieve information")
	}

	return final, nil
}

func merge(dst, src *types.Info) {
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
