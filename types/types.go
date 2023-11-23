package types

import (
	"context"
)

type Schema struct {
	Name        string  // "corebanking"
	Version     int     // 2
	Fields      []Field // [{"name": "MAXDAYS", "type": "int"}, {"name": "MINBAL", "type": "float"}]
	Description string
}

type Field struct {
	Name string `json:"name"`
	Type string `json:"type"`
}
type Storage interface {
	Get(ctx context.Context, key string) (string, error)
	Put(ctx context.Context, key string, value string) error
}
