package main

import (
	"io"
	"os"

	"github.com/containrrr/shoutrrr/pkg/conf"
	"github.com/containrrr/shoutrrr/pkg/ref"
	"github.com/containrrr/shoutrrr/pkg/types"
	"gopkg.in/yaml.v2"
)

type MigratedService interface {
	GetLegacyConfig() types.ServiceConfig
}

func Export(service types.Service, scheme string, w io.Writer) error {
	var configNode *ref.ContainerNode
	if migratedService, ok := service.(MigratedService); ok {
		println("Service config is migrated, using legacy config")
		configNode = ref.GetConfigFormat(migratedService.GetLegacyConfig())
	} else {
		configNode = ref.GetServiceConfigFormat(service)
	}

	configDef := conf.Spec{
		Version: 1,
		Scheme:  scheme,
		Props:   map[string]*conf.SpecProp{},
	}

	for _, item := range configNode.Items {
		field := item.Field()
		values := []string{}
		if ef := field.EnumFormatter; ef != nil {
			values = ef.Names()
		}
		configDef.Props[field.Name] = &conf.SpecProp{
			Type:         conf.ConfigPropTypeFromType(field.Type, item.TokenType()),
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
