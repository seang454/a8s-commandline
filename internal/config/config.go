package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

func Load(cfgFile string) {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err == nil {
			viper.AddConfigPath(home)
		}
		viper.AddConfigPath(".")
		viper.SetConfigName(".a8s")
		viper.SetConfigType("yaml")
	}
	viper.SetEnvPrefix("A8S")
	viper.AutomaticEnv()
	viper.SetDefault("api_url", "http://localhost:8080")
	viper.SetDefault("api_token", "")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Fprintln(os.Stderr, "Warning: could not read config:", err)
		}
	}
}
