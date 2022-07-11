package confwriter

func (cw *ConfWriter) writeUpdate() {
	wl(`// Update updates the Config from a map of it's properties`)
	wl(`func (config *Config) Update(updates map[int]string) error {`)
	if len(cw.propNames) > 0 {
		wl(`	var last_err error`)
		wl(`	for index, value := range updates {`)
		wl(`		switch configProp(index) {`)
		for _, name := range cw.propNames {
			p := cw.spec.Props[name]
			wf(`		case prop%v:`, name)
			wf(`			if val, err := %v; err != nil {`, p.Type.ParserCall(p, "value"))
			wf(`				last_err = err`)
			wf(`			} else {`)
			wf(`				config.%v = val`, name)
			wf(`			}`)
		}
		wl(`		default:`)
		wl(`			return fmt.Errorf("invalid key")`)
		wl(`		}`)
		wl(`		if last_err != nil {`)
		wf(`			return fmt.Errorf("failed to set value for %%v: %%v", propInfo.PropNames[index], last_err)`)
		wl(`		}`)
		wl(`	}`)
	}
	wl(`	return nil`)
	wl(`}`)
	wl()

	wl(`// Update updates the Config from a map of it's properties`)
	wl(`func (config *Config) PropValue(prop int) string {`)
	wl(`	switch configProp(prop) {`)
	for _, name := range cw.propNames {
		p := cw.spec.Props[name]
		value := "config." + name
		wf(`	case prop%v:`, name)
		wf(`		return %v`, p.Type.FormatCall(p, value))
	}
	wl(`	default:`)
	wl(`		return ""`)
	wl(`	}`)
	wl(`}`)
	wl()
}
