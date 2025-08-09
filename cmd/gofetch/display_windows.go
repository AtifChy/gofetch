//go:build windows

package main

import (
	"context"
	"errors"
	"syscall"
	"unsafe"
)

var (
	user32                  = syscall.NewLazyDLL("user32.dll")
	gdi32                   = syscall.NewLazyDLL("gdi32.dll")
	procGetMonitorInfoW     = user32.NewProc("GetMonitorInfoW")
	procCreateDCW           = gdi32.NewProc("CreateDCW")
	procDeleteDC            = gdi32.NewProc("DeleteDC")
	procGetDeviceCaps       = gdi32.NewProc("GetDeviceCaps")
	procEnumDisplayMonitors = user32.NewProc("EnumDisplayMonitors")
)

const (
	HORZRES              = 8
	VERTRES              = 10
	VREFRESH             = 116
	MONITORINFOF_PRIMARY = 0x00000001
	CCHDEVICENAME        = 32
)

type monitorInfoEx struct {
	CbSize    uint32
	RcMonitor rect
	RcWork    rect
	DwFlags   uint32
	SzDevice  [CCHDEVICENAME]uint16
}

func collectDisplayInfo(_ context.Context) (*Info, error) {
	var out []Display

	callback := syscall.NewCallback(func(hMon, hdc, lprc, lParam uintptr) uintptr {
		var mi monitorInfoEx
		mi.CbSize = uint32(unsafe.Sizeof(mi))

		ok, _, _ := procGetMonitorInfoW.Call(hMon, uintptr(unsafe.Pointer(&mi)))
		if ok == 0 {
			return 1 // skip this monitor
		}

		driver, _ := syscall.UTF16PtrFromString("DISPLAY")
		device := &mi.SzDevice[0]

		hdcMon, _, _ := procCreateDCW.Call(
			uintptr(unsafe.Pointer(driver)),
			uintptr(unsafe.Pointer(device)),
			0,
			0,
		)
		if hdcMon == 0 {
			return 1 // skip if failed
		}

		defer procDeleteDC.Call(hdcMon)

		w, _, _ := procGetDeviceCaps.Call(hdcMon, uintptr(HORZRES))
		h, _, _ := procGetDeviceCaps.Call(hdcMon, uintptr(VERTRES))
		f, _, _ := procGetDeviceCaps.Call(hdcMon, uintptr(VREFRESH))

		info := Display{
			Width:       int32(w),
			Height:      int32(h),
			RefreshRate: int32(f),
			IsPrimary:   (mi.DwFlags & MONITORINFOF_PRIMARY) != 0,
		}
		out = append(out, info)

		return 1
	})

	ret, _, _ := procEnumDisplayMonitors.Call(0, 0, callback, 0)
	if ret == 0 {
		return nil, errors.New("failed to enumerate display monitors")
	}

	return &Info{
		Displays: out,
	}, nil
}
