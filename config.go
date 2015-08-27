// Loads Hopper configuration
package main

import (
	"fmt"

	"github.com/EverythingMe/gofigure"
	"github.com/EverythingMe/gofigure/yaml"
	v "github.com/gima/govalid/v1"
	"github.com/mitchellh/mapstructure"
)

// loadHops defined in configFile
func loadHops(configFile string) (hops, error) {
	configMap, err := loadConfigFile(configFile)
	if err != nil {
		return nil, err
	}
	hops, err := parseConfigMap(configMap)
	if err != nil {
		return nil, err
	}
	return hops, nil
}

// Hop name mapped to untyped hop list
type rawHops map[string]([]map[string]interface{})

// Parses YAML structure to typed list of Hops
func parseConfigMap(configMap map[string]interface{}) (hops, error) {
	if err := validateConfig(configMap); err != nil {
		return nil, err
	}

	hopsConfig := make(hops)
	rh := make(rawHops)
	err := mapstructure.Decode(configMap, &rh)
	if err != nil {
		panic(err)
	}

	for name, defs := range rh {
		for _, def := range defs {
			if dockMap, ok := def["docker"]; ok {
				docker := Docker{}
				err := mapstructure.Decode(dockMap, &docker)
				if err != nil {
					panic(err)
				}
				hopsConfig[name] = []Hop{&docker}
			} else {
				return nil, fmt.Errorf(
					"There's no definition of docker for %v", name)
			}
		}
	}
	return hopsConfig, nil
}

// Reads YAML config and returns its content as plain Go generic maps
func loadConfigFile(configFile string) (map[string]interface{}, error) {
	RedirectStandardLog("gofigure")
	defer ResetStandardLog()

	var configMap map[string]interface{}

	loader := gofigure.NewLoader(yaml.Decoder{}, true)

	err := loader.LoadFile(&configMap, configFile)
	if err != nil {
		return nil, err
	}

	return configMap, nil
}

// Validates config structure
func validateConfig(configMap map[string]interface{}) error {
	schema := v.Object(v.ObjKeys(v.String()),
		v.ObjValues(v.Array(v.ArrEach(v.Object(
			v.ObjKV("docker", v.Object(
				v.ObjKV("image", v.String()),
				v.ObjKV("command", v.String()),
				v.ObjKV("permissions", v.Optional(v.Object(
					v.ObjKV("cwd", v.Boolean()),
				))),
			)),
		)))),
	)
	if path, err := schema.Validate(configMap); err != nil {
		return fmt.Errorf("Invalid hop definition at %s. Error (%s)",
			path, err)
	}
	return nil
}
