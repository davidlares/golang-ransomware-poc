package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ransomware",
	Short: "GoLang Ransomware PoC",
	Long:  `A simple Ransomware PoC built with GoLang.`,
}

// adding child commands to the root command (= flags)
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
