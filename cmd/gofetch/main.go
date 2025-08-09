package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/hako/durafmt"

	"github.com/AtifChy/gofetch/internal/collector"
	"github.com/AtifChy/gofetch/internal/cpu"
	"github.com/AtifChy/gofetch/internal/disk"
	"github.com/AtifChy/gofetch/internal/display"
	"github.com/AtifChy/gofetch/internal/gpu"
	"github.com/AtifChy/gofetch/internal/host"
	"github.com/AtifChy/gofetch/internal/memory"
	"github.com/AtifChy/gofetch/internal/types"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// List of collector functions
	collectors := []collector.CollectorFunc{
		host.CollectHostInfo,
		cpu.CollectCPUInfo,
		display.CollectDisplayInfo,
		gpu.CollectGPUInfo,
		memory.CollectMemoryInfo,
		disk.CollectDiskInfo,
	}

	info, err := collector.Collect(ctx, collectors)
	if err != nil {
		log.Fatalf("Error collecting system info: %s", err)
	}

	// Print info
	printInfo(info)
}

func printInfo(info *types.Info) {
	title := fmt.Sprintf("%s@%s", info.Host.Username, info.Host.Hostname)
	fmt.Printf("%s\n", title)
	fmt.Println(strings.Repeat("-", len(title)))

	fmt.Printf("OS: %s\n", info.Host.OS)
	fmt.Printf("Kernel: %s (%s)\n", info.Host.Kernel, info.Host.PlatformVersion)
	fmt.Printf("Uptime: %s\n", durafmt.Parse(info.Host.Uptime).LimitFirstN(2).String())

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
