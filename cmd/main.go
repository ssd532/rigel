package main

import (
	"fmt"
	"github.com/ssd532/rigel"
	"log"

	"github.com/ssd532/rigel/etcd"
)

type Config struct {
	Key1 string `json:"key1"`
	Key2 int    `json:"key2"`
	Key3 bool   `json:"key3"`
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
	err = rigelClient.LoadConfig("schemaName", 1, "configName", &config)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Print the loaded config
	fmt.Printf("Key1: %s\n", config.Key1)
	fmt.Printf("Key2: %d\n", config.Key2)
	fmt.Printf("Key3: %t\n", config.Key3)
}
