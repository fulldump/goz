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

var operation2string = map[bool]string{
	true:  "decrypt",
	false: "encrypt",
}

func main() {

	c := config{}
	goconfig.Read(&c)

	// todo: validate config

	password := readPassword()
	key := GetMD5Hash(password)

	err := gozFile(key, c.File, c.Open)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s: %s\n", operation2string[c.Open], c.File, err.Error())
	}
}

func gozFile(key []byte, filename string, open bool) (err error) {

	f, err := os.Open(filename)
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
