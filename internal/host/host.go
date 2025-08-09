// Package host provides functionality to collect host information
package host

import (
	"context"
	"log"
	"os/user"
	"runtime"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v4/host"

	"github.com/AtifChy/gofetch/internal/types"
)

func CollectHostInfo(ctx context.Context) (*types.Info, error) {
	hostInfo, err := host.InfoWithContext(ctx)
	if err != nil {
		return nil, err
	}

	host := types.Host{
		Username:        getUsername(),
		Hostname:        hostInfo.Hostname,
		OS:              hostInfo.Platform + " " + hostInfo.KernelArch,
		Kernel:          getKernelName() + " " + strings.Split(hostInfo.KernelVersion, " ")[0],
		PlatformVersion: hostInfo.PlatformVersion,
		Uptime:          time.Duration(hostInfo.Uptime) * time.Second,
	}

	return &types.Info{
		Host: host,
	}, nil
}

func getKernelName() string {
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

func getUsername() string {
	user, err := user.Current()
	if err != nil {
		log.Fatalf("Error getting current user: %s\n", err)
	}

	username := user.Username
	if runtime.GOOS == "windows" {
		username = strings.Split(user.Username, "\\")[1]
	}

	return username
}
