package confwriter

import (
	"fmt"
	"strings"

	"github.com/containrrr/shoutrrr/pkg/conf"
	"github.com/containrrr/shoutrrr/pkg/urlpart"
)

func (cw *ConfWriter) writeGetURL() {
	cw.writeSection("GetURL")

	wl(`// GetURL returns a URL representation of it's current field values`)
	wl(`func (config *Config) GetURL() *url.URL {`)

	wl(`	return &url.URL{`)
	writeUserInfoSetter(cw.urlProps)
	writeHostSetter(cw)
	cw.writePathSetter()
	wf(`		RawQuery: conf.QueryValues(config).Encode(),`)
	wf(`		Scheme: Scheme,`)
	wl(`	}`)
	wl(`}`)
	wl()

	if cw.spec.Options.CustomQueryVars {
		if p, found := cw.urlProps[urlpart.Query]; found {
			wl(`func (config *Config) CustomQueryVars() url.Values {`)
			wf(`	return config.%v`, p)
			wl(`}`)
		}
		//
	}
}

func (cw *ConfWriter) writeSetURL() error {

	cw.writeSection("SetURL")

	wl(`// SetURL updates a ServiceConfig from a URL representation of it's field values`)
	wl(`func (config *Config) SetURL(configURL *url.URL) error {`)
	wl(`	if lc, ok := (interface{})(config).(types.ConfigWithLegacyURLSupport); ok {`)
	wl(`		configURL = lc.UpdateLegacyURL(configURL)`)
	wl(`	}`)
	wl(`	updates := make(map[int]string, propCount)`)
	writeHostGetters(cw)

	writeUserInfoGetters(cw.urlProps)

	writePathGetters(cw)
	if cw.spec.Options.CustomQueryVars {
		// wl(`	customQuery := map[string]string{}`)
		wl(`	customQuery := url.Values{}`)
	}
	wl()
	wl(`	for key, value := range configURL.Query() {`)
	wl(`		`)
	wl(`		if propIndex, found := propInfo.PropIndexFor(key); found {`)
	wl(`			updates[propIndex] = value[0]`)
	if cw.spec.Options.CustomQueryVars {
		wl(`		} else {`)
		wf(`			customQuery.Set(conf.UnescapeCustomQueryKey(key), value[0])`)
		wl(`		}`)
	} else {
		wl(`		} else if key != "title" {`)
		wf(`			return fmt.Errorf("invalid key %%q", key)`)
		wl(`		}`)
	}
	wl(`	}`)
	if cw.spec.Options.CustomQueryVars {
		if p, found := cw.urlProps[urlpart.Query]; found {
			wf(`	updates[int(prop%v)] = customQuery.Encode()`, p)
		}
	}
	wl()
	wl(`	err := config.Update(updates); if err != nil {`)
	wl(`		return err`)
	wl(`	}`)
	wl()
	for _, pn := range cw.propNames {
		prop := cw.spec.Props[pn]

		propErrorName := propErrorName(pn)

		for _, validator := range prop.Validators() {
			if err := validator.Verify(); err != nil {
				return err
			}
			variable := "config." + pn
			testCall := validator.TestCall(prop, variable)

			wf(`	if %v {`, testCall)
			wf(`		return fmt.Errorf(%v)`, validator.FailMessage(propErrorName, variable))
			wl(`	}`)
			wl()
		}
	}

	wl(`	return nil`)
	wl(`}`)
	wl()

	return nil
}

func (cw *ConfWriter) writePathSetter() {
	pathVal := `"",`
	pathParts := []string{}

	restBefore := cw.spec.Options.ReversePathPrio
	restPart := ""

	if prop, found := cw.urlProps[urlpart.Path]; found {
		if cw.spec.Props[prop].Type == conf.ListPropType {
			restPart = fmt.Sprintf(`config.%v...`, prop)
		} else {
			restPart = fmt.Sprintf(`string(config.%v)`, prop)
		}
	}

	if restBefore && restPart != "" {
		pathParts = append(pathParts, restPart)
	}

	for _, pp := range urlpart.PathParts {
		if prop, found := cw.urlProps[pp]; found {
			pathParts = append(pathParts, fmt.Sprintf(`string(config.%v)`, prop))
		} else {
			break
		}
	}

	if !restBefore && restPart != "" {
		pathParts = append(pathParts, restPart)
	}

	if len(pathParts) > 0 {
		pathVal = fmt.Sprintf(`conf.JoinPath(%v),`, strings.Join(pathParts, ", "))
	}

	wl(`		Path:`, pathVal)
}

func writeHostGetters(cw *ConfWriter) {

	var prop string
	var host_found bool
	var port_found bool
	host := `""`
	if prop, host_found = cw.urlProps[urlpart.Host]; host_found {
		host = prop
	}
	port := `""`
	if prop, port_found = cw.urlProps[urlpart.Port]; port_found {
		port = prop
	}

	if port_found && host == port {
		wf(`	updates[int(prop%v)] = configURL.Host`, host)
		return
	}

	if port_found {
		wl(`	if port := configURL.Port(); port != "" {`)
		wf(`		updates[int(prop%v)] = port`, port)
		wl(`	}`)
	}
	if host_found {
		wf(`	updates[int(prop%v)] = configURL.Hostname()`, host)
	}
}

