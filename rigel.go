package rigel

import "fmt"

type Rigel struct {
	ServerURL  string
	ConfigName string
	SchemaName string
}

func (r *Rigel) Check() error {
	if r.ServerURL == "" {
		return fmt.Errorf("ServerURL cannot be empty")
	}

	if r.ConfigName == "" {
		return fmt.Errorf("ConfigName cannot be empty")
	}

	if r.SchemaName == "" {
		return fmt.Errorf("SchemaName cannot be empty")
	}

	return nil
}

func (r *Rigel) LoadConfig(config any) error {
	// use rigel client to load config
	return nil
}
