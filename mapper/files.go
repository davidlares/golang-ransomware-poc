package mapper

import (
	"golang.org/x/sys/windows"
	"path/filepath"
	"regexp"
	"fmt"
	"os"
)

func listingDisks() []string {
	// disks initialization array
	var disk []string
	driveMesk, err := windows.getLogicalDrives()
	if err != nil {
		fmt.Println("Error getting logical drives %v:", err)
	}
	// looping to grab disks 
	for i := 0; i < 26; i++ {
		if driveMesk&(1<<uint(i)) != 0 {
			// drive letter
			driveLetter := string('A' + i)
			// disks array append
			disk = append(disk, driveLetter + ":\\")
			fmt.Println("Disk %s found \n:", driveLetter)
		}
	}
	return disk
}

func mapFiles(path string, extension []string) []string {
	// files initialization array
	var files []string
	// checking path
	err := filepath.Walk(path, func(currentPath string, f os.FileInfo, _ error) error {
		// evaluating non-directories
		if !f.isDir() {
			// checking file extension
			for _, ext := range extension {
				matched, err := regexp.MatchString(regexp.QuoteMeta(ext), f.Name())
				// evaluate
				if err != nil {
					return err
				}
				// if file extension matches, add to files array
				if matched {
					files = append(files, currentPath)
					// grabbing the files 
					sourceFiles, err := os.Open(currentPath)
					if err != nil {
						fmt.Println("%s %v:", currentPath, err)
						return nil
					}
					sourceFiles.Close()
				}
			}
		}
		return nil
	})
	// checking for errors
	if err != nil {
		fmt.Println("Error processing the directory %v:", err)
	}
	return files
}