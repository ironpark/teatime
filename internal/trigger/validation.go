package trigger

import (
	"fmt"
	"reflect"

	"github.com/go-viper/mapstructure/v2"
)

// ConfigValidator defines the interface for trigger configuration validation.
// All trigger configuration structs should implement this interface.
type ConfigValidator interface {
	// Validate checks the configuration and returns an error if invalid.
	Validate() error
}

// BindAndValidate binds configuration data to a struct and validates it.
// The config parameter must be a pointer to a struct that implements ConfigValidator.
func BindAndValidate[T ConfigValidator](config T, configMap map[string]any) error {
	// Validate that config is a struct pointer
	configValue := reflect.ValueOf(config)
	if configValue.Kind() != reflect.Pointer {
		return fmt.Errorf("config must be a pointer, got %T", config)
	}
	
	if configValue.IsNil() {
		return fmt.Errorf("config pointer cannot be nil")
	}
	
	elem := configValue.Elem()
	if elem.Kind() != reflect.Struct {
		return fmt.Errorf("config must be a pointer to struct, got pointer to %s", elem.Kind())
	}
	
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  config,
		TagName: "mapstructure",
	})
	if err != nil {
		return fmt.Errorf("failed to create decoder: %w", err)
	}

	if err := decoder.Decode(configMap); err != nil {
		return fmt.Errorf("failed to decode config: %w", err)
	}

	if err := config.Validate(); err != nil {
		return err
	}

	return nil
}