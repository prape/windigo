package api

import (
	"fmt"
	"syscall"
	"unsafe"
	"wingows/api/proc"
	c "wingows/consts"
)

// Wrapper to syscall.UTF16PtrFromString(), panics in error.
func StrToUtf16Ptr(s string) *uint16 {
	// We won't return an uintptr right away because it has no pointer semantics,
	// it's just a number, so pointed memory can be garbage-collected.
	// https://stackoverflow.com/a/51188315
	pstr, err := syscall.UTF16PtrFromString(s)
	if err != nil {
		panic(fmt.Sprintf("ToUtf16Ptr failed \"%s\": %s\n",
			s, err))
	}
	return pstr
}

// Wrapper to syscall.UTF16PtrFromString(), panics in error. A blank string will
// return a null pointer.
func StrToUtf16PtrBlankIsNil(s string) *uint16 {
	if s != "" {
		return StrToUtf16Ptr(s)
	}
	return nil
}

func boolToUintptr(b bool) uintptr {
	if b {
		return uintptr(1)
	}
	return uintptr(0)
}

//------------------------------------------------------------------------------

func HiWord(value uint32) uint16 {
	return uint16(value >> 16 & 0xffff)
}

func LoWord(value uint32) uint16 {
	return uint16(value)
}

func HiByte(value uint16) uint8 {
	return uint8(value >> 8 & 0xff)
}

func LoByte(value uint16) uint8 {
	return uint8(value)
}

//------------------------------------------------------------------------------

func InitCommonControls() {
	syscall.Syscall(proc.InitCommonControls.Addr(), 0,
		0, 0, 0)
}

// Available in Windows 10, version 1607.
func GetDpiForSystem() uint32 {
	ret, _, _ := syscall.Syscall(proc.GetDpiForSystem.Addr(), 0,
		0, 0, 0)
	return uint32(ret)
}

func GetSystemMetrics(index c.SM) int32 {
	ret, _, _ := syscall.Syscall(proc.GetSystemMetrics.Addr(), 1,
		uintptr(index), 0, 0)
	return int32(ret)
}

// Multiplies two 32-bit values and then divides the 64-bit result by a third
// 32-bit value. The final result is rounded to the nearest integer.
func MulDiv(number, numerator, denominator int32) int32 {
	ret, _, _ := syscall.Syscall(proc.MulDiv.Addr(), 3,
		uintptr(number), uintptr(numerator), uintptr(denominator))
	return int32(ret)
}

func PostQuitMessage(exitCode int32) {
	syscall.Syscall(proc.PostQuitMessage.Addr(), 1, uintptr(exitCode), 0, 0)
}

// Available in Windows 10, version 1703.
func SetProcessDpiAwarenessContext(value c.DPI_AWARE_CTX) {
	ret, _, lerr := syscall.Syscall(proc.SetProcessDpiAwarenessContext.Addr(), 1,
		uintptr(value), 0, 0)
	if ret == 0 {
		panic(fmt.Sprintf("SetProcessDpiAwarenessContext failed: %d %s\n",
			lerr, lerr.Error()))
	}
}

// Available in Windows Vista.
func SetProcessDPIAware() {
	ret, _, _ := syscall.Syscall(proc.SetProcessDPIAware.Addr(), 0,
		0, 0, 0)
	if ret == 0 {
		panic("SetProcessDPIAware failed.")
	}
}

func SystemParametersInfo(action c.SPI, param uint32,
	pvParam unsafe.Pointer, winIni uint32) {

	ret, _, lerr := syscall.Syscall6(proc.SystemParametersInfo.Addr(), 4,
		uintptr(action), uintptr(param), uintptr(pvParam), uintptr(winIni),
		0, 0)
	if ret == 0 {
		panic(fmt.Sprintf("SystemParametersInfo failed: %d %s\n",
			lerr, lerr.Error()))
	}
}
