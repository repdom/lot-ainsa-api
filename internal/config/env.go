package config

import "os"

type Env struct {
	Port string
}

func LoadEnv() *Env {
	return &Env{
		Port: getEnv("PORT", "8081"),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

func (env Env) GetEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
