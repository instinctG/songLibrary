package setting

import (
	"github.com/spf13/viper"
	"log"
)

type Server struct {
	Port string `mapstructure:"HTTP_PORT"`
}

var ServerSetting = &Server{}

type Database struct {
	Host     string `mapstructure:"DB_HOST"`
	Port     string `mapstructure:"DB_PORT"`
	User     string `mapstructure:"DB_USER"`
	Name     string `mapstructure:"DB_NAME"`
	Password string `mapstructure:"DB_PASSWORD"`
	SSLMode  string `mapstructure:"SSL_MODE"`
}

var DatabaseSetting = &Database{}

// Setup initialize the configuration instance
func Setup() {
	viper.AddConfigPath("conf")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	if err := viper.Unmarshal(&ServerSetting); err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	if err := viper.Unmarshal(&DatabaseSetting); err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

}
