package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver           string `mapstructure:"DB_DRIVER"`
	DBUser             string `mapstructure:"DB_USER"`
	DBPass             string `mapstructure:"DB_PASS"`
	DBName             string `mapstructure:"DB_NAME"`
	DBHost             string `mapstructure:"DB_HOST"`
	DBPort             int    `mapstructure:"DB_PORT"`
	RedisHost          string `mapstructure:"REDIS_HOST"`
	RedisPort          int    `mapstructure:"REDIS_PORT"`
	RedisDB            int    `mapstructure:"REDIS_DB"`
	APPPort            int    `mapstructure:"APP_PORT"`
	JwtSecretKey       string `mapstructure:"JWT_SECRET_KEY"`
	JwtRefreshKey      string `mapstructure:"JWT_REFRESH_KEY"`
	DatabaseSchemaUser string `mapstructure:"DATABASE_SCHEMA"`
	WhitelistHost      string `mapstructure:"WHITELISTHOST"`
}

func load() *Config {
	// viper.SetConfigType("env")
	viper.AddConfigPath("./.development")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	// viper.AddConfigPath("../../.development")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("error reading config file ", err)
		return &Config{}
	}

	return &Config{
		DBDriver:           viper.GetString("DB_DRIVER"),
		DBUser:             viper.GetString("DB_USER"),
		DBPass:             viper.GetString("DB_PASS"),
		DBName:             viper.GetString("DB_NAME"),
		DBHost:             viper.GetString("DB_HOST"),
		DBPort:             viper.GetInt("DB_PORT"),
		RedisHost:          viper.GetString("REDIS_HOST"),
		RedisPort:          viper.GetInt("REDIS_PORT"),
		RedisDB:            viper.GetInt("REDIS_DB"),
		APPPort:            viper.GetInt("APP_PORT"),
		JwtSecretKey:       viper.GetString("JWT_SECRET_KEY"),
		DatabaseSchemaUser: viper.GetString("DATABASE_SCHEMA"),
		WhitelistHost:      viper.GetString("WHITELISTHOST"),
	}
}

var config = load()

func Cfg() *Config {
	return config
}
