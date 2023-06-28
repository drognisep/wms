package data

import (
	"errors"
	"syscall"
	"unsafe"
)

var (
	kernel32            *syscall.DLL
	getDiskFreeSpaceExW *syscall.Proc
)

func init() {
	kernel32 = syscall.MustLoadDLL("kernel32.dll")
	getDiskFreeSpaceExW = kernel32.MustFindProc("GetDiskFreeSpaceExW")
}

// GetSpaceFree returns the total and free bytes available on the current disk.
// References https://learn.microsoft.com/en-us/windows/win32/api/fileapi/nf-fileapi-getdiskfreespaceexw
func GetSpaceFree() (DiskSpace, DiskSpace, error) {
	var (
		driveSpace uint64
		driveFree  uint64
	)
	rval, _, _ := getDiskFreeSpaceExW.Call(
		uintptr(unsafe.Pointer(nil)),
		uintptr(unsafe.Pointer(&driveFree)),
		uintptr(unsafe.Pointer(&driveSpace)),
		uintptr(unsafe.Pointer(nil)),
	)
	if rval == 0 {
		return 0, 0, errors.New("failed to get total/free disk space")
	}
	return DiskSpace(driveSpace), DiskSpace(driveFree), nil
}
