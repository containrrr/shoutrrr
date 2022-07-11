package ref

import (
	"errors"
	"net/url"

	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/services/standard"
	t "github.com/containrrr/shoutrrr/pkg/types"
)

type testStruct struct {
	Signed          int `key:"signed" default:"0"`
	Unsigned        uint
	Str             string `key:"str" default:"notempty"`
	StrSlice        []string
	StrArray        [3]string
	Sub             subStruct
	TestEnum        int `key:"testenum" default:"None"`
	SubProp         subPropStruct
	SubSlice        []subStruct
	SubPropSlice    []subPropStruct
	SubPropPtrSlice []*subPropStruct
	StrMap          map[string]string
	IntMap          map[string]int
	Int8Map         map[string]int8
	Int16Map        map[string]int16
	Int32Map        map[string]int32
	Int64Map        map[string]int64
	UintMap         map[string]uint
	Uint8Map        map[string]int8
	Uint16Map       map[string]int16
	Uint32Map       map[string]int32
	Uint64Map       map[string]int64
}

func (t *testStruct) GetURL() *url.URL {
	panic("not implemented")
}

func (t *testStruct) SetURL(_ *url.URL) error {
	panic("not implemented")
}

func (t *testStruct) Enums() map[string]t.EnumFormatter {
	return enums
}

type subStruct struct {
	Value string
}

type subPropStruct struct {
	Value string
}

func (s *subPropStruct) SetFromProp(propValue string) error {
	if len(propValue) < 1 || propValue[0] != '@' {
		return errors.New("invalid value")
	}
	s.Value = propValue[1:]
	return nil
}
func (s *subPropStruct) GetPropValue() (string, error) {
	return "@" + s.Value, nil
}

var (
	testEnum = format.CreateEnumFormatter([]string{"None", "Foo", "Bar"})
	enums    = map[string]t.EnumFormatter{
		"TestEnum": testEnum,
	}
)

type testStructBadDefault struct {
	standard.EnumlessConfig
	Value int `key:"value" default:"NaN"`
}

func (t *testStructBadDefault) GetURL() *url.URL {
	panic("not implemented")
}

func (t *testStructBadDefault) SetURL(_ *url.URL) error {
	panic("not implemented")
}
