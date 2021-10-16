package bibfuse

import (
	"github.com/nickng/bibtex"
)

// BibItem holds all the possible fields
type BibItem struct {
	CiteName       string
	CiteType       string
	Title          string
	Author         string
	Booktitle      string
	DOI            string
	Edition        string
	Keyword        string
	Location       string
	ISBN           string
	ISSN           string
	Institution    string
	Journal        string
	Metanote       string
	Note           string
	Number         string
	Numpages       string
	Pages          string
	Publisher      string
	Series         string
	TechreportType string
	URL            string
	Version        string
	Volume         string
	Year           string
}

// ToBibEntry creates a new bibtex.BibEntry from BibItem and return the pointer
func (bi *BibItem) ToBibEntry() *bibtex.BibEntry {
	entry := bibtex.NewBibEntry(bi.CiteType, bi.CiteName)
	fieldMap := map[string]string{
		"title":       bi.Title,
		"author":      bi.Author,
		"booktitle":   bi.Booktitle,
		"doi":         bi.DOI,
		"edition":     bi.Edition,
		"keyword":     bi.Keyword,
		"location":    bi.Location,
		"isbn":        bi.ISBN,
		"issn":        bi.ISSN,
		"institution": bi.Institution,
		"journal":     bi.Journal,
		"metanote":    bi.Metanote,
		"note":        bi.Note,
		"number":      bi.Number,
		"numpages":    bi.Numpages,
		"pages":       bi.Pages,
		"publisher":   bi.Publisher,
		"series":      bi.Series,
		"type":        bi.TechreportType,
		"url":         bi.URL,
		"version":     bi.Version,
		"volume":      bi.Volume,
		"year":        bi.Year,
	}
	for k, v := range fieldMap {
		if v != "" {
			entry.AddField(k, bibtex.NewBibConst(v))
		}
	}
	return entry
}

// Filter has TODO and OPTIONAL fields for a citation type
type Filter map[string][]string

// NewFilter initialize a Filter
func NewFilter() Filter {
	return make(Filter)
}

// Filters have TODO and OPTIONAL fields for each citation type
type Filters map[string]Filter

// HasFilter checks if the filter of a type already exists
func (fs Filters) HasFilter(filterType string) bool {
	_, ok := fs[filterType]
	return ok
}

// Update updates the fields in the given entry with the filter
func (fs Filters) Update(entry *bibtex.BibEntry) {
	for _, f := range fs[entry.Type]["todos"] {
		if _, ok := entry.Fields[f]; !ok {
			entry.AddField(f, bibtex.NewBibConst("(TODO)"))
		}
	}
	for _, f := range fs[entry.Type]["optionals"] {
		if _, ok := entry.Fields[f]; !ok {
			entry.AddField(f, bibtex.NewBibConst("(OPTIONAL)"))
		}
	}
}
