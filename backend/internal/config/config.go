package config

import "os"

type Config struct {
	Port              string
	SupabaseURL       string
	SupabaseJWTSecret string
	DatabaseURL       string
}

func Load() *Config {
	return &Config{
		Port:              getEnv("PORT", "8080"),
		SupabaseURL:       getEnv("SUPABASE_URL", ""),
		SupabaseJWTSecret: getEnv("SUPABASE_JWT_SECRET", ""),
		DatabaseURL:       getEnv("DATABASE_URL", ""),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
