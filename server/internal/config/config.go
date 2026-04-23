package config

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App    AppConfig    `yaml:"app"`
	HTTP   HTTPConfig   `yaml:"http"`
	MySQL  MySQLConfig  `yaml:"mysql"`
	Redis  RedisConfig  `yaml:"redis"`
	Auth   AuthConfig   `yaml:"auth"`
	Upload UploadConfig `yaml:"upload"`
}

type AppConfig struct {
	GinMode     string `yaml:"ginMode"`
	AutoMigrate bool   `yaml:"autoMigrate"`
}

type HTTPConfig struct {
	Addr string `yaml:"addr"`
}

type MySQLConfig struct {
	Host     string            `yaml:"host"`
	Port     int               `yaml:"port"`
	User     string            `yaml:"user"`
	Password string            `yaml:"password"`
	DBName   string            `yaml:"dbName"`
	Params   map[string]string `yaml:"params"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type AuthConfig struct {
	JWTSecret       string `yaml:"jwtSecret"`
	TokenTTLSeconds int    `yaml:"tokenTTLSeconds"`
}

type UploadConfig struct {
	BaseDir      string   `yaml:"baseDir"`
	MaxFileMB    int      `yaml:"maxFileMB"`
	MaxRequestMB int      `yaml:"maxRequestMB"`
	AllowedExt   []string `yaml:"allowedExt"`
}

type fileConfig struct {
	App    fileAppConfig    `yaml:"app"`
	HTTP   fileHTTPConfig   `yaml:"http"`
	MySQL  fileMySQLConfig  `yaml:"mysql"`
	Redis  fileRedisConfig  `yaml:"redis"`
	Auth   fileAuthConfig   `yaml:"auth"`
	Upload fileUploadConfig `yaml:"upload"`
}

type fileAppConfig struct {
	GinMode     *string `yaml:"ginMode"`
	AutoMigrate *bool   `yaml:"autoMigrate"`
}

type fileHTTPConfig struct {
	Addr *string `yaml:"addr"`
}

type fileMySQLConfig struct {
	Host     *string           `yaml:"host"`
	Port     *int              `yaml:"port"`
	User     *string           `yaml:"user"`
	Password *string           `yaml:"password"`
	DBName   *string           `yaml:"dbName"`
	Params   map[string]string `yaml:"params"`
}

type fileRedisConfig struct {
	Host     *string `yaml:"host"`
	Port     *int    `yaml:"port"`
	Password *string `yaml:"password"`
	DB       *int    `yaml:"db"`
}

type fileAuthConfig struct {
	JWTSecret       *string `yaml:"jwtSecret"`
	TokenTTLSeconds *int    `yaml:"tokenTTLSeconds"`
}

type fileUploadConfig struct {
	BaseDir      *string  `yaml:"baseDir"`
	MaxFileMB    *int     `yaml:"maxFileMB"`
	MaxRequestMB *int     `yaml:"maxRequestMB"`
	AllowedExt   []string `yaml:"allowedExt"`
}

func Load() Config {
	cfg := defaultConfig()

	if fileCfg, err := loadFromFile(getEnv("CONFIG_PATH", "config.yaml")); err == nil {
		cfg = mergeConfig(cfg, fileCfg)
	}

	applyEnvOverrides(&cfg)

	return cfg
}

func (m MySQLConfig) DSN() string {
	host := m.Host
	if host == "" {
		host = "127.0.0.1"
	}
	port := m.Port
	if port == 0 {
		port = 3306
	}
	user := m.User
	if user == "" {
		user = "root"
	}
	dbName := m.DBName
	if dbName == "" {
		dbName = "app"
	}

	params := make(map[string]string, len(m.Params))
	for k, v := range m.Params {
		params[k] = v
	}
	if _, ok := params["charset"]; !ok {
		params["charset"] = "utf8mb4"
	}
	if _, ok := params["parseTime"]; !ok {
		params["parseTime"] = "True"
	}
	if _, ok := params["loc"]; !ok {
		params["loc"] = "Local"
	}

	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var b strings.Builder
	for i, k := range keys {
		if i > 0 {
			b.WriteByte('&')
		}
		b.WriteString(k)
		b.WriteByte('=')
		b.WriteString(params[k])
	}

	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", user, m.Password, host, port, dbName, b.String())
}

func (r RedisConfig) Addr() string {
	host := r.Host
	if host == "" {
		host = "127.0.0.1"
	}
	port := r.Port
	if port == 0 {
		port = 6379
	}
	return fmt.Sprintf("%s:%d", host, port)
}

func defaultConfig() Config {
	return Config{
		App: AppConfig{
			GinMode:     "release",
			AutoMigrate: true,
		},
		HTTP: HTTPConfig{
			Addr: ":8080",
		},
		MySQL: MySQLConfig{
			Host:   "127.0.0.1",
			Port:   3306,
			User:   "root",
			DBName: "app",
			Params: map[string]string{
				"charset":   "utf8mb4",
				"parseTime": "True",
				"loc":       "Local",
			},
		},
		Redis: RedisConfig{
			Host: "127.0.0.1",
			Port: 6379,
			DB:   0,
		},
		Auth: AuthConfig{
			JWTSecret:       "dev-secret-change-me",
			TokenTTLSeconds: 7200,
		},
		Upload: UploadConfig{
			BaseDir:      "data/uploads",
			MaxFileMB:    50,
			MaxRequestMB: 100,
			AllowedExt: []string{
				"jpg", "jpeg", "png", "gif", "webp",
				"mp4", "mov",
				"pdf", "zip", "txt", "md",
			},
		},
	}
}

func loadFromFile(path string) (fileConfig, error) {
	var cfg fileConfig

	b, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func mergeConfig(base Config, incoming fileConfig) Config {
	if incoming.App.GinMode != nil {
		base.App.GinMode = *incoming.App.GinMode
	}
	if incoming.App.AutoMigrate != nil {
		base.App.AutoMigrate = *incoming.App.AutoMigrate
	}

	if incoming.HTTP.Addr != nil {
		base.HTTP.Addr = *incoming.HTTP.Addr
	}

	if incoming.MySQL.Host != nil {
		base.MySQL.Host = *incoming.MySQL.Host
	}
	if incoming.MySQL.Port != nil {
		base.MySQL.Port = *incoming.MySQL.Port
	}
	if incoming.MySQL.User != nil {
		base.MySQL.User = *incoming.MySQL.User
	}
	if incoming.MySQL.Password != nil {
		base.MySQL.Password = *incoming.MySQL.Password
	}
	if incoming.MySQL.DBName != nil {
		base.MySQL.DBName = *incoming.MySQL.DBName
	}
	if incoming.MySQL.Params != nil {
		if base.MySQL.Params == nil {
			base.MySQL.Params = map[string]string{}
		}
		for k, v := range incoming.MySQL.Params {
			base.MySQL.Params[k] = v
		}
	}

	if incoming.Redis.Host != nil {
		base.Redis.Host = *incoming.Redis.Host
	}
	if incoming.Redis.Port != nil {
		base.Redis.Port = *incoming.Redis.Port
	}
	if incoming.Redis.Password != nil {
		base.Redis.Password = *incoming.Redis.Password
	}
	if incoming.Redis.DB != nil {
		base.Redis.DB = *incoming.Redis.DB
	}

	if incoming.Auth.JWTSecret != nil {
		base.Auth.JWTSecret = *incoming.Auth.JWTSecret
	}
	if incoming.Auth.TokenTTLSeconds != nil {
		base.Auth.TokenTTLSeconds = *incoming.Auth.TokenTTLSeconds
	}

	if incoming.Upload.BaseDir != nil {
		base.Upload.BaseDir = *incoming.Upload.BaseDir
	}
	if incoming.Upload.MaxFileMB != nil {
		base.Upload.MaxFileMB = *incoming.Upload.MaxFileMB
	}
	if incoming.Upload.MaxRequestMB != nil {
		base.Upload.MaxRequestMB = *incoming.Upload.MaxRequestMB
	}
	if incoming.Upload.AllowedExt != nil {
		base.Upload.AllowedExt = incoming.Upload.AllowedExt
	}

	return base
}

func applyEnvOverrides(cfg *Config) {
	if v := getEnv("HTTP_ADDR", ""); v != "" {
		cfg.HTTP.Addr = v
	}
	if v := getEnv("GIN_MODE", ""); v != "" {
		cfg.App.GinMode = v
	}
	if v := getEnv("AUTO_MIGRATE", ""); v != "" {
		cfg.App.AutoMigrate = getEnvBool("AUTO_MIGRATE", cfg.App.AutoMigrate)
	}

	if v := getEnv("MYSQL_DSN", ""); v != "" {
		user, pass, host, port, dbName, params := parseMySQLDSN(v)
		if user != "" {
			cfg.MySQL.User = user
		}
		cfg.MySQL.Password = pass
		if host != "" {
			cfg.MySQL.Host = host
		}
		if port != 0 {
			cfg.MySQL.Port = port
		}
		if dbName != "" {
			cfg.MySQL.DBName = dbName
		}
		if params != nil {
			if cfg.MySQL.Params == nil {
				cfg.MySQL.Params = map[string]string{}
			}
			for k, val := range params {
				cfg.MySQL.Params[k] = val
			}
		}
	}

	if v := getEnv("REDIS_ADDR", ""); v != "" {
		host, port := splitHostPort(v)
		if host != "" {
			cfg.Redis.Host = host
		}
		if port != 0 {
			cfg.Redis.Port = port
		}
	}
	if v := getEnv("REDIS_PASSWORD", ""); v != "" {
		cfg.Redis.Password = v
	}
	if v := getEnv("REDIS_DB", ""); v != "" {
		cfg.Redis.DB = getEnvInt("REDIS_DB", cfg.Redis.DB)
	}

	if v := getEnv("JWT_SECRET", ""); v != "" {
		cfg.Auth.JWTSecret = v
	}
	if v := getEnv("TOKEN_TTL_SECONDS", ""); v != "" {
		cfg.Auth.TokenTTLSeconds = getEnvInt("TOKEN_TTL_SECONDS", cfg.Auth.TokenTTLSeconds)
	}

	if v := getEnv("UPLOAD_BASE_DIR", ""); v != "" {
		cfg.Upload.BaseDir = v
	}
	if v := getEnv("UPLOAD_MAX_FILE_MB", ""); v != "" {
		cfg.Upload.MaxFileMB = getEnvInt("UPLOAD_MAX_FILE_MB", cfg.Upload.MaxFileMB)
	}
	if v := getEnv("UPLOAD_MAX_REQUEST_MB", ""); v != "" {
		cfg.Upload.MaxRequestMB = getEnvInt("UPLOAD_MAX_REQUEST_MB", cfg.Upload.MaxRequestMB)
	}
	if v := getEnv("UPLOAD_ALLOWED_EXT", ""); v != "" {
		cfg.Upload.AllowedExt = splitCSV(v)
	}
}

func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		out = append(out, p)
	}
	return out
}

func splitHostPort(addr string) (string, int) {
	addr = strings.TrimSpace(addr)
	if addr == "" {
		return "", 0
	}
	host, portStr, ok := strings.Cut(addr, ":")
	if !ok {
		return addr, 0
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return host, 0
	}
	return host, port
}

func parseMySQLDSN(dsn string) (user, pass, host string, port int, dbName string, params map[string]string) {
	beforeAt, afterAt, ok := strings.Cut(dsn, "@")
	if !ok {
		return "", "", "", 0, "", nil
	}
	user, pass, _ = strings.Cut(beforeAt, ":")

	afterProto := afterAt
	if strings.HasPrefix(afterProto, "tcp(") {
		afterProto = strings.TrimPrefix(afterProto, "tcp(")
		afterProto = strings.TrimSuffix(afterProto, ")")
	}

	hostPortPart, dbPart, ok := strings.Cut(afterProto, "/")
	if ok {
		host, port = splitHostPort(hostPortPart)
	} else {
		host, port = splitHostPort(afterProto)
	}

	if dbPart == "" {
		return user, pass, host, port, "", nil
	}

	dbName, query, _ := strings.Cut(dbPart, "?")
	if query == "" {
		return user, pass, host, port, dbName, nil
	}

	params = map[string]string{}
	for _, pair := range strings.Split(query, "&") {
		if pair == "" {
			continue
		}
		k, v, ok := strings.Cut(pair, "=")
		if !ok {
			continue
		}
		params[k] = v
	}
	return user, pass, host, port, dbName, params
}

func getEnv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	v := getEnv(key, "")
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}

func getEnvBool(key string, def bool) bool {
	v := getEnv(key, "")
	if v == "" {
		return def
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return def
	}
	return b
}
