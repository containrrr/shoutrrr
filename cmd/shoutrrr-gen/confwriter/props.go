package confwriter

func (cw *ConfWriter) writeProps() error {
	spec := cw.spec

	cw.writeSection("Props")

	if len(cw.propNames) < 1 {
		wl(`type Config struct {}`)
	} else {
		wl(`type Config struct {`)
	}

	for i, p := range cw.propNames {
		ps := spec.Props[p]
		spec.Props[p].Name = p
		tags := legacyTags(spec.Props[p])
		wf(`	%v %v %v`, cw.propNamesPadded[i], ps.Type.GoType(ps), tags)
	}

	if len(cw.propNames) > 0 {
		wl(`}`)
	}

	wl()
	wl(`type configProp int`)
	wl(`const (`)
	for i, paddedName := range cw.propNamesPadded {
		wf(`	prop%v configProp = %v`, paddedName, i)
	}
	wf(`	propCount = %d`, len(cw.propNames))
	wl(`)`)

	wl(`var propInfo = types.ConfigPropInfo{`)
	wl(`	PropNames: []string{`)
	for _, name := range cw.propNames {
		wf(`		%q,`, name)
	}
	wl(`	},`)

	wl()
	wl(`	// Note that propKeys may not align with propNames, as a property can have no or multiple keys`)
	wl(`	Keys: []string{`)
	for _, key := range cw.keys {
		wf(`		%q,`, key)
	}
	wl(`	},`)

	wl()

	wl(`	DefaultValues: []string{`)
	for _, name := range cw.propNames {
		wf(`		%q,`, spec.Props[name].DefaultValue)
	}
	wl(`	},`)

	wl()

	wl(`	PrimaryKeys: []int{`)
	for _, name := range cw.propNames {
		propKeys := spec.Props[name].Keys
		keyIndex := -1
		if len(propKeys) > 0 {
			for i, key := range cw.keys {
				if key == propKeys[0] {
					keyIndex = i
					break
				}
			}
		}
		wf(`		%v,`, keyIndex)
	}
	wl(`	},`)
	wl()
	wl(`	KeyPropIndexes: map[string]int{`)
	for _, key := range cw.keys {
		wf(`		%q: %v,`, key, cw.keyProps[key])
	}
	wl(`	},`)

	wl()
	wl(`}`)

	wl()
	wl(`func (_ *Config) PropInfo() *types.ConfigPropInfo {`)
	wl(`	return &propInfo`)
	wl(`}`)
	wl()

	return nil
}
