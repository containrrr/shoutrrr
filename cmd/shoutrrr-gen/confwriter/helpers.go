package confwriter

func (cw *ConfWriter) WriteHelpers() {

	cw.writeSection("Helpers")

	wf(`
func propNameFromKey(key string) (string, error) {
	key = strings.ToLower(key)
	for i, pk := range propKeys {
		if key == pk {
			return propNames[keyProp[i]], nil
		}
	}
	return "", fmt.Errorf("invalid key %%q", key)
}

// UpdateFromParams updates the configuration from the supplied params
func (config *Config) UpdateFromParams(params *types.Params) error {
	if params == nil {
		return nil
	}
	updates := make(map[string]string, len(*params))
	for key, value := range *params {
		propName, err := propNameFromKey(key)
		if err == nil {
			updates[propName] = value
		} else if key != "title" {
			return fmt.Errorf("invalid key %%q", key)
		}
	}
	return config.Update(updates)
}


// UpdateFromQuery updates the configuration from the supplied query values
func (config *Config) UpdateFromQuery(values url.Values) error {
	updates := make(map[string]string, len(values))
	for key, value := range values {
		propName, err := propNameFromKey(key)
		if err == nil {
			updates[propName] = value[0]
		} else if key != "title" {
			return fmt.Errorf("invalid key %%q", key)
		}
	}
	return config.Update(updates)
}

// Init sets all the Config properties to their default values
func (config *Config) Init() error {
	updates := make(map[string]string, propCount)
	for i, name := range propNames {
		updates[name] = defaultValues[i]
	}
	return config.Update(updates)
}

// QueryValues returns a url.Values populated from the configuration
func (config *Config) QueryValues() url.Values {
	values := make(url.Values, propCount)
	for i := range propNames {
		if primaryKeys[i] < 0 {
			continue
		}
		value := config.propValue(configProp(i))
		if value == defaultValues[i] {
			continue
		}
		values.Set(propKeys[primaryKeys[i]], config.propValue(configProp(i)))
	}
	return values
}
	`)
}
