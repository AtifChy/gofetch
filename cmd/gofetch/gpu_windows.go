//go:build windows

package main

import (
	"context"
	"fmt"

	"golang.org/x/sys/windows/registry"
)

func collectGPUInfo(_ context.Context) (*Info, error) {
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

	var out []GPU

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

		var gpu GPU

		name, _, err := gpuKey.GetStringValue("DriverDesc")
		if err == nil && name != "" {
			gpu.Name = name
		}

		vram, _, err := gpuKey.GetIntegerValue("HardwareInformation.qwMemorySize")
		if err == nil {
			gpu.VRAM = uint64(vram)
		}

		out = append(out, gpu)

		if err = gpuKey.Close(); err != nil {
			return nil, fmt.Errorf("failed to close registry key: %w", err)
		}
	}

	return &Info{GPUs: out}, nil
}
