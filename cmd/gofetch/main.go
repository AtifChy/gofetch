package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/hako/durafmt"
)

type Info struct {
	OS              string
	Hostname        string
	KernelVersion   string
	PlatformVersion string
	Displays        []Display
	CPU             CPU
	GPUs            []GPU
	Memory          Memory
	Disks           []Disk
	Uptime          time.Duration
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Channel to receive system information
	infoChan := make(chan *Info)
	// WaitGroup to wait for goroutines to finish
	var wg sync.WaitGroup

	// List of collector functions
	collectors := []func(context.Context) (*Info, error){
		collectHostInfo,
		collectCPUInfo,
		collectDisplayInfo,
		collectGPUInfo,
		collectMemoryInfo,
		collectDiskInfo,
	}

	// Start collectors concurrently
	wg.Add(len(collectors))
	for _, collect := range collectors {
		go func(collector func(context.Context) (*Info, error)) {
			defer wg.Done()
			info, err := collector(ctx)
			if err != nil {
				log.Fatalf("Error collecting info: %s", err)
				return
			}
			infoChan <- info
		}(collect)
	}

	// Go routine to close the channel once all collectors are done
	go func() {
		wg.Wait()
		close(infoChan)
	}()

	// Aggregate results from the channel
	final := &Info{}
	for info := range infoChan {
		mergeInfo(final, info)
	}

	// Display final info
	displayInfo(final)
}

func mergeInfo(dst, src *Info) {
	if src.OS != "" {
		dst.OS = src.OS
		dst.KernelVersion = src.KernelVersion
		dst.PlatformVersion = src.PlatformVersion
		dst.Uptime = src.Uptime
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

func displayInfo(info *Info) {
	title := getTitle()
	fmt.Printf("%s\n", title)

	fmt.Println(strings.Repeat("-", len(title)))

	fmt.Printf("OS: %s\n", info.OS)
	fmt.Printf("Kernel: %s %s (%s)\n", getKernelName(), info.KernelVersion, info.PlatformVersion)
	fmt.Printf("Uptime: %s\n", durafmt.Parse(info.Uptime).LimitFirstN(2).String())

	displayStr := "Display: "
	for i, display := range info.Displays {
		if display.IsPrimary {
			displayStr += "*"
		}
		displayStr += fmt.Sprintf("%dx%d @ %dHz", display.Width, display.Height, display.RefreshRate)
		if i < len(info.Displays)-1 {
			displayStr += ", "
		}
	}
	fmt.Println(displayStr)

	fmt.Printf("CPU: %s (%d)\n", info.CPU.Model, info.CPU.Cores)

	for i, gpu := range info.GPUs {
		fmt.Printf("GPU %d: %s (%s)\n", i+1, gpu.Name, humanize.IBytes(gpu.VRAM))
	}

	fmt.Printf("Memory: %s / %s (%d%%)\n",
		humanize.IBytes(info.Memory.Used),
		humanize.IBytes(info.Memory.Total),
		info.Memory.UsedPercent,
	)

	for _, disk := range info.Disks {
		fmt.Printf("Disk (%s): %s / %s (%d%%) - %s\n",
			disk.Mountpoint,
			humanize.IBytes(disk.Used),
			humanize.IBytes(disk.Total),
			disk.UsedPercent,
			disk.FsType,
		)
	}
}
