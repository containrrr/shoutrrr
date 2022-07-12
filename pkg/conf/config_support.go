package conf

import (
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/types"
)

func ParseNumberValue(v string, base int) (int64, error) {
	if v == "" {
		return 0, nil
	}
	return strconv.ParseInt(v, base, 64)
}

func ParseTextValue(v string) (string, error) {
	return v, nil
}

func ParsePathListValue(v string) ([]string, error) {
	v = strings.TrimPrefix(v, "/")
	return ParseListValue(v, "/")
}

func ParseListValue(v string, sep string) ([]string, error) {
	if v == "" {
		// If the source string is empty, the list contains no items
		return []string{}, nil
	}
	return strings.Split(v, sep), nil
}

func ParseColorValue(v string) (uint32, error) {
	if len(v) > 0 && v[0] == '#' {
		v = v[1:]
	}
	if len(v) > 1 && v[:2] == "0x" {
		v = v[2:]
	}
	color, err := strconv.ParseUint(v, 16, 32)
	return uint32(color), err
}

func ParseDurationValue(v string) (time.Duration, error) {
	return time.ParseDuration(v)
}

func ParseToggleValue(v string) (bool, error) {
	if val, ok := format.ParseBool(v, false); ok {
		return val, nil
	} else {
		return false, fmt.Errorf("invalid toggle value %q", v)
	}
}

func ParseMapValue(v string, itemSep string, kvSep string) (map[string]string, error) {
	if v == "" {
		// If the source string is empty, the map contains no items
		return map[string]string{}, nil
	}
	pairs := strings.Split(v, itemSep)
	kvMap := make(map[string]string, len(pairs))
	for _, pair := range pairs {
		kv := strings.SplitN(pair, kvSep, 2)
		if len(kv) < 2 {
			return map[string]string{}, fmt.Errorf("invalid key/value pair %q", pair)
		}
		kvMap[kv[0]] = kv[1]
	}
	return kvMap, nil
}

func FormatNumberValue(v int64, base int) string {
	return strconv.FormatInt(v, base)
}

func FormatTextValue(v string) string {
	return v
}

func FormatListValue(v []string, sep string) string {
	return strings.Join(v, sep)
}

func FormatColorValue(v uint32) string {
	return fmt.Sprintf("0x%06x", v)
}

func FormatDurationValue(v time.Duration) string {
	return v.String()
}

func FormatToggleValue(v bool) string {
	return format.PrintBool(v)
}

func FormatMapValue(v map[string]string, itemSep string, kvSep string) string {
	sb := strings.Builder{}
	keys := make([]string, 0, len(v))
	for key := range v {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for i, key := range keys {
		if i != 0 {
			sb.WriteString(itemSep)
		}
		sb.WriteString(key)
		sb.WriteString(kvSep)
		sb.WriteString(v[key])
	}
	return sb.String()
}

func FormatHost(hostname string, port int64) string {
	return fmt.Sprintf("%v:%v", hostname, port)
}

// UpdateFromParams updates the configuration from the supplied params
func UpdateFromParams(config types.GeneratedConfig, params *types.Params) error {
	if params == nil {
		return nil
	}
	updates := make(map[int]string, len(*params))
	for key, value := range *params {
		if propIndex, found := config.PropInfo().PropIndexFor(key); found {
			updates[propIndex] = value
		} else if key != "title" {
			return fmt.Errorf("invalid key %q", key)
		}
	}
	return config.Update(updates)
}

// UpdateFromQuery updates the configuration from the supplied query values
func UpdateFromQuery(config types.GeneratedConfig, values url.Values) error {
	updates := make(map[int]string, len(values))
	for key, value := range values {
		if propIndex, found := config.PropInfo().PropIndexFor(key); found {
			updates[propIndex] = value[0]
		} else if key != "title" {
			return fmt.Errorf("invalid key %q", key)
		}
	}
	return config.Update(updates)
}

// SetDefaults sets all the Config properties to their default values
func SetDefaults(config types.GeneratedConfig) error {
	info := config.PropInfo()
	updates := make(map[int]string, len(info.PropNames))
	for i := range info.PropNames {
		updates[i] = info.DefaultValues[i]
	}
	return config.Update(updates)
}

func SetPartialDefaults(config types.GeneratedConfig) error {
	info := config.PropInfo()
	updates := make(map[int]string, len(info.PropNames))
	for i := range info.PropNames {
		currValue := config.PropValue(i)
		if currValue == "" || currValue == "0" {
			updates[i] = info.DefaultValues[i]
		}
	}
	return config.Update(updates)
}

// Init sets properties to default values and then updates them according to the configURL
func Init(config types.GeneratedConfig, configURL *url.URL) error {
	if err := SetDefaults(config); err != nil {
		return err
	}

	return config.SetURL(configURL)
}

// QueryValues returns a url.Values populated from the configuration
func QueryValues(config types.GeneratedConfig) url.Values {
	info := config.PropInfo()
	values := make(url.Values, len(info.PropNames))

	for i := range info.PropNames {
		if info.PrimaryKeys[i] < 0 {
			continue
		}
		value := config.PropValue(i)
		if value == info.DefaultValues[i] {
			continue
		}
		values.Set(info.Keys[info.PrimaryKeys[i]], value)
	}

	if cqc, ok := (interface{})(config).(types.CustomQueryConfig); ok {
		for key, value := range cqc.CustomQueryVars() {
			values.Set(EscapeCustomQueryKey(config, key), value[0])
		}
	}

	return values
}

// SplitPath splits a url.Path, removing the initial slash if present
func SplitPath(path string) []string {
	return strings.Split(strings.Trim(path, "/"), "/")
}

// JoinPath joins multiple path elements into a URL path
func JoinPath(parts ...string) string {
	return strings.Join(parts, "/")
}

func ValueMatchesPattern(value, pattern string) bool {
	matches, _ := regexp.MatchString(pattern, value)
	return matches
}

func UnescapeCustomQueryKey(key string) string {
	return strings.TrimPrefix(key, "__")
}

func EscapeCustomQueryKey(config types.GeneratedConfig, key string) string {
	lowerKey := strings.ToLower(key)
	for _, reserved := range config.PropInfo().Keys {
		if lowerKey == reserved {
			return "__" + lowerKey
		}
	}
	return key
}

func UserInfoOrNil(ui *url.Userinfo) *url.Userinfo {
	if strings.TrimPrefix(ui.String(), ":") == "" {
		return nil
	}
	return ui
}
