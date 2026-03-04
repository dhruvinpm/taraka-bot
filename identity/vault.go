package identity

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Vault struct {
	path string
	key  []byte
	data map[string]string
}

func NewVault(path, keyStr string) (*Vault, error) {
	h := sha256.Sum256([]byte(keyStr))
	v := &Vault{
		path: path,
		key:  h[:],
		data: make(map[string]string),
	}
	if err := v.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return v, nil
}

func (v *Vault) Set(key, value string) error {
	v.data[key] = value
	return v.save()
}

func (v *Vault) Get(key string) string {
	return v.data[key]
}

func (v *Vault) save() error {
	plaintext, err := json.Marshal(v.data)
	if err != nil {
		return err
	}
	block, err := aes.NewCipher(v.key)
	if err != nil {
		return err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return os.WriteFile(v.path, ciphertext, 0600)
}

func (v *Vault) load() error {
	ciphertext, err := os.ReadFile(v.path)
	if err != nil {
		return err
	}
	block, err := aes.NewCipher(v.key)
	if err != nil {
		return err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return fmt.Errorf("ciphertext too short")
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return err
	}
	return json.Unmarshal(plaintext, &v.data)
}
