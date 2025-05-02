package config

import (
	"encoding/json"
	"os"
)

type CategorySource struct {
	URL      string `json:"url"`
	Category string `json:"category"`
}

type Config struct {
	DBStorePath string           `json:"dbstore_path"`
	Categories  []CategorySource `json:"categories"`
}

func LoadConfig(filename string) (*Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config //[]CategorySource
	err = json.NewDecoder(f).Decode(&cfg)
	return &cfg, err
}
