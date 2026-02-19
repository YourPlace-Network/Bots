package src

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type WalletFile struct {
	Address    string `json:"address"`
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
}

func CreateWallet() (*ecdsa.PrivateKey, string, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, "", fmt.Errorf("could not generate private key: %w", err)
	}
	address := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
	return privateKey, address, nil
}
func LoadOrCreateWallet(path string) (*ecdsa.PrivateKey, string, error) {
	if _, err := os.Stat(path); err == nil {
		return LoadWallet(path)
	}
	privateKey, address, err := CreateWallet()
	if err != nil {
		return nil, "", err
	}
	if err := SaveWallet(path, privateKey); err != nil {
		return nil, "", fmt.Errorf("could not save wallet: %w", err)
	}
	fmt.Println("New wallet created!")
	fmt.Println("Address:", address)
	fmt.Println("Fund this address with ETH on Base to pay for transaction fees.")
	fmt.Println("Wallet saved to:", path)
	return privateKey, address, nil
}
func LoadWallet(path string) (*ecdsa.PrivateKey, string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, "", fmt.Errorf("could not read wallet file: %w", err)
	}
	var wf WalletFile
	if err := json.Unmarshal(data, &wf); err != nil {
		return nil, "", fmt.Errorf("could not parse wallet file: %w", err)
	}
	privateKey, err := crypto.HexToECDSA(wf.PrivateKey[2:])
	if err != nil {
		return nil, "", fmt.Errorf("could not parse private key: %w", err)
	}
	address := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
	return privateKey, address, nil
}
func SaveWallet(path string, privateKey *ecdsa.PrivateKey) error {
	privateKeyHex := hexutil.Encode(crypto.FromECDSA(privateKey))
	publicKeyHex := hexutil.Encode(crypto.FromECDSAPub(&privateKey.PublicKey))
	address := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
	wf := WalletFile{
		Address:    address,
		PrivateKey: privateKeyHex,
		PublicKey:  publicKeyHex,
	}
	data, err := json.MarshalIndent(wf, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal wallet file: %w", err)
	}
	return os.WriteFile(path, data, 0600)
}
