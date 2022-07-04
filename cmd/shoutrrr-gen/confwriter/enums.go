package confwriter

func (cw *ConfWriter) WriteEnums() {

	cw.writeSection("Enums / Options")

	wl(`func (config *Config) Enums() map[string]types.EnumFormatter {`)
	wl(`	return map[string]types.EnumFormatter{`)
	for p := range cw.enumProps {
		wf(`		%q: %sOptions.Formatter,`, p, p)
	}
	wl(`	}`)
	wl(`}`)
	wl()
	// wl(`var (`)
	for p, vals := range cw.enumProps {

		cw.writeSubSection(p + " Option")

		typeName := optionTypeName(p)
		structName := typeName + "Vals"
		wf(`type %v int`, typeName)
		wl()
		wf(`type %v struct {`, structName)
		for _, val := range vals {
			wf(`	%s	%s`, val, typeName)
		}
		wl(`	Formatter types.EnumFormatter`)
		wl(`}`)
		wl()
		wf(`var %vOptions = &%v {`, p, structName)
		for i, val := range vals {
			wf(`	%s:	%v,`, val, i)
		}
		wl(`	Formatter: format.CreateEnumFormatter([]string{`)
		for _, val := range vals {
			wf(`		%q,`, val)
		}
		wl(`	}),`)
		wl(`}`)
		wl()
		wf(`func (ov *%v) Parse(v string) (%v, error) {`, structName, typeName)
		wf(`	if val := ov.Formatter.Parse(v); val != format.EnumInvalid {`)
		wf(`		return %v(val), nil`, typeName)
		wf(`	} else {`)
		wf(`		return %v(val), fmt.Errorf("invalid option %%q for %v", v)`, typeName, p)
		wf(`	}`)
		wf(`}`)
		wl()
	}
}
