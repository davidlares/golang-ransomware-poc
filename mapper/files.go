package mapper

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"golang.org/x/sys/windows"
)

var (
	Wg sync.WaitGroup
)

func ListingDisks() []string {
	// disks initialization array
	var disks []string
	driveMesk, err := windows.GetLogicalDrives()
	if err != nil {
		fmt.Printf("Error getting logical drives %v:", err)
	}
	// looping to grab disks
	for i := 0; i < 26; i++ {
		if driveMesk&(1<<uint(i)) != 0 {
			// drive letter
			driveLetter := string('A' + i)
			// disks array append
			disks = append(disks, driveLetter+":\\")
			fmt.Printf("Disk %s found \n:", driveLetter)
		}
	}
	return disks
}

func MapFiles(path string, blackList []string) []string {
	// files initialization array
	var files []string
	// excluding directories for testing
	excludedDir := []string{"Windows", "Program Files", "Program Data (x86)", "boot", "source", "AppData", "ProgramData"}
	// checking path
	err := filepath.Walk(path, func(currentPath string, f os.FileInfo, err error) error {
		// evaluating errors
		if err != nil {
			return nil
		}
		// looping and finding files in excluded directories
		for _, dir := range excludedDir {
			if f.IsDir() && f.Name() == dir {
				return filepath.SkipDir
			}
		}
		// evaluating non-directories
		if !f.IsDir() {
			// checking file extension
			ext := filepath.Ext(f.Name())
			for _, blackExt := range blackList {
				if ext == blackExt {
					return nil
				}
			}
			// appending
			files = append(files, currentPath)
		}
		return nil
	})
	// checking for errors
	if err != nil {
		fmt.Printf("Error processing the directory %v:", err)
	}
	return files
}
