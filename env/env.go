package env

import (
	"github.com/spf13/viper"
)

// Environment environment
type Environment struct {
	Release           bool   `mapstructure:"RELEASE"`
	Production        bool   `mapstructure:"PRODUCTION"`
	ServerPort        string `mapstructure:"SERVER_PORT"`
	DatabaseURL       string `mapstructure:"DATA_BASE_URL"`
	DatabasePort      int    `mapstructure:"DATA_BASE_PORT"`
	DatabaseName      string `mapstructure:"DATA_BASE_NAME"`
	DatabaseRoot      bool   `mapstructure:"DATA_BASE_ROOT"`
	DatabaseUsername  string `mapstructure:"DATA_BASE_USERNAME"`
	DatabasePassword  string `mapstructure:"DATA_BASE_PASSWORD"`
	JWTSecret         string `mapstructure:"JWT_SECRET"`
	PartnerJWTSecret  string `mapstructure:"PARTNER_JWT_SECRET"`
	SwaggerHost       string `mapstructure:"SWAGGER_HOST"`
	SwaggerBasePath   string `mapstructure:"SWAGGER_BASE_PATH"`
	WebURL            string `mapstructure:"WEB_URL"`
	WebECommerceURL   string `mapstructure:"WEB_ECOMMERCE_URL"`
	CookieHashKey     string `mapstructure:"COOKIE_HASH_KEY"`
	CookieBlockKey    string `mapstructure:"COOKIE_BLOCK_KEY"`
	AntsSenderName    string `mapstructure:"ANTS_SENDER_NAME"`
	AntsAuth          string `mapstructure:"ANTS_AUTH"`
	OTPExpirationTime uint   `mapstructure:"OTP_EXPIRATION_TIME"`
	OTPIssuer         string `mapstructure:"OTP_ISSUER"`
	OTPAccountName    string `mapstructure:"OTP_ACCOUNT_NAME"`
	OrderTimeout      int64  `mapstructure:"ORDER_TIME_OUT"`
	RedisOn           bool   `mapstructure:"REDIS_ON"`
	RedisHost         string `mapstructure:"REDIS_HOST"`
	RedisPassword     string `mapstructure:"REDIS_PASSWORD"`
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
