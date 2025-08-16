package main

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"dario.cat/mergo"
	"github.com/dustin/go-humanize"
	"github.com/hako/durafmt"

	"github.com/AtifChy/gofetch/internal/collector"
	data "github.com/AtifChy/gofetch/internal/config"
	"github.com/AtifChy/gofetch/internal/cpu"
	"github.com/AtifChy/gofetch/internal/disk"
	"github.com/AtifChy/gofetch/internal/display"
	"github.com/AtifChy/gofetch/internal/gpu"
	"github.com/AtifChy/gofetch/internal/host"
	"github.com/AtifChy/gofetch/internal/logo"
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

	config := data.LoadDefaultConfig()

	if userConfig, err := data.LoadConfig("config.json"); err == nil {
		err = mergo.Merge(&config, userConfig, mergo.WithOverride)
		if err != nil {
			log.Fatalf("Error contructing config: %s", err)
		}
	} else {
		log.Printf("Error loading config: %s", err)
	}

	info, err := collector.Collect(ctx, collectors)
	if err != nil {
		log.Fatalf("Error collecting system info: %s", err)
	}

	logo, err := logo.GetLogo(config.Logo)
	if err != nil {
		log.Printf("Error getting logo: %s", err)
	}

	// Print info
	printInfo(info, logo)
}

func printInfo(info *types.Info, logo string) {
	infoLines := getInfoLines(info)
	logoLines := strings.Split(logo, "\n")

	// Regex to strip ANSI color codes for accurate width calculation
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`)

	logoLineWidths := make([]int, len(logoLines))
	maxLogoWidth := 0
	for i, line := range logoLines {
		visibleLine := ansiRegex.ReplaceAllString(line, "")
		logoLineWidths[i] = len(visibleLine)
		maxLogoWidth = max(maxLogoWidth, logoLineWidths[i])
	}

	// Find max lines
	var out strings.Builder
	maxLines := max(len(infoLines), len(logoLines))

	for i := range maxLines {
		var logoPart string
		var logoPartWidth int
		if i < len(logoLines) {
			logoPart = logoLines[i]
			logoPartWidth = logoLineWidths[i]
		}

		var infoPart string
		if i < len(infoLines) {
			infoPart = infoLines[i]
		}

		padding := max(maxLogoWidth-logoPartWidth, 0)
		padding += 5 // Add extra padding for better separation

		out.WriteString(logoPart)
		out.WriteString(strings.Repeat(" ", padding))
		out.WriteString(infoPart)
		out.WriteString("\n")
	}

	fmt.Println(out.String())
}

func getInfoLines(info *types.Info) []string {
	var lines []string

	labelColor := logo.ColorMap["yellow"] + logo.ColorMap["bold"]
	resetColor := logo.ColorMap["reset"]

	title := fmt.Sprintf("%s@%s", info.Host.Username, info.Host.Hostname)
	lines = append(lines, title)
	lines = append(lines, strings.Repeat("-", len(title)))

	lines = append(lines, fmt.Sprintf("%sOS%s: %s", labelColor, resetColor, info.Host.OS))
	lines = append(lines, fmt.Sprintf("%sKernel%s: %s (%s)", labelColor, resetColor, info.Host.Kernel, info.Host.PlatformVersion))
	lines = append(lines, fmt.Sprintf("%sUptime%s: %s", labelColor, resetColor, durafmt.Parse(info.Host.Uptime).LimitFirstN(2).String()))

	displayStr := fmt.Sprintf("%sDisplay%s: ", labelColor, resetColor)
	for i, display := range info.Displays {
		if display.IsPrimary {
			displayStr += "*"
		}
		displayStr += fmt.Sprintf("%dx%d @ %dHz", display.Width, display.Height, display.RefreshRate)
		if i < len(info.Displays)-1 {
			displayStr += ", labelColor, resetColor, "
		}
	}
	lines = append(lines, displayStr)

	lines = append(lines, fmt.Sprintf("%sCPU%s: %s (%d)", labelColor, resetColor, info.CPU.Model, info.CPU.Cores))

	for i, gpu := range info.GPUs {
		lines = append(lines, fmt.Sprintf("%sGPU %d%s: %s (%s)", labelColor, i+1, resetColor, gpu.Name, humanize.IBytes(gpu.VRAM)))
	}

	lines = append(lines, fmt.Sprintf("%sMemory%s: %s / %s (%d%%)", labelColor, resetColor,
		humanize.IBytes(info.Memory.Used),
		humanize.IBytes(info.Memory.Total),
		info.Memory.UsedPercent,
	))

	for _, disk := range info.Disks {
		lines = append(lines, fmt.Sprintf("%sDisk (%s)%s: %s / %s (%d%%) - %s", labelColor,
			disk.Mountpoint, resetColor,
			humanize.IBytes(disk.Used),
			humanize.IBytes(disk.Total),
			disk.UsedPercent,
			disk.FsType,
		))
	}

	return lines
}
