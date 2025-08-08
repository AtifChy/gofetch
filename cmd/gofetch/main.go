package main

import (
	"context"
	"fmt"
	"log"
	"os/user"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
	"golang.org/x/sys/windows/registry"
)

type Info struct {
	OS              string
	Hostname        string
	KernelVersion   string
	PlatformVersion string
	CPUModel        string
	CPUCount        int32
	GPUs            []string
	Memory          Memory
	Disks           []Disk
	Uptime          time.Time
}

type Disk struct {
	Mountpoint string
	Total      uint64
	Free       uint64
	Used       uint64
}

type Memory struct {
	Total uint64
	Free  uint64
	Used  uint64
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
		collectGPUInfo,
		collectMemoryInfo,
		collectDiskInfo,
	}

	// Start collectors concurrently
	wg.Add(len(collectors))
	for _, collector := range collectors {
		go func(collector func(context.Context) (*Info, error)) {
			defer wg.Done()
			info, err := collector(ctx)
			if err != nil {
				log.Fatalf("Error collecting info: %s", err)
				return
			}
			infoChan <- info
		}(collector)
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
	}
	if src.Hostname != "" {
		dst.Hostname = src.Hostname
	}
	if src.KernelVersion != "" {
		dst.KernelVersion = src.KernelVersion
	}
	if src.PlatformVersion != "" {
		dst.PlatformVersion = src.PlatformVersion
	}
	if src.CPUModel != "" {
		dst.CPUModel = src.CPUModel
		dst.CPUCount = src.CPUCount
	}
	if src.GPUs != nil {
		dst.GPUs = append(dst.GPUs, src.GPUs...)
	}
	if src.Memory.Total > 0 {
		dst.Memory.Total = src.Memory.Total
		dst.Memory.Used = src.Memory.Used
		dst.Memory.Free = src.Memory.Free
	}
	if len(src.Disks) > 0 {
		dst.Disks = append(dst.Disks, src.Disks...)
	}
	if !src.Uptime.IsZero() {
		dst.Uptime = src.Uptime
	}
}

func collectHostInfo(ctx context.Context) (*Info, error) {
	h, err := host.InfoWithContext(ctx)
	if err != nil {
		return nil, err
	}
	return &Info{
		OS:              h.Platform + " " + h.KernelArch,
		Hostname:        h.Hostname,
		KernelVersion:   strings.Split(h.KernelVersion, " ")[0],
		PlatformVersion: h.PlatformVersion,
		Uptime:          time.Unix(int64(h.BootTime), 0),
	}, nil
}

func collectCPUInfo(ctx context.Context) (*Info, error) {
	infos, err := cpu.InfoWithContext(ctx)
	if err != nil {
		return nil, err
	}
	return &Info{
		CPUModel: strings.TrimSpace(infos[0].ModelName),
		CPUCount: infos[0].Cores,
	}, nil
}

func collectGPUInfo(_ context.Context) (*Info, error) {
	if runtime.GOOS != "windows" {
		return &Info{}, nil // GPU info collection is not implemented for non-Windows systems
	}

	var gpus []string
	key, err := registry.OpenKey(
		registry.LOCAL_MACHINE,
		`SYSTEM\CurrentControlSet\Control\Class\{4d36e968-e325-11ce-bfc1-08002be10318}`,
		registry.ENUMERATE_SUB_KEYS|registry.QUERY_VALUE,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to open registry key: %w", err)
	}

	subkeys, err := key.ReadSubKeyNames(-1)
	if err != nil {
		return nil, fmt.Errorf("failed to read subkeys: %w", err)
	}

	for _, subkey := range subkeys {
		if len(subkey) != 4 || subkey[0] < '0' || subkey[0] > '9' {
			continue
		}

		subkeyPath := `SYSTEM\CurrentControlSet\Control\Class\{4d36e968-e325-11ce-bfc1-08002be10318}\` + subkey
		gpuKey, err := registry.OpenKey(
			registry.LOCAL_MACHINE,
			subkeyPath,
			registry.QUERY_VALUE,
		)
		if err != nil {
			continue
		}

		desc, _, err := gpuKey.GetStringValue("DriverDesc")
		if err == nil && desc != "" {
			gpus = append(gpus, desc)
		}

		if err = gpuKey.Close(); err != nil {
			return nil, fmt.Errorf("failed to close registry key: %w", err)
		}
	}

	return &Info{GPUs: gpus}, nil
}

func collectMemoryInfo(ctx context.Context) (*Info, error) {
	m, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		return nil, err
	}
	return &Info{
		Memory: Memory{
			Total: m.Total,
			Used:  m.Used,
			Free:  m.Free,
		},
	}, nil
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
			Mountpoint: part.Mountpoint,
			Total:      usage.Total,
			Free:       usage.Free,
			Used:       usage.Used,
		})
	}
	return &Info{Disks: out}, nil
}

func kernelName() string {
	switch runtime.GOOS {
	case "linux":
		return "Linux"
	case "darwin":
		return "Darwin"
	case "windows":
		return "Windows_NT"
	default:
		return "Unknown"
	}
}

func displayInfo(info *Info) {
	user, err := user.Current()
	if err != nil {
		log.Fatalf("Error getting current user: %s\n", err)
	}
	username := user.Username
	if runtime.GOOS == "windows" {
		username = strings.Split(user.Username, "\\")[1]
	}
	title := username + "@" + info.Hostname

	fmt.Printf("%s\n", title)
	fmt.Println(strings.Repeat("-", len(title)))
	fmt.Printf("OS: %s\n", info.OS)
	fmt.Printf("Kernel: %s %s (%s)\n", kernelName(), info.KernelVersion, info.PlatformVersion)
	fmt.Printf("Uptime: %s\n", humanize.Time(info.Uptime))
	fmt.Printf("CPU: %s (%d)\n", info.CPUModel, info.CPUCount)
	for i, gpu := range info.GPUs {
		fmt.Printf("GPU %d: %s\n", i+1, gpu)
	}
	fmt.Printf("Memory: %s total, %s used, %s free\n",
		humanize.IBytes(info.Memory.Total),
		humanize.IBytes(info.Memory.Used),
		humanize.IBytes(info.Memory.Free),
	)
	for _, disk := range info.Disks {
		fmt.Printf("Disk (%s): %s total, %s used, %s free\n",
			disk.Mountpoint,
			humanize.IBytes(disk.Total),
			humanize.IBytes(disk.Used),
			humanize.IBytes(disk.Free),
		)
	}
}
