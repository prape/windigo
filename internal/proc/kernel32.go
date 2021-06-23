package proc

import (
	"syscall"
)

var (
	kernel32 = syscall.NewLazyDLL("kernel32.dll")

	CloseHandle                     = kernel32.NewProc("CloseHandle")
	CreateDirectory                 = kernel32.NewProc("CreateDirectoryW")
	CreateFile                      = kernel32.NewProc("CreateFileW")
	CreateFileMapping               = kernel32.NewProc("CreateFileMappingW")
	CreateProcess                   = kernel32.NewProc("CreateProcessW")
	DeleteFile                      = kernel32.NewProc("DeleteFileW")
	ExitProcess                     = kernel32.NewProc("ExitProcess")
	FileTimeToSystemTime            = kernel32.NewProc("FileTimeToSystemTime")
	FindClose                       = kernel32.NewProc("FindClose")
	FindFirstFile                   = kernel32.NewProc("FindFirstFileW")
	FindNextFile                    = kernel32.NewProc("FindNextFileW")
	FreeLibrary                     = kernel32.NewProc("FreeLibrary")
	GetCurrentProcessId             = kernel32.NewProc("GetCurrentProcessId")
	GetCurrentThreadId              = kernel32.NewProc("GetCurrentThreadId")
	GetDynamicTimeZoneInformation   = kernel32.NewProc("GetDynamicTimeZoneInformation")
	GetExitCodeProcess              = kernel32.NewProc("GetExitCodeProcess")
	GetExitCodeThread               = kernel32.NewProc("GetExitCodeThread")
	GetFileAttributes               = kernel32.NewProc("GetFileAttributesW")
	GetFileSizeEx                   = kernel32.NewProc("GetFileSizeEx")
	GetModuleFileName               = kernel32.NewProc("GetModuleFileNameW")
	GetModuleHandle                 = kernel32.NewProc("GetModuleHandleW")
	GetProcAddress                  = kernel32.NewProc("GetProcAddress")
	GetProcessId                    = kernel32.NewProc("GetProcessId")
	GetProcessIdOfThread            = kernel32.NewProc("GetProcessIdOfThread")
	GetProcessTimes                 = kernel32.NewProc("GetProcessTimes")
	GetStartupInfo                  = kernel32.NewProc("GetStartupInfoW")
	GetSystemInfo                   = kernel32.NewProc("GetSystemInfo")
	GetSystemTime                   = kernel32.NewProc("GetSystemTime")
	GetSystemTimeAsFileTime         = kernel32.NewProc("GetSystemTimeAsFileTime")
	GetSystemTimePreciseAsFileTime  = kernel32.NewProc("GetSystemTimePreciseAsFileTime")
	GetSystemTimes                  = kernel32.NewProc("GetSystemTimes")
	GetThreadId                     = kernel32.NewProc("GetThreadId")
	GetThreadTimes                  = kernel32.NewProc("GetThreadTimes")
	GetTickCount64                  = kernel32.NewProc("GetTickCount64")
	GetTimeZoneInformation          = kernel32.NewProc("GetTimeZoneInformation")
	GetTimeZoneInformationForYear   = kernel32.NewProc("GetTimeZoneInformationForYear")
	LoadLibrary                     = kernel32.NewProc("LoadLibraryW")
	MapViewOfFile                   = kernel32.NewProc("MapViewOfFile")
	MulDiv                          = kernel32.NewProc("MulDiv")
	QueryPerformanceCounter         = kernel32.NewProc("QueryPerformanceCounter")
	QueryPerformanceFrequency       = kernel32.NewProc("QueryPerformanceFrequency")
	ReadFile                        = kernel32.NewProc("ReadFile")
	SetEndOfFile                    = kernel32.NewProc("SetEndOfFile")
	SetFilePointerEx                = kernel32.NewProc("SetFilePointerEx")
	Sleep                           = kernel32.NewProc("Sleep")
	SystemTimeToFileTime            = kernel32.NewProc("SystemTimeToFileTime")
	SystemTimeToTzSpecificLocalTime = kernel32.NewProc("SystemTimeToTzSpecificLocalTime")
	TzSpecificLocalTimeToSystemTime = kernel32.NewProc("TzSpecificLocalTimeToSystemTime")
	UnmapViewOfFile                 = kernel32.NewProc("UnmapViewOfFile")
	VerifyVersionInfo               = kernel32.NewProc("VerifyVersionInfoW")
	VerSetConditionMask             = kernel32.NewProc("VerSetConditionMask")
	WaitForSingleObject             = kernel32.NewProc("WaitForSingleObject")
	WriteFile                       = kernel32.NewProc("WriteFile")
)
