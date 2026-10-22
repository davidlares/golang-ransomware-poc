package kill

import (
	"fmt"
	"log"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	modKernel32          = windows.NewLazyDLL("kernel32.dll")
	procOpenProcess      = modKernel32.NewProc("OpenProcess")
	procTerminateProcess = modKernel32.NewProc("TerminateProcess")
)

const (
	PROCESS_TERMINATE = 0x0001
)

func TerminateProcess(pid uint32) error {
	handle, _, err := procOpenProcess.Call(uintptr(PROCESS_TERMINATE), uintptr(0), uintptr(pid))
	// check if could be opened
	if handle == 0 {
		return fmt.Errorf("Failed to open process with PID %d: %v", pid, err)
	}
	// closing handle
	defer syscall.CloseHandle(syscall.Handle(handle))
	// terminating the process (call instruction)
	ret, _, err := procTerminateProcess.Call(handle, uintptr(0))
	if ret == 0 {
		return fmt.Errorf("failed to terminate process with PID %d: %v", pid, err)
	}

	return nil
}

func GetProcessList() ([]string, error) {
	// custom process list
	processList := []string{
		"Chrome", "msedge", "calc",
	}
	return processList, nil
}

func KillProcessesByName(targetProcesses []string) {
	// snapshot created
	handle, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		log.Fatalf("Error creating the snapshot: %v", err)
	}
	// then (closing syscall)
	defer syscall.CloseHandle(syscall.Handle(handle))
	// calling the entry
	var entry windows.ProcessEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))
	// then
	err = windows.Process32First(handle, &entry)
	if err != nil {
		log.Fatalf("Error getting the first process: %v", err)
	}
	// getting the first process of the list
	for {
		// converting the null terminate uint16 array to standard Go string
		exeName := windows.UTF16ToString(entry.ExeFile[:])
		// checking matching processes
		for _, target := range targetProcesses {
			// compare case-insensitivity
			if strings.EqualFold(exeName, target) || strings.EqualFold(exeName, target+".exe") {
				// terminating
				fmt.Printf("Terminating process '%s' (PID: %d)...\n", exeName, entry.ProcessID)
				if err := TerminateProcess(entry.ProcessID); err != nil {
					log.Printf("Failed to terminate %s: %v", exeName, err)
				}
			}
			// moving to the next proess
			err = windows.Process32Next(handle, &entry)
			if err != nil {
				break
			}
		}
	}
}
