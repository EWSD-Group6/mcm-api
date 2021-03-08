package config

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

type Config struct {
	RedisAddr         string `mapstructure:"redis_addr"`
	RedisPassword     string `mapstructure:"redis_password"`
	RedisDb           int    `mapstructure:"redis_db"`
	RedisQueueName    string `mapstructure:"redis_queue_name"`
	S3MediaBucket     string `mapstructure:"s3_media_bucket"`
	DatabaseHost      string `mapstructure:"database_host"`
	DatabasePort      string `mapstructure:"database_port"`
	DatabaseUsername  string `mapstructure:"database_username"`
	DatabasePassword  string `mapstructure:"database_password"`
	DatabaseName      string `mapstructure:"database_name"`
	WebAppUrl         string `mapstructure:"web_app_url"`
	JwtSecret         string `mapstructure:"jwt_secret"`
	AdminEmail        string `mapstructure:"admin_email"`
	AdminPassword     string `mapstructure:"admin_password"`
	SesSenderEmail    string `mapstructure:"ses_sender_email"`
	MediaBucket       string `mapstructure:"media_bucket"`
	ConverterService  string `mapstructure:"converter_service"`
	ImageProxyService string `mapstructure:"image_proxy_service"`
}

func init() {
	_ = viper.BindEnv("redis_addr", strings.ToUpper("redis_addr"))
	_ = viper.BindEnv("redis_password", strings.ToUpper("redis_password"))
	_ = viper.BindEnv("redis_db", strings.ToUpper("redis_db"))
	_ = viper.BindEnv("redis_queue_name", strings.ToUpper("redis_queue_name"))
	_ = viper.BindEnv("s3_media_bucket", strings.ToUpper("s3_media_bucket"))
	_ = viper.BindEnv("database_host", strings.ToUpper("database_host"))
	_ = viper.BindEnv("database_port", strings.ToUpper("database_port"))
	_ = viper.BindEnv("database_username", strings.ToUpper("database_username"))
	_ = viper.BindEnv("database_password", strings.ToUpper("database_password"))
	_ = viper.BindEnv("database_name", strings.ToUpper("database_name"))
	_ = viper.BindEnv("web_app_url", strings.ToUpper("web_app_url"))
	_ = viper.BindEnv("jwt_secret", strings.ToUpper("jwt_secret"))
	_ = viper.BindEnv("admin_email", strings.ToUpper("admin_email"))
	_ = viper.BindEnv("admin_password", strings.ToUpper("admin_password"))
	_ = viper.BindEnv("ses_sender_email", strings.ToUpper("ses_sender_email"))
	_ = viper.BindEnv("media_bucket", strings.ToUpper("media_bucket"))
	_ = viper.BindEnv("converter_service", strings.ToUpper("converter_service"))
	_ = viper.BindEnv("image_proxy_service", strings.ToUpper("image_proxy_service"))
}

func (config *Config) GetDatabaseDsn() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		config.DatabaseHost,
		config.DatabaseUsername,
		config.DatabasePassword,
		config.DatabaseName,
		config.DatabasePort,
	)
}

func (config *Config) GetDatabaseUrl() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.DatabaseUsername,
		config.DatabasePassword,
		config.DatabaseHost,
		config.DatabasePort,
		config.DatabaseName,
	)
}
