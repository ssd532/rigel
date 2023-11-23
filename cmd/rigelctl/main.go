package main

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/ssd532/rigel"
	"github.com/ssd532/rigel/etcd"
)

func main() {
	// Create a new EtcdStorage instance
	etcdStorage, err := etcd.NewEtcdStorage([]string{"localhost:2379"})
	if err != nil {
		log.Fatalf("Failed to create EtcdStorage: %v", err)
	}
	// Create a new Rigel instance
	rigelClient := rigel.New(etcdStorage)

	// Create the root command
	rootCmd := &cobra.Command{
		Use:   "rigelctl",
		Short: "CLI for managing Rigel schemas and configs",
	}

	// Create the addSchema command
	addSchemaCmd := &cobra.Command{
		Use:   "addSchema [name] [version] [fields]",
		Short: "Add a new schema",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return addSchemaCommand(rigelClient, args)
		},
	}

	// Add the addSchema command to the root command
	rootCmd.AddCommand(addSchemaCmd)

	// Execute the root command
	rootCmd.Execute()
}
