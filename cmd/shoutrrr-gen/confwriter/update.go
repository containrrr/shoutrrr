package confwriter

func (cw *ConfWriter) WriteUpdate() {
	wl(`// Update updates the Config from a map of it's properties`)
	wl(`func (config *Config) Update(updates map[string]string) error {`)
	wl(`	var last_err error`)
	wl(`	for key, value := range updates {`)
	wl(`		switch key {`)
	for _, name := range cw.propNames {
		p := cw.spec.Props[name]
		wf(`		case %q:`, name)
		wf(`			if val, err := %v; err != nil {`, p.Type.ParserCall(p, "value"))
		wf(`				last_err = err`)
		wf(`			} else {`)
		wf(`				config.%v = val`, name)
		wf(`			}`)
	}
	wl(`		default:`)
	wl(`			last_err = fmt.Errorf("invalid key")`)
	wl(`		}`)
	wl(`		if last_err != nil {`)
	wf(`			return fmt.Errorf("failed to set value for %%q: %%v", key, last_err)`)
	wl(`		}`)
	wl(`	}`)
	wl(`	return nil`)
	wl(`}`)
	wl()

	wl(`// Update updates the Config from a map of it's properties`)
	wl(`func (config *Config) propValue(prop configProp) string {`)
	wl(`	switch prop {`)
	for _, name := range cw.propNames {
		p := cw.spec.Props[name]
		wf(`	case prop%v:`, name)
		wf(`		return %v`, p.Type.FormatCall(p, "config."+name))
	}
	wl(`	default:`)
	wl(`		return ""`)
	wl(`	}`)
	wl(`}`)
	wl()
}
