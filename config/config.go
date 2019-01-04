package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Host    string
	Key     string
	Secret  string
	Timeout int64
	DB      struct {
		Host     string
		Login    string
		Password string
		Name     string
	}
	Neural struct {
		Iterations int
		Predict    float64
	}
	Strategy struct {
		Profit   float64
		StopLose float64
		Quantity float32
	}
}

type MasterConfig struct {
	IsDev  bool
	Master Config
	Dev    Config
}

func LoadConfig(path string) (Config, error) {
	config, err := LoadMasterConfig(path)
	if err != nil {
		return Config{}, err
	}
	if config.IsDev {
		return config.Dev, nil
	}

	return config.Master, nil
}

func LoadMasterConfig(path string) (MasterConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return MasterConfig{}, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	config := MasterConfig{}
	err = decoder.Decode(&config)
	if err != nil {
		return MasterConfig{}, err
	}
	return config, nil
}
