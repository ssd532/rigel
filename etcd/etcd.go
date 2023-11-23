package etcd

import (
	"context"
	"fmt"
	"time"

	"github.com/ssd532/rigel/types"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const DialTimeout = 5 * time.Second

// EtcdStorage is implements the Storage interface using etcd
type EtcdStorage struct {
	Client *clientv3.Client
}

var _ types.Storage = &EtcdStorage{}

// NewEtcdStorage creates a new instance of EtcdStorage using the provided endpoints
// with default settings from the package. If an optional clientv3.Config is supplied,
// it is used to configure the etcd client, overriding the default settings.
func NewEtcdStorage(endpoints []string, config ...clientv3.Config) (*EtcdStorage, error) {
	var cfg clientv3.Config
	if len(config) > 0 {
		cfg = config[0]
	} else {
		cfg = clientv3.Config{
			Endpoints:   endpoints,
			DialTimeout: DialTimeout,
		}
	}

	cli, err := clientv3.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create etcd client: %w", err)
	}

	return &EtcdStorage{Client: cli}, nil
}

// Get retrieves a value from etcd based on the provided key.
func (e *EtcdStorage) Get(ctx context.Context, key string) (string, error) {
	resp, err := e.Client.Get(ctx, key)
	if err != nil {
		return "", fmt.Errorf("failed to get key from etcd: %w", err)
	}

	// Assuming the value is a string
	var value string
	for _, ev := range resp.Kvs {
		value = string(ev.Value)
	}

	return value, nil
}

func (e *EtcdStorage) Put(ctx context.Context, key string, value string) error {
	_, err := e.Client.Put(ctx, key, value)
	if err != nil {
		return err
	}
	return nil
}
