## Golang Ransomware PoC

This is a modular command-line tool built in Go designed to demonstrate security concepts, cryptographic operations (AES-GCM, RSA, XOR), environment validation, and system utility interactions on Windows Operating Systems.

Internally performs symmetric encryption using AES-256 (GCM mode) combined with custom XOR preprocessing and then an asymmetric wrapping of the AES session key using RSA public/private key pairs.

## Ransomware as a malware mechanism

Ransomware is a type of malware that blocks access to your data and demands payments to restore it. This type of attack has become a major problem in recent years, and the idea behind it is to collect a ransom for encrypted data via asymmetric encryption

The ransomware enters your system through phishing emails, infected file downloads, or vulnerability exploits, which basically follow MITRE tactics such as (recon, initial access, persistence, escalation, lateral movement, etc.). Once installed, it encrypts your files, making them inaccessible. It uses a strong encryption algorithm, making decryption impossible without the private keys.

Later, a message is displayed demanding payment in crypto, and the idea behind it is to "pay" to recover access to the encrypted files; if not, usually the victim is threatened with permanent destruction within a certain time frame

## Project folder architecture and responsibility

1. `cmd/`: Manages command-line interface routing and flags using `spf13/cobra` (encrypt, decrypt, and file listing utilities).

2. `aes/`: Handles symmetric AES key generation, persistence, and file storage.

3. `encrypt/`: Implements multi-layered file transformation using XOR encoding combined with AES-GCM encryption.

4. `decrypt/`: Reverses the encryption pipeline to restore original files securely.

5. `rsa/`: Manages asymmetric key pair generation (PKCS#1) and wraps/unwraps the symmetric AES key using RSA-OAEP/PKCS#1v15.

6. `layout/`: Interacts with Windows system DLLs (user32.dll) via golang.org/x/sys/windows to inspect active keyboard layouts (e.g., bypassing execution based on regional language layouts).

7. `sandbox/`: Performs basic environment and RAM validation checks for sandbox evasion analysis.

8. `mapper/`: Handles system drive discovery and recursive file-system traversal filtered by target file extensions.

9. `logs/`: Manages execution logging and post-execution log sanitization.

## Hybrid encryption

Why hybrid? Because we have symmetric encryption (AES) and asymmetric encryption (RSA). We gain speed with file encryption because we use the same key with AES allowing ransomware to lock thousands of files in seconds rather than hours.

So, if we stick to a program that uses pure AES, it either needs a hardcoded key or a static key file (like `aes_key.txt`).The key could reverse-engineered and write a decryption tool instantly, so, with that in mind, RSA comes to the rescue. 

The malware generates a unique, random AES session key for that specific victim and then encrypts files using that AES key. Finally, it encrypts only the AES key using an RSA Public Key embedded in the binary, producing a file like `aes_key.txt.enc_rsa`

So, even if the victim recovers the encrypted files and the encrypted AES key, they cannot decrypt the AES key without the corresponding RSA Private Key, which is securely kept offline by the attacker.

### The decryption side 

The pipeline works in reverse, to successfully restore the files, the process must first recover the master AES key using the asymmetric RSA private key, and then use that recovered AES key along with the XOR key to unseal each encrypted file.

Because the AES session key was wrapped using your RSA public key during encryption, you need the corresponding RSA private key (`private.pem`) to unwrap it before you can decrypt any files. And once the AES key is recovered, the `decrypt` command iterates over the locked files, reads the stored nonce and ciphertext, validates the auth tag via AES-GCM, and strips the XOR layer.

## Useful concepts 

1. *Marshaling* is just a fancy word for Serialization. It is the process of taking a memory object (like a Go struct) and turning it into a flat stream of bytes that can be saved to a disk or sent over a network.

2. *X.509* is the standard for public key infrastructure (PKI), you can see it as the "umbrella" standard.

3. *PKCS#1* stands for Public Key Cryptography Standards #1 and is a specific format designed only for RSA private keys. It is the "old school" but very standard way to store an RSA key.

4. *PKIX* (Public Key Infrastructure (X.509)) is a modern, flexible format. It’s used for public keys because it can hold RSA keys, ECDSA keys, or Ed25519 keys without changing the structure

5. The Private Key is the "Master". It contains all the mathematical components needed to derive the Public Key; this is what is shared with the world

## Dependencies

Check `go.mod` and use `go mod tidy` to download and install the dependencies

## Usage 

You can simply run the same `main.go` with two possible args: `encrypt` or `decrypt`

Example:

`go run main.go encrypt` or `go run main.go decrypt`

## Credits
[David Lares S](https://davidlares.com)

## License
[MIT](https://opensource.org/licenses/MIT)