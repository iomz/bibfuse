package bibfuse

import (
	"errors"
	"reflect"

	"github.com/nickng/bibtex"
)

// BibItemOption specifies the option for AllFields()
type BibItemOption int64

const (
	// ByFieldName use filed names as keys
	ByFieldName BibItemOption = iota
	// ByBibTexName use bibtex names as keys
	ByBibTexName
)

// BibItem holds all the possible fields
type BibItem struct {
	CiteName       string `default:"" bibtex:"cite_name"`
	CiteType       string `default:"" bibtex:"cite_type"`
	Title          string `default:"" bibtex:"title"`
	Author         string `default:"" bibtex:"author"`
	Booktitle      string `default:"" bibtex:"booktitle"`
	DOI            string `default:"" bibtex:"doi"`
	Edition        string `default:"" bibtex:"edition"`
	ISBN           string `default:"" bibtex:"isbn"`
	ISSN           string `default:"" bibtex:"issn"`
	Institution    string `default:"" bibtex:"institution"`
	Journal        string `default:"" bibtex:"journal"`
	Keyword        string `default:"" bibtex:"keyword"`
	Location       string `default:"" bibtex:"location"`
	Metanote       string `default:"" bibtex:"metanote"`
	Note           string `default:"" bibtex:"note"`
	Number         string `default:"" bibtex:"number"`
	Numpages       string `default:"" bibtex:"numpages"`
	Pages          string `default:"" bibtex:"pages"`
	Publisher      string `default:"" bibtex:"publisher"`
	Series         string `default:"" bibtex:"series"`
	TechreportType string `default:"" bibtex:"type"`
	URL            string `default:"" bibtex:"url"`
	Version        string `default:"" bibtex:"version"`
	Volume         string `default:"" bibtex:"volume"`
	Year           string `default:"" bibtex:"year"`
}

// AllFields returns a string map of all the fileds
func (bi BibItem) AllFields(bio BibItemOption) map[string]string {
	fieldMap := make(map[string]string)
	switch bio {
	case ByBibTexName:
		v := reflect.ValueOf(&bi).Elem()
		t := v.Type()
		for i := 0; i < t.NumField(); i++ {
			if fieldName := t.Field(i).Tag.Get("bibtex"); fieldName != "-" {
				fieldMap[fieldName] = v.Field(i).String()
			}
		}
	case ByFieldName:
		v := reflect.ValueOf(bi)
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			fieldMap[t.Field(i).Name] = v.Field(i).String()
		}
	}
	return fieldMap
}

// SetFieldByBibTexName update the field value specified by bibtex field name
func (bi *BibItem) SetFieldByBibTexName(fieldName, fieldValue string) error {
	v := reflect.ValueOf(bi).Elem()
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		if fieldName == t.Field(i).Tag.Get("bibtex") {
			field := v.Field(i)
			field.Set(reflect.ValueOf(fieldValue).Convert(field.Type()))
			return nil
		}
	}
	return errors.New("SetFieldByBibTexName: no such field")
}

// ToBibEntry creates a new bibtex.BibEntry from BibItem and return the pointer
func (bi BibItem) ToBibEntry() *bibtex.BibEntry {
	entry := bibtex.NewBibEntry(bi.CiteType, bi.CiteName)
	for k, v := range bi.AllFields(ByBibTexName) {
		if k != "cite_name" && k != "cite_type" {
			entry.AddField(k, bibtex.NewBibConst(v))
		}
	}
	return entry
}

// NewBibItem returns a BibItem with default field value
func NewBibItem() BibItem {
	bi := BibItem{}
	v := reflect.ValueOf(&bi).Elem()
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		if defaultVal := t.Field(i).Tag.Get("default"); defaultVal != "-" {
			field := v.Field(i).String()
			bi.SetFieldByBibTexName(field, defaultVal)
		}
	}
	return bi
}

// Filter has TODO and OPTIONAL fields for a citation type
type Filter map[string][]string

// HasField checks if the field exists in a given field type (e.g., TODO or OPTIONAL)
func (f Filter) HasField(fieldType, fieldName string) bool {
	_, ok := f[fieldType]
	if !ok {
		return false
	}
	for _, name := range f[fieldType] {
		if name == fieldName {
			return true
		}
	}
	return false
}

// NewFilter initialize a Filter
func NewFilter() Filter {
	return make(Filter)
}

// Filters have TODO and OPTIONAL fields for each citation type
type Filters map[string]Filter

// HasFilter checks if the filter of a citation type exists
func (fs Filters) HasFilter(filterType string) bool {
	_, ok := fs[filterType]
	return ok
}

// ConvertFromBibEntryToBibItem returns BibItem with the filter
func (fs Filters) ConvertFromBibEntryToBibItem(entry *bibtex.BibEntry) BibItem {
	bi := NewBibItem()
	bi.CiteName = entry.CiteName
	bi.CiteType = entry.Type

	for k, v := range entry.Fields {
		bi.SetFieldByBibTexName(k, v.String())
	}

	for k, v := range bi.AllFields(ByBibTexName) {
		if v == "" {
			if fs[entry.Type].HasField("todos", k) {
				bi.SetFieldByBibTexName(k, "(TODO)")
			} else if fs[entry.Type].HasField("optionals", k) {
				bi.SetFieldByBibTexName(k, "(OPTIONAL)")
			}
		}
	}

	return bi
}
