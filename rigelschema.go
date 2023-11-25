package rigel

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ssd532/rigel/types"
)

func (r *Rigel) AddSchema(ctx context.Context, name string, version int, fields []types.Field) error {
	// Convert fields to JSON
	fieldsJson, err := json.Marshal(fields)
	if err != nil {
		return fmt.Errorf("failed to marshal fields: %v", err)
	}

	// Store fields
	fieldsKey := getSchemaFieldsPath(name, version)
	fmt.Printf("Storing fields at %s\n", fieldsKey)
	fmt.Printf("Fields: %s\n", string(fieldsJson))
	err = r.Storage.Put(ctx, fieldsKey, string(fieldsJson))
	if err != nil {
		return fmt.Errorf("failed to store fields: %v", err)
	}

	return nil
}
