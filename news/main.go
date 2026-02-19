package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"yourplace-news-bot/src"
)

func main() {
	configPath := flag.String("config", "config.json", "path to config.json")
	dataDir := flag.String("data", "data", "path to data directory")
	metadata := flag.Bool("metadata", false, "set profile metadata onchain and exit")
	once := flag.Bool("once", false, "run once and exit")
	single := flag.Bool("single", false, "post only the latest article from each feed, then exit")
	flag.Parse()
	cfg, err := src.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}
	if err := os.MkdirAll(*dataDir, 0755); err != nil {
		log.Fatalf("Could not create data directory: %v", err)
	}
	walletPath := filepath.Join("wallet.json")
	privateKey, address, err := src.LoadOrCreateWallet(walletPath)
	if err != nil {
		log.Fatalf("Wallet error: %v", err)
	}
	fmt.Printf("Bot wallet address: %s\n", address)
	if *metadata {
		if cfg.Avatar != "" {
			txHash, err := src.SendMetadataAvatar(cfg.RpcUrl, privateKey, cfg.Avatar)
			if err != nil {
				fmt.Printf("Error setting avatar: %v\n", err)
			} else {
				fmt.Printf("Avatar set: %s\n", txHash)
			}
			time.Sleep(5 * time.Second)
		}
		if cfg.Banner != "" {
			txHash, err := src.SendMetadataBanner(cfg.RpcUrl, privateKey, cfg.Banner)
			if err != nil {
				fmt.Printf("Error setting banner: %v\n", err)
			} else {
				fmt.Printf("Banner set: %s\n", txHash)
			}
			time.Sleep(5 * time.Second)
		}
		if cfg.Description != "" {
			txHash, err := src.SendMetadataDescription(cfg.RpcUrl, privateKey, cfg.Description)
			if err != nil {
				fmt.Printf("Error setting description: %v\n", err)
			} else {
				fmt.Printf("Description set: %s\n", txHash)
			}
			time.Sleep(5 * time.Second)
		}
		if cfg.Username != "" {
			txHash, err := src.SendMetadataName(cfg.RpcUrl, privateKey, cfg.Username)
			if err != nil {
				fmt.Printf("Error setting username: %v\n", err)
			} else {
				fmt.Printf("Username set: %s\n", txHash)
			}
			time.Sleep(5 * time.Second)
		}
		if cfg.Vertical != "" {
			txHash, err := src.SendMetadataVertical(cfg.RpcUrl, privateKey, cfg.Vertical)
			if err != nil {
				fmt.Printf("Error setting vertical: %v\n", err)
			} else {
				fmt.Printf("Vertical set: %s\n", txHash)
			}
		}
		fmt.Println("Metadata complete. Exiting.")
		return
	}
	dbPath := filepath.Join(*dataDir, "news.db")
	dedup, err := src.OpenDedupDB(dbPath)
	if err != nil {
		log.Fatalf("Database error: %v", err)
	}
	defer dedup.Close()
	for {
		for _, feedUrl := range cfg.Feeds {
			fmt.Printf("Fetching feed: %s\n", feedUrl)
			items, err := src.FetchFeed(feedUrl)
			if err != nil {
				fmt.Printf("Error fetching feed %s: %v\n", feedUrl, err)
				continue
			}
			fmt.Printf("Found %d items in feed\n", len(items))
			if *single {
				for _, item := range items {
					if item.GUID == "" {
						continue
					}
					posted, err := dedup.IsPosted(item.GUID)
					if err != nil {
						fmt.Printf("Error checking dedup for %s: %v\n", item.GUID, err)
						break
					}
					if posted {
						continue
					}
					payload := src.BuildPostPayload(item.Title, item.Link, item.Description, cfg.MaxPostLength)
					fmt.Printf("Posting: %s\n", item.Title)
					txHash, err := src.SendPostTransaction(cfg.RpcUrl, privateKey, payload)
					if err != nil {
						fmt.Printf("Error sending transaction: %v\n", err)
					} else {
						fmt.Printf("Transaction sent: %s\n", txHash)
						if err := dedup.MarkPosted(item.GUID, feedUrl, item.Title, txHash); err != nil {
							fmt.Printf("Error marking as posted: %v\n", err)
						}
					}
					break
				}
			} else {
				for i := len(items) - 1; i >= 0; i-- {
					item := items[i]
					if item.GUID == "" {
						continue
					}
					posted, err := dedup.IsPosted(item.GUID)
					if err != nil {
						fmt.Printf("Error checking dedup for %s: %v\n", item.GUID, err)
						continue
					}
					if posted {
						continue
					}
					payload := src.BuildPostPayload(item.Title, item.Link, item.Description, cfg.MaxPostLength)
					fmt.Printf("Posting: %s\n", item.Title)
					txHash, err := src.SendPostTransaction(cfg.RpcUrl, privateKey, payload)
					if err != nil {
						fmt.Printf("Error sending transaction: %v\n", err)
						continue
					}
					fmt.Printf("Transaction sent: %s\n", txHash)
					if err := dedup.MarkPosted(item.GUID, feedUrl, item.Title, txHash); err != nil {
						fmt.Printf("Error marking as posted: %v\n", err)
					}
					time.Sleep(5 * time.Second)
				}
			}
		}
		if err := dedup.CleanOld(30); err != nil {
			fmt.Printf("Error cleaning old entries: %v\n", err)
		}
		if *once || *single {
			fmt.Println("Single run complete. Exiting.")
			return
		}
		fmt.Printf("Sleeping for %d seconds...\n", cfg.PollIntervalSeconds)
		time.Sleep(time.Duration(cfg.PollIntervalSeconds) * time.Second)
	}
}
