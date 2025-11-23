package bibfuse

import (
	"fmt"
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
	Metanote       string `default:"" bibtex:"metanote"`
	Note           string `default:"" bibtex:"note"`
	Number         string `default:"" bibtex:"number"`
	Numpages       string `default:"" bibtex:"numpages"`
	Pages          string `default:"" bibtex:"pages"`
	Publisher      string `default:"" bibtex:"publisher"`
	School         string `default:"" bibtex:"school"`
	Series         string `default:"" bibtex:"series"`
	TechreportType string `default:"" bibtex:"type"`
	URL            string `default:"" bibtex:"url"`
	Version        string `default:"" bibtex:"version"`
	Volume         string `default:"" bibtex:"volume"`
	Year           string `default:"" bibtex:"year"`
}

// AllFields returns a string map of all the fileds
func (bi BibItem) AllFields(bio BibItemOption) map[string]string {
	fieldMap := make(map[string]string, len(bibItemFieldMetas))
	value := reflect.ValueOf(bi)
	switch bio {
	case ByBibTexName:
		for _, meta := range bibItemFieldMetas {
			if !meta.hasBibtex {
				continue
			}
			fieldMap[meta.bibtexName] = meta.valueFrom(value)
		}
	default:
		for _, meta := range bibItemFieldMetas {
			fieldMap[meta.fieldName] = meta.valueFrom(value)
		}
	}
	return fieldMap
}

// SetFieldByBibTexName update the field value specified by bibtex field name
func (bi *BibItem) SetFieldByBibTexName(fieldName, fieldValue string) error {
	metaIndex, ok := bibItemBibtexIndex[fieldName]
	if !ok {
		return fmt.Errorf("SetFieldByBibTexName: no such field %q", fieldName)
	}
	bibItemFieldMetas[metaIndex].setString(bi, fieldValue)
	return nil
}

// FieldValueByBibTexName returns the field value specified by bibtex field name.
func (bi BibItem) FieldValueByBibTexName(fieldName string) (string, bool) {
	metaIndex, ok := bibItemBibtexIndex[fieldName]
	if !ok {
		return "", false
	}
	return bibItemFieldMetas[metaIndex].valueFrom(reflect.ValueOf(bi)), true
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
	for _, meta := range bibItemFieldMetas {
		if !meta.hasDefault {
			continue
		}
		meta.setString(&bi, meta.defaultVal)
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

func (fs Filters) filterFor(citeType string) Filter {
	if filter, ok := fs[citeType]; ok {
		return filter
	}
	if filter, ok := fs["default"]; ok {
		return filter
	}
	return Filter{}
}

// BuildBibItem returns BibItem with the filter
func (fs Filters) BuildBibItem(entry *bibtex.BibEntry, smart bool, oneofs Oneofs) (BibItem, error) {
	bi := NewBibItem()
	bi.CiteName = entry.CiteName
	bi.CiteType = entry.Type

	for k, v := range entry.Fields {
		switch k {
		case "author":
			authors, err := NewAuthors(v.String())
			if err != nil {
				return bi, fmt.Errorf("[%v] %w", bi.CiteName, err)
			}
			_ = bi.SetFieldByBibTexName(k, authors.String())
		default:
			_ = bi.SetFieldByBibTexName(k, v.String())
		}
	}

	filter := fs.filterFor(entry.Type)

	for _, fieldName := range filter["todos"] {
		if value, ok := bi.FieldValueByBibTexName(fieldName); ok && value == "" {
			_ = bi.SetFieldByBibTexName(fieldName, "(TODO)")
		}
	}
	for _, fieldName := range filter["optionals"] {
		if value, ok := bi.FieldValueByBibTexName(fieldName); ok && value == "" {
			_ = bi.SetFieldByBibTexName(fieldName, "(OPTIONAL)")
		}
	}

	// smart mode: use oneof_ filters to discard unnecessary fields
	if smart && oneofs.HasOneof(entry.Type) {
		for _, of := range *(oneofs[entry.Type]) {
			var keep string
			for _, fieldName := range of {
				value, ok := bi.FieldValueByBibTexName(fieldName)
				if !ok {
					continue
				}
				if value == "" || value == "(TODO)" || value == "(OPTIONAL)" {
					continue
				}
				keep = fieldName
				break
			}
			if keep == "" {
				continue
			}
			for _, fieldName := range of {
				if fieldName == keep {
					continue
				}
				_ = bi.SetFieldByBibTexName(fieldName, "")
			}
		}
	}

	return bi, nil
}

// Oneof defines the one of the fields
type Oneof [][]string

// AddOneof adds a oneof fields
func (of *Oneof) AddOneof(fields []string) {
	*of = append(*of, fields)
}

// NewOneof initialize a Oneof
func NewOneof() *Oneof {
	return &Oneof{}
}

// Oneofs is a oneof map for each citeType
type Oneofs map[string]*Oneof

// HasOneof checks if the filter of a citation type exists
func (os Oneofs) HasOneof(citeType string) bool {
	_, ok := os[citeType]
	return ok
}
