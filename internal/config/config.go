package config

import (
    "os"
    "strconv"

    "github.com/joho/godotenv"
)

type Config struct {
    Port        string
    DBUrl       string
    JWTSecret   string
    AdminToken  string
    DefaultRTP  float64
}

func Load() *Config {
    _ = godotenv.Load()

    rtp, _ := strconv.ParseFloat(getEnv("DEFAULT_RTP", "97.0"), 64)

    return &Config{
        Port:       getEnv("APP_PORT", "8080"),
        DBUrl:      os.Getenv("DB_URL"),
        JWTSecret:  os.Getenv("JWT_SECRET"),
        AdminToken: os.Getenv("ADMIN_TOKEN"),
        DefaultRTP: rtp,
    }
}

func getEnv(k, def string) string {
    if v := os.Getenv(k); v != "" {
        return v
    }
    return def
}