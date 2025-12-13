package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	HTTPServer HTTPServer `mapstructure:"http_server"`
	URLDatabase   URLDatabase     `mapstructure:"url_database"`
	UserDatabase   UserDatabase     `mapstructure:"user_database"`
}

type HTTPServer struct {
	Host        string `mapstructure:"host"`
	Port        string    `mapstructure:"port"`
	Timeout     int    `mapstructure:"timeout"`
	IdleTimeout int    `mapstructure:"idle_timeout"`
}

type URLDatabase struct {
	Host string `mapstructure:"host"`
	Port string    `mapstructure:"port"`
	User string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName string `mapstructure:"dbname"`	
	SSL  string `mapstructure:"ssl"`
}

type UserDatabase struct {
	Host string `mapstructure:"host"`
	Port string    `mapstructure:"port"`
	User string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName string `mapstructure:"dbname"`	
	SSL  string `mapstructure:"ssl"`
}

func LoadConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./config/")	

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		panic(err)
	}
	godotenv.Load()
	config.URLDatabase.Password = os.Getenv("DB_PASSWORD")
	config.UserDatabase.Password = os.Getenv("DB_PASSWORD")

	return &config
}