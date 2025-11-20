package config

import (
	"log"
	"log/slog"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/joho/godotenv"

	"github.com/spf13/viper"
)

type Config struct {
	AppEnv string
	Port   string

	JWTSecret     string
	RefreshSecret string
	JWTIssuer     string
	JWTAudience   string

	APIRateLimit int
	Timeout      time.Duration

	DatabaseConfig *DatabaseConfig
	EmailConfig    *EmailConfig
	RedisConfig    *RedisConfig
	FirebaseConfig *FirebaseConfig
	TwilioConfig   *TwilioConfig
	AWSConfig      *AWSConfig
	FrontendURL    string
}

type TwilioConfig struct {
	MessageServiceSID string
	AccountSID        string
	AuthToken         string
	FromNumber        string
}

type FirebaseConfig struct {
	ProjectID   string
	ClientEmail string
	PrivateKey  string
}

type DatabaseConfig struct {
	Host       string
	Port       string
	User       string
	Password   string
	Name       string
	SslMode    string
	ConnString string
}

type EmailConfig struct {
	Port string
	Host string
	User string
	Pass string
	From string
	Ssl  bool
	Tls  bool
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type AWSConfig struct {
	Region          string
	Bucket          string
	AccessKeyID     string
	SecretAccessKey string
	Endpoint        string
	AwsConfig       aws.Config
}

func NewDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:       viper.GetString("DB_HOST"),
		Port:       viper.GetString("DB_PORT"),
		User:       viper.GetString("DB_USER"),
		Password:   viper.GetString("DB_PASSWORD"),
		Name:       viper.GetString("DB_NAME"),
		SslMode:    viper.GetString("DB_SSLMODE"),
		ConnString: viper.GetString("DB_CONN_STRING"),
	}
}

func NewEmailConfig() *EmailConfig {
	return &EmailConfig{
		Port: viper.GetString("EMAIL_PORT"),
		Host: viper.GetString("EMAIL_HOST"),
		User: viper.GetString("EMAIL_USER"),
		Pass: viper.GetString("EMAIL_PASSWORD"),
		From: viper.GetString("EMAIL_FROM"),
		Ssl:  viper.GetBool("EMAIL_SSL"),
		Tls:  viper.GetBool("EMAIL_TLS"),
	}
}

func NewFirebaseConfig() *FirebaseConfig {
	return &FirebaseConfig{
		ProjectID:   viper.GetString("FIREBASE_PROJECT_ID"),
		ClientEmail: viper.GetString("FIREBASE_CLIENT_EMAIL"),
		PrivateKey:  viper.GetString("FIREBASE_PRIVATE_KEY"),
	}
}

func NewTwilioConfig() *TwilioConfig {
	return &TwilioConfig{
		MessageServiceSID: viper.GetString("TWILIO_MESSAGE_SERVICE_SID"),
		AccountSID:        viper.GetString("TWILIO_ACCOUNT_SID"),
		AuthToken:         viper.GetString("TWILIO_AUTH_TOKEN"),
		FromNumber:        viper.GetString("TWILIO_FROM_NUMBER"),
	}
}

func NewRedisConfig() *RedisConfig {
	db := viper.GetInt("REDIS_DB")
	if db < 0 {
		db = 0
	}
	return &RedisConfig{
		Host:     viper.GetString("REDIS_HOST"),
		Port:     viper.GetInt("REDIS_PORT"),
		Password: viper.GetString("REDIS_PASSWORD"),
		DB:       db,
	}
}

func NewAWSConfig() *AWSConfig {
	return &AWSConfig{
		Region:          viper.GetString("AWS_REGION"),
		Bucket:          viper.GetString("AWS_S3_BUCKET"),
		AccessKeyID:     viper.GetString("ACCESS_KEY_ID"),
		SecretAccessKey: viper.GetString("SECRET_ACCESS_KEY"),
		Endpoint:        viper.GetString("AWS_S3_ENDPOINT"),
	}
}

func (c *Config) IsDevelopment() bool {
	return c.AppEnv == "development" || c.AppEnv == "dev"
}

