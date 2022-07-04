package confwriter

func (cw *ConfWriter) WriteProps() error {
	spec := cw.spec

	cw.writeSection("Props")

	wl(`type Config struct {`)
	for i, p := range cw.propNames {
		ps := spec.Props[p]
		spec.Props[p].Name = p
		tags := legacyTags(spec.Props[p])
		wf(`	%v %v %v`, cw.propNamesPadded[i], goType(ps), tags)
	}
	wl(`}`)

	wl()
	wl(`type configProp int`)
	wl(`const (`)
	for i, paddedName := range cw.propNamesPadded {
		wf(`	prop%v configProp = %v`, paddedName, i)
	}
	wf(`	propCount = %d`, len(cw.propNames))
	wl(`)`)

	wl(`var propNames = []string{`)
	for _, name := range cw.propNames {
		wf(`	%q,`, name)
	}
	wl(`}`)

	wl()
	wl(`// Note that propKeys may not align with propNames, as a property can have no or multiple keys`)
	wl(`var propKeys = []string{`)
	for _, key := range cw.keys {
		wf(`	%q,`, key)
	}
	wl(`}`)

	wl()

	wl(`var keyProp = []configProp{`)
	for _, key := range cw.keys {
		wf(`	prop%v,`, cw.keyProps[key])
	}
	wl(`}`)

	wl()

	wl(`var defaultValues = []string{`)
	for _, name := range cw.propNames {
		wf(`	%q,`, spec.Props[name].DefaultValue)
	}
	wl(`}`)

	wl()

	wl(`var primaryKeys = []int{`)
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
		wf(`	%v,`, keyIndex)
	}
	wl(`}`)

	wl()

	return nil
}
