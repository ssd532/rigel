package rigel

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ssd532/rigel/types"
)

// Rigel represents a client for Rigel configuration manager server.
type Rigel struct {
	Storage types.Storage
}

// New creates a new instance of Rigel with the provided Storage interface.
// The Storage interface is used by Rigel to interact with the underlying storage system.
// Currently, only etcd is supported as a storage system.
func New(storage types.Storage) *Rigel {
	return &Rigel{
		Storage: storage,
	}
}

// LoadConfig retrieves the configuration data associated with the provided configName.
// It then unmarshals this data into the provided configStruct.
// The configStruct parameter should be a pointer to a config struct used in the application.
func (r *Rigel) LoadConfig(schemaName string, schemaVersion int, configName string, configStruct any) error {
	// Retrieve the schema
	schema, err := r.getSchema(schemaName, schemaVersion)
	if err != nil {
		return err
	}

	// Construct the configuration map
	configMap, err := r.constructConfigMap(schema, schemaName, schemaVersion)
	if err != nil {
		return err
	}

	// Marshal the configuration map into a JSON string
	configJSON, err := json.Marshal(configMap)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Unmarshal the JSON string into the provided configStruct
	err = json.Unmarshal(configJSON, configStruct)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config value: %w", err)
	}

	return nil
}

// getSchema retrieves a schema from the storage based on the provided schemaName and schemaVersion.
func (r *Rigel) getSchema(schemaName string, schemaVersion int) (*types.Schema, error) {
	// Construct the base key for the schema
	baseKey := fmt.Sprintf("/remiges/rigel/schema/%s/%d", schemaName, schemaVersion)

	// Retrieve the fields from the storage
	fieldsKey := baseKey + "/fields"
	fieldsStr, err := r.Storage.Get(context.Background(), fieldsKey)
	if err != nil {
		return nil, err
	}
	var fields []types.Field
	err = json.Unmarshal([]byte(fieldsStr), &fields)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal fields: %w", err)
	}

	// Construct the schema
	schema := &types.Schema{
		Name:    schemaName,
		Version: schemaVersion,
		Fields:  fields,
	}

	return schema, nil
}

// getConfigValue retrieves a configuration value from Rigel based on the provided schemaName, schemaVersion, and paramName.
func (r *Rigel) getConfigValue(schemaName string, schemaVersion int, paramName string) (interface{}, error) {
	// Construct the key for the parameter
	key := fmt.Sprintf("/remiges/rigel/conf/%s/%d/%s", schemaName, schemaVersion, paramName)

	// Retrieve the parameter value from the storage
	paramValue, err := r.Storage.Get(context.Background(), key)
	if err != nil {
		return nil, err
	}

	// Unmarshal the parameter value
	var value interface{}
	err = json.Unmarshal([]byte(paramValue), &value)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config value: %w", err)
	}

	return value, nil
}

// constructConfigMap constructs a configuration map based on the provided schema, schemaName, and schemaVersion.
func (r *Rigel) constructConfigMap(schema *types.Schema, schemaName string, schemaVersion int) (map[string]any, error) {
	// Construct the configuration map
	config := make(map[string]interface{})
	for _, field := range schema.Fields {
		// Retrieve the configuration value for the field
		value, err := r.getConfigValue(schemaName, schemaVersion, field.Name)
		if err != nil {
			return nil, err
		}

		// Add the value to the configuration map
		config[field.Name] = value
	}

	return config, nil
}
