package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config 是应用配置的根结构
// 可通过不同环境的 YAML 文件提供不同值
// 示例见 config/dev.yaml
type Config struct {
	Server struct {
		Addr string `yaml:"addr"` // 例如 :8080
	} `yaml:"server"`

	Database struct {
		DSN         string `yaml:"dsn"`          // 为空则使用内存仓库
		AutoMigrate bool   `yaml:"auto_migrate"` // MySQL 模式下是否自动迁移
	} `yaml:"database"`

	Kafka struct {
		Brokers []string `yaml:"brokers"` // 为空则使用内存事件总线
		Topic   string   `yaml:"topic"`
		GroupID string   `yaml:"group"`
	} `yaml:"kafka"`

	Redis struct {
		Addr     string `yaml:"addr"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	} `yaml:"redis"`
}

// Load 从给定的 YAML 文件路径加载配置
func Load(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file failed: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal yaml failed: %w", err)
	}
	// 默认值填充
	if cfg.Server.Addr == "" {
		cfg.Server.Addr = ":8080"
	}
	return &cfg, nil
}
