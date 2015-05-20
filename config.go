package main

import (
	"log"
	"reflect"

	"github.com/EverythingMe/gofigure"
	"github.com/EverythingMe/gofigure/yaml"
)

func LoadHops(configFile string) (hops, error) {
	configMap, err := loadConfigFile(configFile)
	if err != nil {
		return nil, err
	}
	hops, err := loadConfigMap(configMap)
	if err != nil {
		return nil, err
	}
	return hops, nil
}

func loadConfigMap(configMap map[string]interface{}) (hops, error) {
	hopsConfig := make(hops)
	for name, def := range configMap {
		h := make([]Hop, 0)
		switch deflT := def.(type) {
		case []interface{}:
			firstDef := deflT[0]
			switch defT := firstDef.(type) {
			case map[interface{}]interface{}:
				for hopType, params := range defT {
					switch hopType {
					case "docker":
						switch pT := params.(type) {
						case map[interface{}]interface{}:
							d := &Docker{}
							switch image := pT["image"].(type) {
							case string:
								d.Image = image
							}
							switch command := pT["command"].(type) {
							case string:
								d.Command = command
							}
							h = append(h, d)
						default:
							log.Println("Wrong Params: ", reflect.TypeOf(params))
						}
					default:
						log.Println("Wrong hop type: ", reflect.TypeOf(hopType))
					}
				}
			default:
				log.Println("Wrong hop def: ", reflect.TypeOf(firstDef))
			}
		}
		hopsConfig[name] = h
	}
	return hopsConfig, nil
}

func loadConfigFile(configFile string) (map[string]interface{}, error) {
	var configMap map[string]interface{}

	loader := gofigure.NewLoader(yaml.Decoder{}, true)

	err := loader.LoadFile(&configMap, configFile)
	if err != nil {
		return nil, err
	}

	return configMap, nil
}
