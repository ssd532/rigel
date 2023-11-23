package main

import (
	"fmt"
	"github.com/ssd532/rigel"
	"github.com/ssd532/rigel/etcd"
	"log"
)

type Config struct {
	DatabaseURL string `json:"database_url"`
	MaxRetries  int    `json:"max_retries"`
	EnableSSL   bool   `json:"enable_ssl"`
}

func main() {
	// Create a new EtcdStorage instance
	etcdStorage, err := etcd.NewEtcdStorage([]string{"localhost:2379"})
	if err != nil {
		log.Fatalf("Failed to create EtcdStorage: %v", err)
	}

	// Create a new Rigel instance
	rigelClient := rigel.New(etcdStorage)

	// Define a config struct
	var config Config

	// Load the config
	err = rigelClient.LoadConfig("appConfig", 1, "appConfig", &config)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Print the loaded config
	fmt.Printf("DatabaseURL: %s\n", config.DatabaseURL)
	fmt.Printf("MaxRetries: %d\n", config.MaxRetries)
	fmt.Printf("EnableSSL: %t\n", config.EnableSSL)
}
