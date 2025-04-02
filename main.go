package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dogeorg/wow-payment/internal/config"
	"github.com/dogeorg/wow-payment/internal/handler"
)

func main() {
	cfg, err := config.LoadConfig("config.toml")
	if err != nil {
		log.Fatalf("Much Sad! Config error: %v", err)
	}

	http.HandleFunc("/register", handler.RegisterHandler(cfg))

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Wow! Starting Wow Payment Server on %s", addr)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("Such Sad! Wow Payment Server failed: %v", err)
	}
}
