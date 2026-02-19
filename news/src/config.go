package src

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Avatar              string   `json:"avatar"`
	Banner              string   `json:"banner"`
	Description         string   `json:"description"`
	Feeds               []string `json:"feeds"`
	MaxPostLength       int      `json:"maxPostLength"`
	PollIntervalSeconds int      `json:"pollIntervalSeconds"`
	RpcUrl              string   `json:"rpcUrl"`
	Username            string   `json:"username"`
	Vertical            string   `json:"vertical"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read config file: %w", err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("could not parse config file: %w", err)
	}
	if len(cfg.Feeds) == 0 {
		return nil, fmt.Errorf("config must contain at least one feed URL")
	}
	if cfg.RpcUrl == "" {
		return nil, fmt.Errorf("config must contain a rpcUrl")
	}
	if cfg.MaxPostLength <= 0 {
		cfg.MaxPostLength = 500
	}
	if cfg.PollIntervalSeconds <= 0 {
		cfg.PollIntervalSeconds = 300
	}
	return &cfg, nil
}
