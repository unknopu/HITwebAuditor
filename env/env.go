package env

import (
	"github.com/spf13/viper"
)

// Environment environment
type Environment struct {
	Release          bool   `mapstructure:"RELEASE"`
	Production       bool   `mapstructure:"PRODUCTION"`
	ServerPort       string `mapstructure:"SERVER_PORT"`
	DatabaseURL      string `mapstructure:"DATA_BASE_URL"`
	DatabasePort     int    `mapstructure:"DATA_BASE_PORT"`
	DatabaseName     string `mapstructure:"DATA_BASE_NAME"`
	DatabaseRoot     bool   `mapstructure:"DATA_BASE_ROOT"`
	DatabaseUsername string `mapstructure:"DATA_BASE_USERNAME"`
	DatabasePassword string `mapstructure:"DATA_BASE_PASSWORD"`
	RedisOn          bool   `mapstructure:"REDIS_ON"`
	RedisHost        string `mapstructure:"REDIS_HOST"`
	RedisPassword    string `mapstructure:"REDIS_PASSWORD"`
}

// Read init env
func Read(path string) (*Environment, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.AddConfigPath(path)
	v.AutomaticEnv()
	v.SetConfigType("yml")
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	env := &Environment{}
	err := v.Unmarshal(&env)
	if err != nil {
		return nil, err
	}
	return env, nil
}
