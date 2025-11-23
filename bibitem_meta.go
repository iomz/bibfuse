package bibfuse

import (
	"fmt"
	"reflect"
)

type bibItemFieldMeta struct {
	index      int
	fieldName  string
	bibtexName string
	defaultVal string
	hasBibtex  bool
	hasDefault bool
}

var (
	bibItemType        = reflect.TypeOf(BibItem{})
	bibItemFieldMetas  []bibItemFieldMeta
	bibItemBibtexIndex map[string]int
)

func init() {
	bibItemBibtexIndex = make(map[string]int)
	for i := 0; i < bibItemType.NumField(); i++ {
		field := bibItemType.Field(i)
		bibtexName := field.Tag.Get("bibtex")
		defaultVal := field.Tag.Get("default")
		meta := bibItemFieldMeta{
			index:      i,
			fieldName:  field.Name,
			bibtexName: bibtexName,
			defaultVal: defaultVal,
			hasBibtex:  bibtexName != "" && bibtexName != "-",
			hasDefault: defaultVal != "" && defaultVal != "-",
		}
		bibItemFieldMetas = append(bibItemFieldMetas, meta)
		if meta.hasBibtex {
			if _, exists := bibItemBibtexIndex[bibtexName]; exists {
				panic(fmt.Sprintf("duplicate bibtex tag %q on BibItem", bibtexName))
			}
			bibItemBibtexIndex[bibtexName] = i
		}
	}
}

func (meta bibItemFieldMeta) valueFrom(v reflect.Value) string {
	return v.Field(meta.index).String()
}

func (meta bibItemFieldMeta) setString(bi *BibItem, value string) {
	v := reflect.ValueOf(bi).Elem().Field(meta.index)
	v.SetString(value)
}
