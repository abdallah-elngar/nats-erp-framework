package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config يمثل جميع إعدادات النظام
type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Cache    CacheConfig    `mapstructure:"cache"`
	Auth     AuthConfig     `mapstructure:"auth"`
	Session  SessionConfig  `mapstructure:"session"`
	Log      LogConfig      `mapstructure:"log"`
	Mail     MailConfig     `mapstructure:"mail"`
	Storage  StorageConfig  `mapstructure:"storage"`
}

// AppConfig إعدادات التطبيق
type AppConfig struct {
	Name          string `mapstructure:"name"`
	Env           string `mapstructure:"env"`
	Debug         bool   `mapstructure:"debug"`
	DeveloperMode bool   `mapstructure:"developer_mode"`
	URL           string `mapstructure:"url"`
	Timezone      string `mapstructure:"timezone"`
	Locale        string `mapstructure:"locale"`
}

// ServerConfig إعدادات الخادم
type ServerConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	ReadTimeout  string `mapstructure:"read_timeout"`
	WriteTimeout string `mapstructure:"write_timeout"`
	IdleTimeout  string `mapstructure:"idle_timeout"`
}

// DatabaseConfig إعدادات قاعدة البيانات
type DatabaseConfig struct {
	Default     string                      `mapstructure:"default"`
	Connections map[string]ConnectionConfig `mapstructure:"connections"`
}

// ConnectionConfig إعدادات اتصال قاعدة البيانات
type ConnectionConfig struct {
	Driver   string            `mapstructure:"driver"`
	Host     string            `mapstructure:"host"`
	Port     int               `mapstructure:"port"`
	Database string            `mapstructure:"database"`
	Username string            `mapstructure:"username"`
	Password string            `mapstructure:"password"`
	Options  map[string]string `mapstructure:"options"`
}

// CacheConfig إعدادات التخزين المؤقت
type CacheConfig struct {
	Default     string                     `mapstructure:"default"`
	Connections map[string]CacheConnection `mapstructure:"connections"`
}

// CacheConnection إعدادات اتصال التخزين المؤقت
type CacheConnection struct {
	Driver   string `mapstructure:"driver"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Database int    `mapstructure:"db"`
}

// AuthConfig إعدادات المصادقة
type AuthConfig struct {
	JWT JWTConfig `mapstructure:"jwt"`
}

// JWTConfig إعدادات JWT
type JWTConfig struct {
	Secret            string `mapstructure:"secret"`
	Expiration        string `mapstructure:"expiration"`
	RefreshExpiration string `mapstructure:"refresh_expiration"`
}

// SessionConfig إعدادات الجلسات
type SessionConfig struct {
	Driver   string `mapstructure:"driver"`
	Lifetime int    `mapstructure:"lifetime"`
	Secure   bool   `mapstructure:"secure"`
}

// LogConfig إعدادات التسجيل
type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

// MailConfig إعدادات البريد
type MailConfig struct {
	Host       string `mapstructure:"host"`
	Port       int    `mapstructure:"port"`
	Username   string `mapstructure:"username"`
	Password   string `mapstructure:"password"`
	Encryption string `mapstructure:"encryption"`
}

// StorageConfig إعدادات التخزين
type StorageConfig struct {
	Disk  string                `mapstructure:"disk"`
	Disks map[string]DiskConfig `mapstructure:"disks"`
}

// DiskConfig إعدادات وسيط التخزين
type DiskConfig struct {
	Root   string `mapstructure:"root"`
	URL    string `mapstructure:"url"`
	Key    string `mapstructure:"key"`
	Secret string `mapstructure:"secret"`
	Region string `mapstructure:"region"`
	Bucket string `mapstructure:"bucket"`
}

func Load() (*Config, error) {
	viper.SetConfigName("app")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// ✅ تعيين القيم الافتراضية
	viper.SetDefault("app.developer_mode", false)
	viper.SetDefault("app.debug", false)

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// GetConnection يعيد إعدادات اتصال قاعدة البيانات
func (d *DatabaseConfig) GetConnection(name string) (ConnectionConfig, bool) {
	conn, ok := d.Connections[name]
	return conn, ok
}

// GetDefaultConnection يعيد إعدادات الاتصال الافتراضية
func (d *DatabaseConfig) GetDefaultConnection() (ConnectionConfig, bool) {
	return d.GetConnection(d.Default)
}
