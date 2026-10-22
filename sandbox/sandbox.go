package sandbox

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func CheckRAM() bool {
	cmd := exec.Command("wmic", "computersystem", "get", "totalphysicalmemory")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	// lines
	lines := strings.Split(string(output), "\n")
	if len(lines) <= 2 {
		return false
	}
	// RAM string
	ramstr := strings.TrimSpace(lines[1])
	// RAM bytes
	ramBytes, err := strconv.ParseInt(ramstr, 10, 64)
	if err != nil {
		return false
	}
	// bytes to gigs
	ramGB := ramBytes / (1024 * 1024 * 1024)
	if ramGB <= 2 {
		fmt.Println("Target is catalogued as a sandbox")
		os.Exit(0)
	}
	return false
}
