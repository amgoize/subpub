package config

import (
	"encoding/json"
	"os"
	"time"
)

type Config struct {
	GRPC struct {
		Port                    int           `json:"port"`
		ShutdownTimeout         string        `json:"shutdown_timeout"` // читается из JSON
		ShutdownTimeoutDuration time.Duration `json:"-"`                // НЕ читается из JSON, задается вручную
	} `json:"grpc"`

	Log struct {
		Level string `json:"level"`
	} `json:"log"`
}

func LoadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	// Преобразуем строку в time.Duration и сохраняем в отдельное поле
	duration, err := time.ParseDuration(config.GRPC.ShutdownTimeout)
	if err != nil {
		return nil, err
	}
	config.GRPC.ShutdownTimeoutDuration = duration

	return &config, nil
}
