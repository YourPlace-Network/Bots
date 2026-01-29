package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/sha3"
)

type WalletFile struct {
	Address             string `json:"address"`
	EncryptedPrivateKey string `json:"encryptedPrivateKey"`
	PublicKey           string `json:"publicKey"`
}

func CreateWallet() (*ecdsa.PrivateKey, string, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, "", fmt.Errorf("could not generate private key: %w", err)
	}
	address := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
	return privateKey, address, nil
}
func LoadOrCreateWallet(path string, password string) (*ecdsa.PrivateKey, string, error) {
	if _, err := os.Stat(path); err == nil {
		return LoadWallet(path, password)
	}
	privateKey, address, err := CreateWallet()
	if err != nil {
		return nil, "", err
	}
	if err := SaveWallet(path, password, privateKey); err != nil {
		return nil, "", fmt.Errorf("could not save wallet: %w", err)
	}
	fmt.Println("New wallet created!")
	fmt.Println("Address:", address)
	fmt.Println("Fund this address with ETH on Base to pay for transaction fees.")
	fmt.Println("Wallet saved to:", path)
	return privateKey, address, nil
}
func LoadWallet(path string, password string) (*ecdsa.PrivateKey, string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, "", fmt.Errorf("could not read wallet file: %w", err)
	}
	var wf WalletFile
	if err := json.Unmarshal(data, &wf); err != nil {
		return nil, "", fmt.Errorf("could not parse wallet file: %w", err)
	}
	privateKeyHex, err := decryptString(password, wf.EncryptedPrivateKey)
	if err != nil {
		return nil, "", fmt.Errorf("could not decrypt private key: %w", err)
	}
	privateKey, err := crypto.HexToECDSA(privateKeyHex[2:])
	if err != nil {
		return nil, "", fmt.Errorf("could not parse private key: %w", err)
	}
	address := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
	return privateKey, address, nil
}
func SaveWallet(path string, password string, privateKey *ecdsa.PrivateKey) error {
	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateKeyHex := hexutil.Encode(privateKeyBytes)
	publicKeyBytes := crypto.FromECDSAPub(&privateKey.PublicKey)
	publicKeyHex := hexutil.Encode(publicKeyBytes)
	address := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
	encrypted, err := encryptString(password, privateKeyHex)
	if err != nil {
		return fmt.Errorf("could not encrypt private key: %w", err)
	}
	wf := WalletFile{
		Address:             address,
		EncryptedPrivateKey: encrypted,
		PublicKey:           publicKeyHex,
	}
	data, err := json.MarshalIndent(wf, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal wallet file: %w", err)
	}
	return os.WriteFile(path, data, 0600)
}

func deriveKey(password string) ([]byte, []byte, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return nil, nil, err
	}
	key := pbkdf2.Key([]byte(password), salt, 100000, 32, sha3.New512)
	return key, salt, nil
}
func deriveKeyWithSalt(password string, salt []byte) []byte {
	return pbkdf2.Key([]byte(password), salt, 100000, 32, sha3.New512)
}
func decryptString(password, cipherTextB64 string) (string, error) {
	cipherText, err := base64.StdEncoding.DecodeString(cipherTextB64)
	if err != nil {
		return "", fmt.Errorf("could not decode base64: %w", err)
	}
	nonceSize := 12
	saltSize := 16
	if len(cipherText) < saltSize+nonceSize+1 {
		return "", fmt.Errorf("ciphertext too short")
	}
	salt := cipherText[:saltSize]
	data := cipherText[saltSize : len(cipherText)-nonceSize]
	nonce := cipherText[len(cipherText)-nonceSize:]
	key := deriveKeyWithSalt(password, salt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	plainText, err := aesGCM.Open(nil, nonce, data, nil)
	if err != nil {
		return "", fmt.Errorf("decryption failed: %w", err)
	}
	return string(plainText), nil
}
func encryptString(password, plainText string) (string, error) {
	key, salt, err := deriveKey(password)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, 12)
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	cipherText := aesGCM.Seal(nil, nonce, []byte(plainText), nil)
	result := append(salt, cipherText...)
	result = append(result, nonce...)
	return base64.StdEncoding.EncodeToString(result), nil
}
