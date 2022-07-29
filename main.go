package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"github.com/fulldump/goconfig"
	"golang.org/x/crypto/ssh/terminal"
)

type config struct {
	File string
	Open bool
}

func readPassword() []byte {
	fmt.Fprint(os.Stderr, "Enter password: ")
	password, err := terminal.ReadPassword(0)
	fmt.Fprintln(os.Stderr, "")
	if err != nil {
		panic(err)
	}
	return password
}

func GetMD5Hash(text []byte) []byte {
	hasher := md5.New()
	hasher.Write(text)
	return []byte(hex.EncodeToString(hasher.Sum(nil)))
}

func main() {

	c := config{}
	goconfig.Read(&c)

	// todo: validate config

	password := readPassword()
	key := GetMD5Hash(password)

	// content, err := os.ReadFile(c.File)
	// if err != nil {
	// 	panic(err)
	// }

	gozFile(key, c.File, c.Open)

}

func gozFile(key []byte, filename string, open bool) {

	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	s, err := f.Stat()
	if err != nil {
		panic(err)
	}

	content, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	var transformed []byte
	if open {
		transformed = decrypt(key, content)
	} else {
		transformed = encrypt(key, content)
	}

	err = os.WriteFile(filename, transformed, s.Mode())
	if err != nil {
		panic(err)
	}
}

func encrypt(key, plaintext []byte) (ciphertext []byte) {

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err)
	}

	return gcm.Seal(nonce, nonce, plaintext, nil)
}

func decrypt(key, ciphertext []byte) (plaintext []byte) {

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}

	nonce := ciphertext[:gcm.NonceSize()]
	ciphertext = ciphertext[gcm.NonceSize():]
	plaintext, err = gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err)
	}

	return plaintext
}
