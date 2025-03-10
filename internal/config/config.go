package config

import (
	"fmt"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/spf13/viper"
)

type Config struct {
	Env     string `yaml:"env" env-default:"local"`
	Storage `yaml:"storage" env-required:"true"`
	// HTTPServer  `yaml:"http_server"`
}

type Storage struct {
	Host     string `yaml:"host" env-default:"go_db"`
	Port     string `yaml:"port" env-default:"5432"`
	User     string `yaml:"user" env-default:"postgres"`
	Password string `yaml:"password" env-default:"root"`
	DBname   string `yaml:"dbname" env-default:"postgres"`
	SSLMode  string `yaml:"sslmode" env-default:"disable"`
}

func MustLoad() Config {
	// if err := initConfig(); err != nil {
	// 	log.Fatalf("error initializing configs: %s", err.Error())
	// }

	// configPath := os.Getenv("CONFIG_PATH")

	configPath := "./config/prod.yaml"
	entries, err := os.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}
	for _, e := range entries {
		fmt.Println(e.Name())
	}

	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	// fmt.Println(viper.GetString("storage.db"))
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return cfg
}


func initConfig() error {
	viper.AddConfigPath("config")
	viper.SetConfigName("prod")
	viper.SetConfigType("yaml")
	return viper.ReadInConfig()
}