package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/fulldump/goconfig"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/ssh/terminal"
)

var VERSION = "dev"

// Changing these values will make this tool backwards incompatible (unsuitable
// for decrypting files already encrypted with previous versions)
var (
	salt = []byte("848405a0-14df-11ed-9d64-23f5c6c8bd9c")
	hash = sha1.New
	iter = 4096
)

type config struct {
	File    string
	Dir     string
	Open    bool
	Version bool
}

// Source: https://go.dev/play/p/l-9IP1mrhA
func readPassword(prompt string) []byte {
	fmt.Fprint(os.Stderr, prompt)
	password, err := terminal.ReadPassword(0)
	fmt.Fprintln(os.Stderr, "")
	if err != nil {
		panic(err)
	}
	return password
}

var operation2string = map[bool]string{
	true:  "decrypt",
	false: "encrypt",
}

func main() {

	c := config{}
	goconfig.Read(&c)

	if c.Version {
		fmt.Println("Version:", VERSION)
		return
	}

	// todo: validate config

	password := readPassword("Enter password: ")
	if !c.Open {
		confirm := readPassword("Confirm password: ")
		if string(confirm) != string(password) {
			fmt.Fprintln(os.Stderr, "Password does not match")
			os.Exit(7)
		}
	}

	// Using Pasword Based Key Derivation Function to get a key based on password
	// Thanks to @psanford (github id) Peter Stanford for this: https://gophers.slack.com/archives/C02A3DRK6/p1659185759603659
	key := pbkdf2.Key(password, salt, iter, 32, hash)

	files := []string{}

	if c.Dir == "" {
		files = append(files, c.File)
	} else {
		filepath.Walk(c.Dir, func(path string, info fs.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			if strings.HasPrefix(info.Name(), ".") {
				return nil
			}

			files = append(files, path)
			return nil
		})
	}

	for _, file := range files {
		err := gozFile(key, file, c.Open)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s: %s\n", operation2string[c.Open], file, err.Error())
		}
	}

}

func gozFile(key []byte, filename string, open bool) (err error) {

	f, err := os.Open(filename) // TODO: open for exclusive access
	if err != nil {
		return fmt.Errorf("os.Open: %w", err)
	}

	s, err := f.Stat()
	if err != nil {
		return fmt.Errorf("f.Stat: %w", err)
	}

	content, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("io.ReadAll: %w", err)
	}
	f.Close() // TODO: handle err

	var transformed []byte
	if open {
		transformed, err = decrypt(key, content)
	} else {
		transformed, err = encrypt(key, content)
	}
	if err != nil {
		return fmt.Errorf("f.Stat: %w", err)
	}

	err = os.WriteFile(filename, transformed, s.Mode())
	if err != nil {
		return fmt.Errorf("os.WriteFile: %w", err)
	}

	return
}

func encrypt(key, plaintext []byte) (ciphertext []byte, err error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("aes.NewCipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("cipher.NewGCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("nonce: %w", err)
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func decrypt(key, ciphertext []byte) (plaintext []byte, err error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("aes.NewCipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("cipher.NewGCM: %w", err)
	}

	nonce := ciphertext[:gcm.NonceSize()]
	ciphertext = ciphertext[gcm.NonceSize():]
	plaintext, err = gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("gcm.Open: %w", err)
	}

	return
}
