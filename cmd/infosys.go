package cmd

import (
	"davidlares/ransomware/infosys"
	"fmt"

	"github.com/spf13/cobra"
)

var infoSysCmd = &cobra.Command{
	Use:   "infosys",
	Short: "Getting information of the system",
	Long:  `A simple script to get information of the system built with GoLang.`,
	Run: func(cmd *cobra.Command, args []string) {
		username := infosys.GetUserName()
		hostname := infosys.GetHostName()
		ips, _ := infosys.GetIPAddress()
		// printing
		fmt.Println("System information")
		fmt.Println("Username: ", username)
		fmt.Println("Hostname: ", hostname)
		fmt.Println("Address: ", ips)
	},
}

// adding child commands to the encrypt command (= flags)
func init() {
	rootCmd.AddCommand(infoSysCmd)
}
