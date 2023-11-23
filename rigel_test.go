package rigel

import (
	"context"
	"fmt"
	"testing"

	"github.com/ssd532/rigel/etcd"
	"github.com/ssd532/rigel/types"
)

func TestNewRigelClient(t *testing.T) {
	etcdStorage := &etcd.EtcdStorage{} // Mocked EtcdStorage
	rigelClient := New(etcdStorage)

	if rigelClient == nil {
		t.Fatalf("Expected rigelClient to be not nil")
	}

	if rigelClient.Storage != etcdStorage {
		t.Errorf("Expected rigelClient.Storage to be equal to provided etcdStorage")
	}
}

type mockStorage struct {
	getFunc func(ctx context.Context, key string) (string, error)
}

func (m *mockStorage) Get(ctx context.Context, key string) (string, error) {
	return m.getFunc(ctx, key)
}

func TestGetSchema(t *testing.T) {
	// Mocked Storage
	mockStorage := &mockStorage{
		getFunc: func(ctx context.Context, key string) (string, error) {
			// Return a predefined schema JSON string
			if key == "/remiges/rigel/schema/schemaName/1/fields" {
				return `[{"name": "key1", "type": "string"}, {"name": "key2", "type": "int"}, {"name": "key3", "type": "bool"}]`, nil
			}
			return "", fmt.Errorf("unexpected key: %s", key)
		},
	}

	rigelClient := New(mockStorage)

	// Call getSchema
	schema, err := rigelClient.getSchema("schemaName", 1)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check if the returned schema is correct
	if schema.Name != "schemaName" || schema.Version != 1 || len(schema.Fields) != 3 {
		t.Errorf("Returned schema is incorrect")
	}
	if schema.Fields[0].Name != "key1" || schema.Fields[0].Type != "string" {
		t.Errorf("Field 1 is incorrect")
	}
	if schema.Fields[1].Name != "key2" || schema.Fields[1].Type != "int" {
		t.Errorf("Field 2 is incorrect")
	}
	if schema.Fields[2].Name != "key3" || schema.Fields[2].Type != "bool" {
		t.Errorf("Field 3 is incorrect")
	}
}

func TestGetConfigValue(t *testing.T) {
	// Mocked Storage
	mockStorage := &mockStorage{
		getFunc: func(ctx context.Context, key string) (string, error) {
			// Return a predefined config value JSON string
			return "value", nil
		},
	}

	rigelClient := New(mockStorage)

	// Call getConfigValue
	value, err := rigelClient.getConfigValue("schemaName", 1, "key")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check if the returned value is correct
	if value != "value" {
		t.Errorf("Expected value to be 'value', got '%v'", value)
	}
}

func TestConstructConfigMap(t *testing.T) {
	// Mocked Storage
	mockStorage := &mockStorage{
		getFunc: func(ctx context.Context, key string) (string, error) {
			// Return a predefined config value JSON string
			return "value", nil
		},
	}

	rigelClient := New(mockStorage)

	// Define a schema
	schema := &types.Schema{
		Name:    "schemaName",
		Version: 1,
		Fields: []types.Field{
			{Name: "key", Type: "string"},
		},
	}

	// Call constructConfigMap
	configMap, err := rigelClient.constructConfigMap(schema, "schemaName", 1)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check if the returned config map is correct
	if value, ok := configMap["key"]; !ok || value != "value" {
		t.Errorf("Expected configMap['key'] to be 'value', got '%v'", value)
	}
}

func TestLoadConfig(t *testing.T) {
	// Mocked Storage
	mockStorage := &mockStorage{
		getFunc: func(ctx context.Context, key string) (string, error) {
			// Return a predefined schema JSON string for getSchema
			if key == "/remiges/rigel/schema/schemaName/1/fields" {
				return `[{"name": "key1", "type": "string"}, {"name": "key2", "type": "int"}, {"name": "key3", "type": "bool"}]`, nil
			}
			// Return a predefined config value JSON string for getConfigValue
			switch key {
			case "/remiges/rigel/conf/schemaName/1/key1":
				return "value1", nil
			case "/remiges/rigel/conf/schemaName/1/key2":
				return `2`, nil
			case "/remiges/rigel/conf/schemaName/1/key3":
				return `true`, nil
			default:
				return "", fmt.Errorf("unexpected key: %s", key)
			}
		},
	}

	rigelClient := New(mockStorage)

	var config struct {
		Key1 string `json:"key1"`
		Key2 int    `json:"key2"`
		Key3 bool   `json:"key3"`
	}
	err := rigelClient.LoadConfig("schemaName", 1, "configName", &config)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if config.Key1 != "value1" {
		t.Errorf("Expected config.Key1 to be 'value1', got '%s'", config.Key1)
	}
	if config.Key2 != 2 {
		t.Errorf("Expected config.Key2 to be 2, got '%d'", config.Key2)
	}
	if config.Key3 != true {
		t.Errorf("Expected config.Key3 to be true, got '%t'", config.Key3)
	}

}
func ExampleRigel_LoadConfig() {
	//// Create a new EtcdStorage instance
	//etcdStorage, err := etcd.NewEtcdStorage([]string{"localhost:2379"})
	//if err != nil {
	//	log.Fatalf("Failed to create EtcdStorage: %v", err)
	//}
	//
	//// Create a new Rigel instance
	//rigelClient := New(etcdStorage)
	//
	//// Define a config struct
	//var config struct {
	//	DatabaseURL string `json:"database_url"`
	//	APIKey      string `json:"api_key"`
	//	IsDebug     bool   `json:"is_debug"`
	//}
	//
	//// Load the config
	//err = rigelClient.LoadConfig("AppConfig", 1, "Production", &config)
	//if err != nil {
	//	log.Fatalf("Failed to load config: %v", err)
	//}
	//
	//// Print the loaded config
	//fmt.Printf("DatabaseURL: %s\n", config.DatabaseURL)
	//fmt.Printf("APIKey: %s\n", config.APIKey)
	//fmt.Printf("IsDebug: %t\n", config.IsDebug)
	//
	//// Output:
	//// DatabaseURL: postgres://user:pass@localhost:5432/dbname
	//// APIKey: abc123
	//// IsDebug: false
}