func Load() Config {
	_ = godotenv.Load(".env")

	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	err := viper.ReadInConfig()
	if err != nil {
		log.Println("⚠️  .env file not found, using environment variables")
	}

	if !viper.IsSet("FRONTEND_URL") {
		viper.Set("FRONTEND_URL", "http://localhost:3000")
	}

	cfg := Config{
		AppEnv:         viper.GetString("APP_ENV"),
		Port:           viper.GetString("PORT"),
		DatabaseConfig: NewDatabaseConfig(),
		EmailConfig:    NewEmailConfig(),
		RedisConfig:    NewRedisConfig(),
		FirebaseConfig: NewFirebaseConfig(),
		TwilioConfig:   NewTwilioConfig(),
		AWSConfig:      NewAWSConfig(),

		FrontendURL: viper.GetString("FRONTEND_URL"),

		JWTSecret:    viper.GetString("JWT_SECRET"),
		JWTIssuer:    viper.GetString("JWT_ISSUER"),
		JWTAudience:  viper.GetString("JWT_AUDIENCE"),
		APIRateLimit: viper.GetInt("API_RATE_LIMIT"),
		Timeout:      viper.GetDuration("TIMEOUT"),
	}

	validateConfig(cfg)

	return cfg
}

func validateConfig(cfg Config) {
	var missing []string

	if cfg.Port == "" {
		missing = append(missing, "PORT")
	}
	if cfg.JWTSecret == "" {
		missing = append(missing, "JWT_SECRET")
	}
	if cfg.AppEnv == "" {
		missing = append(missing, "APP_ENV")
	}
	if cfg.DatabaseConfig.Host == "" {
		missing = append(missing, "DB_HOST")
	}
	if cfg.DatabaseConfig.Port == "" {
		missing = append(missing, "DB_PORT")
	}
	if cfg.DatabaseConfig.User == "" {
		missing = append(missing, "DB_USER")
	}
	if cfg.DatabaseConfig.Password == "" {
		missing = append(missing, "DB_PASSWORD")
	}
	if cfg.DatabaseConfig.Name == "" {
		missing = append(missing, "DB_NAME")
	}
	if cfg.EmailConfig.Host == "" {
		missing = append(missing, "EMAIL_HOST")
	}
	if cfg.EmailConfig.Port == "" {
		missing = append(missing, "EMAIL_PORT")
	}
	if cfg.EmailConfig.User == "" {
		missing = append(missing, "EMAIL_USER")
	}
	if cfg.EmailConfig.Pass == "" {
		missing = append(missing, "EMAIL_PASSWORD")
	}
	if cfg.EmailConfig.From == "" {
		missing = append(missing, "EMAIL_FROM")
	}

	if cfg.RedisConfig.Host == "" {
		missing = append(missing, "REDIS_HOST")
	}
	if cfg.RedisConfig.Port == 0 {
		missing = append(missing, "REDIS_PORT")
	}
	if cfg.FirebaseConfig.ProjectID == "" {
		missing = append(missing, "FIREBASE_PROJECT_ID")
	}
	if cfg.FirebaseConfig.ClientEmail == "" {
		missing = append(missing, "FIREBASE_CLIENT_EMAIL")
	}
	if cfg.FirebaseConfig.PrivateKey == "" {
		missing = append(missing, "FIREBASE_PRIVATE_KEY")
	}

	if cfg.TwilioConfig.MessageServiceSID == "" {
		missing = append(missing, "TWILIO_MESSAGE_SERVICE_SID")
	}
	if cfg.TwilioConfig.AccountSID == "" {
		missing = append(missing, "TWILIO_ACCOUNT_SID")
	}
	if cfg.TwilioConfig.AuthToken == "" {
		missing = append(missing, "TWILIO_AUTH_TOKEN")
	}

	if len(missing) > 0 {
		slog.Error("❌ Missing required variables", slog.Any("missing", missing))
		log.Fatalf("❌ Missing required variables: %s", strings.Join(missing, ", "))
	}
}
