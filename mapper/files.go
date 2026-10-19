package mapper

import (
	"path/filepath"
	"regexp"
	"fmt"
	"os"
)

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