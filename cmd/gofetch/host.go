package main

import (
	"context"
	"log"
	"os"
	"os/user"
	"runtime"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v4/host"
)

func collectHostInfo(ctx context.Context) (*Info, error) {
	hostInfo, err := host.InfoWithContext(ctx)
	if err != nil {
		return nil, err
	}
	return &Info{
		OS:              hostInfo.Platform + " " + hostInfo.KernelArch,
		KernelVersion:   strings.Split(hostInfo.KernelVersion, " ")[0],
		PlatformVersion: hostInfo.PlatformVersion,
		Uptime:          time.Duration(hostInfo.Uptime) * time.Second,
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

func getTitle() string {
	user, err := user.Current()
	if err != nil {
		log.Fatalf("Error getting current user: %s\n", err)
	}
	username := user.Username
	if runtime.GOOS == "windows" {
		username = strings.Split(user.Username, "\\")[1]
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("Error getting hostname: %s\n", err)
	}

	return username + "@" + hostname
}
