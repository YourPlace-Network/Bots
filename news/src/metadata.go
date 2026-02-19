package src

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"strings"
)

var validVerticals = []string{
	"anime", "art", "asmr", "automotive", "beauty", "books", "business",
	"career", "comedy", "cooking", "culture", "diy", "education",
	"entertainment", "fashion", "finance", "fitness", "food", "gaming",
	"health", "history", "lifestyle", "luxury", "movies", "music", "nature",
	"news", "other", "outdoors", "parenting", "pets", "photography",
	"politics", "realestate", "science", "spirituality", "sports",
	"sustainability", "technology", "television", "travel",
}

func isValidVertical(vertical string) bool {
	v := strings.ToLower(vertical)
	for _, valid := range validVerticals {
		if valid == v {
			return true
		}
	}
	return false
}
func SendMetadataAvatar(rpcUrl string, privateKey *ecdsa.PrivateKey, avatar string) (string, error) {
	avatar = sanitizeNonPrintable(avatar)
	payload, _ := json.Marshal(map[string]string{"a": avatar})
	return SendPostTransaction(rpcUrl, privateKey, "yp/1/ma:"+string(payload))
}
func SendMetadataBanner(rpcUrl string, privateKey *ecdsa.PrivateKey, banner string) (string, error) {
	banner = sanitizeNonPrintable(banner)
	payload, _ := json.Marshal(map[string]string{"b": banner})
	return SendPostTransaction(rpcUrl, privateKey, "yp/1/mb:"+string(payload))
}
func SendMetadataDescription(rpcUrl string, privateKey *ecdsa.PrivateKey, description string) (string, error) {
	description = sanitizeNonPrintable(description)
	payload, _ := json.Marshal(map[string]string{"d": description})
	return SendPostTransaction(rpcUrl, privateKey, "yp/1/md:"+string(payload))
}
func SendMetadataName(rpcUrl string, privateKey *ecdsa.PrivateKey, name string) (string, error) {
	name = sanitizeNonPrintable(name)
	payload, _ := json.Marshal(map[string]string{"n": name})
	return SendPostTransaction(rpcUrl, privateKey, "yp/1/mn:"+string(payload))
}
func SendMetadataVertical(rpcUrl string, privateKey *ecdsa.PrivateKey, vertical string) (string, error) {
	vertical = strings.ToLower(vertical)
	if !isValidVertical(vertical) {
		return "", fmt.Errorf("invalid vertical %q", vertical)
	}
	payload, _ := json.Marshal(map[string]string{"v": vertical})
	return SendPostTransaction(rpcUrl, privateKey, "yp/1/mv:"+string(payload))
}
