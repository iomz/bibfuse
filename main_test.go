package main

import (
	"bufio"
	"strings"
	"testing"

	"github.com/nickng/bibtex"
)

var bibtests = []struct {
	in  string
	out string
}{
	{
		"@book{mizutani2021book,\ntitle={{Title of the Book}},\nauthor=\"Mizutani, Iori\",\n}",
		"@book{mizutani2021book,\n  publisher = {(TODO)},\n  year = {(TODO)},\n  isbn = {(OPTIONAL)},\n  issn = {(OPTIONAL)},\n  url = {(OPTIONAL)},\n  title = {{Title of the Book}},\n  author = {Mizutani, Iori},\n  doi = {(OPTIONAL)},\n  edition = {(OPTIONAL)},\n  metanote = {(OPTIONAL)},\n}",
	},
}

// BibEntryEqual compares 2 bib entries
func BibEntryEqual(t *testing.T, from, to string) bool {
	fromScanner := bufio.NewScanner(strings.NewReader(from))
	for fromScanner.Scan() {
		toScanner := bufio.NewScanner(strings.NewReader(to))
		found := false
		for toScanner.Scan() {
			if strings.HasPrefix(toScanner.Text(), fromScanner.Text()) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func TestUpdateBibEntry(t *testing.T) {
	for _, tt := range bibtests {
		parsed, err := bibtex.Parse(strings.NewReader(tt.in))
		if err != nil {
			t.Error(err)
		}
		entry := parsed.Entries[0]
		updateBibEntry(entry)
		if !BibEntryEqual(t, entry.RawString(), tt.out) {
			t.Errorf("BibEntryTemplate => \n%v, want \n%v", entry.RawString(), tt.out)
		}
	}
}
