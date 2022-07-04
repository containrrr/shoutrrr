package confwriter

import "github.com/containrrr/shoutrrr/pkg/format"

func (cw *ConfWriter) WriteGetURL() {
	cw.writeSection("GetURL")

	wl(`// GetURL returns a URL representation of it's current field values`)
	wl(`func (config *Config) GetURL() *url.URL {`)

	wl(`	return &url.URL{`)
	wl(getUserInfoSetter(cw.urlProps))
	if p, found := cw.urlProps[format.URLHost]; found {
		wf(`		Host: config.%v,`, p)
	}
	if p, found := cw.urlProps[format.URLPath1]; found {
		wf(`		Path: config.%v,`, p)
	} else {
		wl(`		Path: "/",`)
	}
	wf(`		RawQuery: config.QueryValues().Encode(),`)
	wf(`		Scheme: Scheme,`)
	wl(`	}`)
	wl(`}`)
	wl()
}

func (cw *ConfWriter) WriteSetURL() {

	cw.writeSection("SetURL")

	wl(`// SetURL updates a ServiceConfig from a URL representation of it's field values`)
	wl(`func (config *Config) SetURL(url *url.URL) error {`)
	wl(`	updates := make(map[string]string, propCount)`)
	if p, found := cw.urlProps[format.URLHost]; found {
		wf(`	updates[%q] = url.Host`, p)
	}
	if p, found := cw.urlProps[format.URLUser]; found {
		wf(`	updates[%q] = url.User.Username()`, p)
	}
	if p, found := cw.urlProps[format.URLPassword]; found {
		wl(`	if pwd, found := url.User.Password(); found {`)
		wf(`		updates[%q] = pwd`, p)
		wl(`	}`)
	}
	if p, found := cw.urlProps[format.URLPath1]; found {
		wf(`	updates[%q] = url.Path`, p)
	}
	wl()
	wl(`	for key, value := range url.Query() {`)
	wl(`		propName, err := propNameFromKey(key)`)
	wl(`		if err == nil {`)
	wl(`			updates[propName] = value[0]`)
	wl(`		} else if key != "title" {`)
	wf(`			return fmt.Errorf("invalid key %%q", key)`)
	wl(`		}`)
	wl(`	}`)
	wl()
	wl(`	err := config.Update(updates); if err != nil {`)
	wl(`		return err`)
	wl(`	}`)
	wl()
	for _, pn := range cw.propNames {
		prop := cw.spec.Props[pn]
		if !prop.Required {
			continue
		}
		wf(`	if %v {`, prop.Type.EmptyCall(prop, "config."+pn))
		wf(`		return fmt.Errorf("%v missing from config URL")`, pn)
		wl(`	}`)
		wl()
	}

	wl(`	return nil`)
	wl(`}`)
	wl()
}
