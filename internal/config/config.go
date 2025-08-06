package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

const (
	ConfigPath = "./config.yml"
)

type ServerConfig struct {
	HttpPort     string `yaml:"http_port"`
	ReadTimeout  int    `yaml:"read_timeout"`
	WriteTimeout int    `yaml:"write_timeout"`
}

type StageConfig struct {
	IsDev       bool   `yaml:"is_dev"`
	LogFilePath string `yaml:"log_file_path"`
	Stage       string `yaml:"stage"`
}
type MetricsConfig struct {
	Enabled bool `yaml:"enabled"`
	Port    int  `yaml:"port"`
}

type RuleConfig struct {
	From            string   `yaml:"from"`
	RedirectTo      string   `yaml:"redirectTo"`
	AuthRequired    bool     `yaml:"auth_required"`
	RolesRequired   []string `yaml:"roles_required"`
	AllowedHeaders  []string `yaml:"allowed_headers"`
	HeadersRequired []string `yaml:"headers_required"`
}

type Config struct {
	Server  ServerConfig  `yaml:"server"`
	IDP     IDPConfig     `yaml:"idp"`
	Stage   StageConfig   `yaml:"stage"`
	Rules   []RuleConfig  `yaml:"rules"`
	Metrics MetricsConfig `yaml:"metrics"`
	APIKeys APIKeysConfig `yaml:"api_keys"`
}
type APIKeysConfig struct {
	Source   string   `yaml:"source"`
	Keys     []string `yaml:"keys"`
	FilePath string   `yaml:"file_path"`
}
type IDPConfig struct {
	RegistrationURL string `yaml:"registration_url"`
}

func Load() (*Config, error) {
	cfg, err := os.ReadFile(ConfigPath)
	if err != nil {
		return nil, err
	}

	var c Config
	err = yaml.Unmarshal(cfg, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
