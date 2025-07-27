package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DBSource  string `mapstructure:"DB_SOURCE"`
	JWTSecret string `mapstructure:"JWT_SECRET"`
	Port      string `mapstructure:"PORT"`
}

var AppConfig Config

func LoadConfig(path string) error {
	viper.SetConfigName(".env") // 🔁 now uses .env
	viper.SetConfigType("env")
	viper.AddConfigPath(path) // look for config in the path
	viper.AutomaticEnv()      // override with ENV vars

	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("No config file found: %v\n", err)
		// still continue to allow env vars only
	}

	err = viper.Unmarshal(&AppConfig)
	if err != nil {
		return err
	}
	return nil
}
