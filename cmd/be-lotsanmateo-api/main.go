package main

import (
	"be-lotsanmateo-api/internal/adapter/http"
	"be-lotsanmateo-api/internal/config"
	"log"
	"strings"
)

func main() {
	config.SetupLogger()
	cfg := config.LoadEnv()

	router := http.NewRouter(cfg)
	log.Println("Server started on port", cfg.Port)
	var builder strings.Builder
	builder.WriteString("0.0.0.0:")
	builder.WriteString(cfg.Port)
	err := router.Run(builder.String())
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
		return
	}
}
