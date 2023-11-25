package rigel

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/ssd532/rigel/types"
)

const RigelPrefix = "/remiges/rigel"

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
	schemaFieldsKey := getSchemaFieldsPath(schemaName, schemaVersion)

	fieldsStr, err := r.Storage.Get(context.Background(), schemaFieldsKey)
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
func (r *Rigel) getConfigValue(schemaName string, schemaVersion int, paramName string) (string, error) {
	// Construct the key for the parameter
	key := getConfKeyPath(schemaName, schemaVersion, paramName)

	// Retrieve the parameter value from the storage
	value, err := r.Storage.Get(context.Background(), key)
	if err != nil {
		return "", err
	}

	return value, nil
}

// constructConfigMap constructs a configuration map based on the provided schema, schemaName, and schemaVersion.
func (r *Rigel) constructConfigMap(schema *types.Schema, schemaName string, schemaVersion int) (map[string]any, error) {
	// Construct the configuration map
	config := make(map[string]interface{})
	for _, field := range schema.Fields {
		// Retrieve the configuration value for the field
		valueStr, err := r.getConfigValue(schemaName, schemaVersion, field.Name)
		if err != nil {
			return nil, err
		}

		// Convert the value to the correct type based on the field type
		var value interface{}
		switch field.Type {
		case "int":
			value, err = strconv.Atoi(valueStr)
			if err != nil {
				return nil, fmt.Errorf("failed to convert value to int: %w", err)
			}
		case "bool":
			value, err = strconv.ParseBool(valueStr)
			if err != nil {
				return nil, fmt.Errorf("failed to convert value to bool: %w", err)
			}
		default:
			// Assume the value is a string if the field type is not "int" or "bool"
			value = valueStr
		}

		// Add the value to the configuration map
		config[field.Name] = value
	}
	return config, nil
}
