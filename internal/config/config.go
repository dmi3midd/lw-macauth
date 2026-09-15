package config

import (
	"crypto/rsa"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.yaml.in/yaml/v3"
)

type Server struct {
	Address      string        `yaml:"address"`
	IdleTimeout  time.Duration `yaml:"idleTimeout"`
	ReadTimeout  time.Duration `yaml:"readTimeout"`
	WriteTimeout time.Duration `yaml:"writeTimeout"`
}

type PEM struct {
	PrivPath string `yaml:"privPath"`
	PubPath  string `yaml:"pubPath"`
}

type Log struct {
	Level string `yaml:"level"`
}

type Config struct {
	HTTPServer Server    `yaml:"server"`
	PEM        PEM       `yaml:"pem"`
	Log        Log       `yaml:"log"`
	Keys       *KeysPair `yaml:"-"`
}

func LoadConfig() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./config.yaml"
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	keys, err := LoadKeys(cfg.PEM.PrivPath, cfg.PEM.PubPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load keys: %w", err)
	}
	cfg.Keys = keys

	return cfg, nil
}

type KeysPair struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

func LoadKeys(privPath, pubPath string) (*KeysPair, error) {
	privKeyData, err := os.ReadFile(privPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key from %s: %w", privPath, err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privKeyData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	pubKeyData, err := os.ReadFile(pubPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key from %s: %w", pubPath, err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(pubKeyData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	return &KeysPair{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}, nil
}
