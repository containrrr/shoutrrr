package main

import (
	"io"
	"os"

	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/types"
	"gopkg.in/yaml.v2"
)

type MigratedService interface {
	GetLegacyConfig() types.ServiceConfig
}

func Export(service types.Service, scheme string, w io.Writer) error {
	var configNode *format.ContainerNode
	if migratedService, ok := service.(MigratedService); ok {
		println("Service config is migrated, using legacy config")
		configNode = format.GetConfigFormat(migratedService.GetLegacyConfig())
	} else {
		configNode = format.GetServiceConfigFormat(service)
	}

	configDef := format.ConfigSpec{
		Version: 1,
		Scheme:  scheme,
		Props:   map[string]*format.ConfigSpecProp{},
	}

	for _, item := range configNode.Items {
		field := item.Field()
		values := []string{}
		if ef := field.EnumFormatter; ef != nil {
			values = ef.Names()
		}
		configDef.Props[field.Name] = &format.ConfigSpecProp{
			Type:         format.ConfigPropTypeFromType(field.Type, item.TokenType()),
			Description:  field.Description,
			DefaultValue: field.DefaultValue,
			Template:     field.Template,
			Required:     field.Required,
			URLParts:     field.URLParts,
			Title:        field.Title,
			Base:         field.Base,
			Keys:         field.Keys,
			Values:       values,
		}
	}

	enc := yaml.NewEncoder(w)
	defer enc.Close()
	return enc.Encode(configDef)

	// bytes, err := yaml.Marshal(configDef)
	// if err != nil {
	// 	return err
	// }

	// if

	// os.Stdout.Write(bytes)

	// return nil
}

func writeConfigDef(fileName string, v interface{}) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	enc := yaml.NewEncoder(file)
	defer enc.Close()
	return enc.Encode(v)
}
