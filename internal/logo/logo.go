// Package logo provides access to embedded logos.
package logo

import (
	"embed"
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"github.com/AtifChy/gofetch/internal/config"
)

//go:embed ascii/*.txt
var logoFiles embed.FS

var ColorMap = map[string]string{
	"red":     "\033[31m",
	"green":   "\033[32m",
	"yellow":  "\033[33m",
	"blue":    "\033[34m",
	"magenta": "\033[35m",
	"cyan":    "\033[36m",

	"bold":  "\033[1m",
	"reset": "\033[0m",
}

func GetLogo(logoConfig config.LogoConfig) (string, error) {
	var logoName string
	switch runtime.GOOS {
	case "windows":
		logoName = "windows_11.txt"
	// Add other OS cases here
	default:
		// Fallback logo or empty string
		return "", fmt.Errorf("no logo for OS: %s", runtime.GOOS)
	}

	data, err := logoFiles.ReadFile("ascii/" + logoName)
	if err != nil {
		return "", fmt.Errorf("could not read logo file %s: %w", logoName, err)
	}

	// Normalize newlines
	logo := strings.ReplaceAll(string(data), "\r\n", "\n")

	for key, color := range logoConfig.Colors {
		placeholder := "${" + strconv.Itoa(key) + "}"
		if code, ok := ColorMap[strings.ToLower(color)]; ok {
			code = ColorMap["bold"] + code // Apply bold style
			logo = strings.ReplaceAll(logo, placeholder, code)
		}
	}

	logo = strings.ReplaceAll(logo, "\n", ColorMap["reset"]+"\n")

	return logo, nil
}
