package src

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	baseChainID    = 8453
	burnAddressHex = "0x0000000000000000000000000000000000000000"
)

func BuildPostPayload(title, link, description string, maxLen int) string {
	content := "\U0001F4F0 " + title
	linkSection := ""
	if link != "" {
		linkSection = "<br><br>\U0001F517 " + link
	}
	if description != "" {
		desc := sanitizeNonPrintable(description)
		remaining := maxLen - len(content) - len(linkSection) - 8
		if remaining > 0 && len(desc) > 0 {
			if len(desc) > remaining {
				desc = desc[:remaining] + "..."
			}
			content += "<br><br>" + desc
		}
	}
	content += linkSection
	if len(content) > maxLen {
		content = content[:maxLen-3] + "..."
	}
	payload := map[string]string{"p": content}
	payloadJSON, _ := json.Marshal(payload)
	return "yp/1/p:" + string(payloadJSON)
}
func SendPostTransaction(rpcUrl string, privateKey *ecdsa.PrivateKey, payload string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	client, err := ethclient.DialContext(ctx, rpcUrl)
	if err != nil {
		return "", fmt.Errorf("could not connect to RPC: %w", err)
	}
	defer client.Close()
	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return "", fmt.Errorf("could not get nonce: %w", err)
	}
	burnAddress := common.HexToAddress(burnAddressHex)
	data := []byte(payload)
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return "", fmt.Errorf("could not get gas price: %w", err)
	}
	msg := ethereum.CallMsg{
		From:  fromAddress,
		To:    &burnAddress,
		Value: big.NewInt(0),
		Data:  data,
	}
	gasLimit, err := client.EstimateGas(ctx, msg)
	if err != nil {
		return "", fmt.Errorf("could not estimate gas: %w", err)
	}
	gasLimit = gasLimit * 120 / 100
	tx := types.NewTransaction(nonce, burnAddress, big.NewInt(0), gasLimit, gasPrice, data)
	chainID := big.NewInt(baseChainID)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return "", fmt.Errorf("could not sign transaction: %w", err)
	}
	if err := client.SendTransaction(ctx, signedTx); err != nil {
		if strings.Contains(err.Error(), "insufficient funds") {
			return "", fmt.Errorf("insufficient ETH for gas fees - fund address %s on Base: %w", fromAddress.Hex(), err)
		}
		return "", fmt.Errorf("could not send transaction: %w", err)
	}
	return signedTx.Hash().Hex(), nil
}
