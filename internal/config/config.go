package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	SMS      SMSConfig      `mapstructure:"sms"`
	COS      COSConfig      `mapstructure:"cos"`
	IM       IMConfig       `mapstructure:"im"`
	Maps     MapsConfig     `mapstructure:"maps"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"sslmode"`
}

func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		d.User, d.Password, d.Host, d.Port, d.Name, d.SSLMode,
	)
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type JWTConfig struct {
	Secret      string `mapstructure:"secret"`
	ExpireHours int    `mapstructure:"expire_hours"`
}

func (j JWTConfig) ExpireDuration() time.Duration {
	return time.Duration(j.ExpireHours) * time.Hour
}

type SMSConfig struct {
	SecretID   string `mapstructure:"secret_id"`
	SecretKey  string `mapstructure:"secret_key"`
	SDKAppID   string `mapstructure:"sdk_app_id"`
	TemplateID string `mapstructure:"template_id"`
	SignName   string `mapstructure:"sign_name"`
}

type COSConfig struct {
	SecretID  string `mapstructure:"secret_id"`
	SecretKey string `mapstructure:"secret_key"`
	BucketURL string `mapstructure:"bucket_url"`
	Region    string `mapstructure:"region"`
}

type IMConfig struct {
	AppID  string `mapstructure:"app_id"`
	Secret string `mapstructure:"secret"`
}

type MapsConfig struct {
	APIKey string `mapstructure:"api_key"`
}

func Load(configPath string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	v.AutomaticEnv()
	v.BindEnv("database.host", "DB_HOST")
	v.BindEnv("database.port", "DB_PORT")
	v.BindEnv("database.user", "DB_USER")
	v.BindEnv("database.password", "DB_PASSWORD")
	v.BindEnv("database.name", "DB_NAME")
	v.BindEnv("redis.addr", "REDIS_ADDR")
	v.BindEnv("redis.password", "REDIS_PASSWORD")
	v.BindEnv("jwt.secret", "JWT_SECRET")
	v.BindEnv("jwt.expire_hours", "JWT_EXPIRE_HOURS")
	v.BindEnv("sms.secret_id", "SMS_SECRET_ID")
	v.BindEnv("sms.secret_key", "SMS_SECRET_KEY")
	v.BindEnv("sms.sdk_app_id", "SMS_SDK_APP_ID")
	v.BindEnv("sms.template_id", "SMS_TEMPLATE_ID")
	v.BindEnv("cos.secret_id", "COS_SECRET_ID")
	v.BindEnv("cos.secret_key", "COS_SECRET_KEY")
	v.BindEnv("cos.bucket_url", "COS_BUCKET_URL")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &cfg, nil
}