func writePathGetters(cw *ConfWriter) {

	restProp, restFound := cw.urlProps[urlpart.Path]
	reversePrio := cw.spec.Options.ReversePathPrio
	pathProps := make([]string, 0, len(urlpart.PathParts))

	for ix := range urlpart.PathParts {
		if reversePrio {
			ix = (len(urlpart.PathParts) - ix) - 1
		}

		if prop, found := cw.urlProps[urlpart.PathParts[ix]]; found {
			pathProps = append(pathProps, prop)
		}
	}

	if len(pathProps) > 0 {
		wl()
		wl(`	pathParts := conf.SplitPath(configURL.Path)`)
	}

	for ix, prop := range pathProps {
		pos := fmt.Sprintf(`%v`, ix)
		if reversePrio {
			pos = fmt.Sprintf(`len(pathParts)-%v`, ix+1)
		}
		wf(`	if len(pathParts) > %v {`, ix)
		wf(`		updates[int(prop%v)] = pathParts[%v]`, prop, pos)
		wf(`	}`)
	}

	if restFound {
		if len(pathProps) > 0 {
			var rhs string
			if reversePrio {
				rhs = fmt.Sprintf("pathParts[:len(pathParts)-%v]", len(pathProps))
			} else {
				rhs = fmt.Sprintf("pathParts[%v:]", len(pathProps))
			}
			if cw.spec.Props[restProp].Type != conf.ListPropType {
				rhs = fmt.Sprintf(`conf.JoinPath(%v...)`, rhs)
			}

			wf(`	if len(pathParts) > %v {`, len(pathProps))
			wf(`		updates[int(prop%v)] = %v`, restProp, rhs)
			wf(`	}`)
		} else {
			wf(`	updates[int(prop%v)] = configURL.Path`, restProp)
		}
	} else {
		if len(pathProps) > 0 {
			wf(`	if len(pathParts) > %v {`, len(pathProps))
			wf(`		return fmt.Errorf("too many path items: %%v, expected %v", len(pathParts))`, len(pathProps))
			wl(`	}`)
		} else {
			wl(`	if configURL.Path != "" && configURL.Path != "/" {`)
			wf(`		return fmt.Errorf("unexpected path in config URL: %%v", configURL.Path)`)
			wl(`	}`)
		}
	}

}

func writeUserInfoGetters(urlProps map[urlpart.URLPart]string) {
	var prop string
	var user_found bool
	var pass_found bool
	user := `""`
	if prop, user_found = urlProps[urlpart.User]; user_found {
		user = prop
	}
	pass := `""`
	if prop, pass_found = urlProps[urlpart.Password]; pass_found {
		pass = prop
	}

	if pass_found && pass == user {
		wf(`	updates[int(prop%v)] = configURL.User.String()`, user)
	} else {
		if pass_found {
			wl(`	if pwd, found := configURL.User.Password(); found {`)
			wf(`		updates[int(prop%v)] = pwd`, pass)
			wl(`	}`)
		}
		if user_found {
			wf(`	updates[int(prop%v)] = configURL.User.Username()`, user)
		}
	}
}

func writeUserInfoSetter(urlProps map[urlpart.URLPart]string) {
	var prop string
	var user_found bool
	var pass_found bool
	user := `""`
	if prop, user_found = urlProps[urlpart.User]; user_found {
		user = "config." + prop
	}
	pass := `""`
	if prop, pass_found = urlProps[urlpart.Password]; pass_found {
		pass = "config." + prop
	}

	if pass_found && pass == user {
		wf(`		User: conf.UserInfoOrNil(%v.UserInfo()),`, user)
	} else if pass_found {
		wf(`		User: conf.UserInfoOrNil(url.UserPassword(%v, %v)),`, user, pass)
	} else if user_found {
		wf(`		User: conf.UserInfoOrNil(url.User(%v)),`, user)
	} else {
		wl("		// Userinfo fields are not used for configuration")
	}
}

func writeHostSetter(cw *ConfWriter) {
	var prop string
	var host_found bool
	var port_found bool
	host := `""`
	if prop, host_found = cw.urlProps[urlpart.Host]; host_found {
		host = "config." + prop
	}
	port := `""`
	if prop, port_found = cw.urlProps[urlpart.Port]; port_found {
		port = "config." + prop
	}

	if !host_found {
		wl("		// Host fields are not used for configuration")
		return
	}

	if host == port || !port_found {
		// if the same property is used for hostname and port, or if no property is used for port
		wf(`		Host: %v,`, host)
	} else {
		wf(`		Host: conf.FormatHost(%v, %v),`, host, port)
	}
}

func propErrorName(name string) string {
	lower := strings.ToLower(name)
	split := 1
	if len(name) < 2 || name[1] != lower[1] {
		// initial word is an initialism, iterate until start of second word

		for i := range name {
			if name[i] == lower[i] {
				// start of second word
				split = i - 1
				break
			}
		}
	}

	if split < 0 || split >= len(name) {
		// no good split point found, just return the whole name in lower case
		return lower
	}

	return lower[:split] + name[split:]
}
