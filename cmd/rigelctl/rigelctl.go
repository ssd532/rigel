package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/ssd532/rigel"
	"github.com/ssd532/rigel/types"
)

type RigelClient interface {
	AddSchema(ctx context.Context, name string, version int, fields []types.Field) error
}

func addSchemaCommand(client *rigel.Rigel, args []string) error {

	// Parse arguments
	if len(args) != 3 {
		return fmt.Errorf("expected 3 arguments, got %d", len(args))
	}
	name := args[0]
	version, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("invalid version: %v", err)
	}

	// Parse fields
	var fields []types.Field
	err = json.Unmarshal([]byte(args[2]), &fields)
	if err != nil {
		return fmt.Errorf("invalid fields: %v", err)
	}

	// Call AddSchema
	err = client.AddSchema(context.Background(), name, version, fields)
	if err != nil {
		return fmt.Errorf("failed to add schema: %v", err)
	}

	return nil
}
