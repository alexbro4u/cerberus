package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env      string `yaml:"env" env-default:"local"`
	GRPC     `yaml:"grpc_server"`
	Storage  `yaml:"storage"`
	TokenTTL time.Duration `yaml:"token_ttl" env-default:"1h"`
}

type GRPC struct {
	Host         string        `yaml:"host" env-default:"localhost"`
	Port         int           `yaml:"port" env-default:"8080"`
	ReadTimeout  time.Duration `yaml:"read_timeout" env-default:"4s"`
	WriteTimeout time.Duration `yaml:"write_timeout" env-default:"4s"`
}

type Storage struct {
	DbName   string `yaml:"db_name" env-required:"true"`
	User     string `yaml:"user" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	Host     string `yaml:"host" env-required:"true"`
}

func LoadConfig(configPath string, cfg interface{}) {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file not found: %s", configPath)
	}

	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	var cfg Config
	LoadConfig(configPath, &cfg)
	return &cfg
}

type Migrator struct {
	Storage         Storage `yaml:"storage"`
	MigrationsPath  string  `yaml:"migrations_path" env-required:"true"`
	MigrationsTable string  `yaml:"migrations_table" env-required:"true"`
}

func MustLoadMigrator() *Migrator {
	configPath := os.Getenv("CONFIG_PATH_MIGRATOR")
	if configPath == "" {
		log.Fatal("CONFIG_PATH_MIGRATOR is not set")
	}

	var migrator Migrator
	LoadConfig(configPath, &migrator)
	return &migrator
}
