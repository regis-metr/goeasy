package obj

import (
	"fmt"
	"reflect"
)

// FieldMapConfig contains configuration on how to transform a given field of
// a struct.
type FieldMapConfig struct {
	Source              string
	Destination         string
	GetDestinationValue func(source any) (any, error)
}

type structMapKey struct {
	source      reflect.Type
	destination reflect.Type
}

type MapperConfig struct {
	fieldMaps map[structMapKey]map[string]*FieldMapConfig
}

// ConfigureFieldMaps allows overriding of how fields are mapped for sourceT and destinationT
func ConfigureFieldMaps[sourceT any, destinationT any](mapperConfig *MapperConfig,
	fieldMapConfigs ...FieldMapConfig) error {
	var zeroSource sourceT
	var zeroDestination destinationT
	sourceType := reflect.TypeOf(zeroSource)
	destinationType := reflect.TypeOf(zeroDestination)
	if sourceType.Kind() != reflect.Struct || destinationType.Kind() != reflect.Struct {
		return fmt.Errorf("sourceT and destinationT must be structs")
	}

	structKey := structMapKey{
		source:      sourceType,
		destination: destinationType,
	}

	fieldMap := mapperConfig.fieldMaps[structKey]
	if fieldMap == nil {
		fieldMap = make(map[string]*FieldMapConfig)
	}

	for _, cfg := range fieldMapConfigs {
		if cfg.Destination == "" {
			return fmt.Errorf("destination field names must be provided")
		}

		fieldMap[cfg.Destination] = &cfg
	}
	mapperConfig.fieldMaps[structKey] = fieldMap
	return nil
}
