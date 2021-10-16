package bibfuse

import (
	"bufio"
	"os"
	"strings"
	"testing"

	"github.com/nickng/bibtex"
	"github.com/spf13/viper"
)

var bibtests = []struct {
	in  string
	out string
}{
	{
		"@article{mizutani2021article,\ntitle={{Title of the Article}},\nauthor=\"Mizutani, Iori\",\n}",
		`@article{mizutani2021article,
    title     = {{Title of the Article}},
    author    = "Mizutani, Iori",
    url       = "(OPTIONAL)",
    doi       = "(OPTIONAL)",
    isbn      = "(OPTIONAL)",
    issn      = "(OPTIONAL)",
    journal   = "(TODO)",
    keyword   = "(OPTIONAL)",
    metanote  = "(OPTIONAL)",
    number    = "(OPTIONAL)",
    numpages  = "(OPTIONAL)",
    pages     = "(OPTIONAL)",
    publisher = "(OPTIONAL)",
    volume    = "(OPTIONAL)",
    year      = "(TODO)",
}
`,
	}, {
		"@book{mizutani2021book,\ntitle={{Title of the Book}},\nauthor=\"Mizutani, Iori\",\n}",
		`@book{mizutani2021book,
    title     = {{Title of the Book}},
    author    = "Mizutani, Iori",
    url       = "(OPTIONAL)",
    doi       = "(OPTIONAL)",
    edition   = "(OPTIONAL)",
    isbn      = "(OPTIONAL)",
    issn      = "(OPTIONAL)",
    metanote  = "(OPTIONAL)",
    publisher = "(TODO)",
    year      = "(TODO)",
}
`,
	}, {
		"@incollection{mizutani2021incollection,\ntitle={{Title of the Book Chapter}},\nauthor=\"Mizutani, Iori\",\n}",
		`@incollection{mizutani2021incollection,
    title     = {{Title of the Book Chapter}},
    author    = "Mizutani, Iori",
    url       = "(OPTIONAL)",
    booktitle = "(TODO)",
    doi       = "(OPTIONAL)",
    isbn      = "(OPTIONAL)",
    issn      = "(OPTIONAL)",
    keyword   = "(OPTIONAL)",
    metanote  = "(OPTIONAL)",
    numpages  = "(OPTIONAL)",
    pages     = "(OPTIONAL)",
    publisher = "(TODO)",
    year      = "(TODO)",
}
`,
	}, {
		"@inproceedings{mizutani2021inproceedings,\ntitle={{Title of the Conference Paper}},\nauthor=\"Mizutani, Iori\",\n}",
		`@inproceedings{mizutani2021inproceedings,
    title     = {{Title of the Conference Paper}},
    author    = "Mizutani, Iori",
    url       = "(OPTIONAL)",
    booktitle = "(TODO)",
    doi       = "(OPTIONAL)",
    isbn      = "(OPTIONAL)",
    issn      = "(OPTIONAL)",
    keyword   = "(OPTIONAL)",
    location  = "(OPTIONAL)",
    metanote  = "(OPTIONAL)",
    numpages  = "(OPTIONAL)",
    pages     = "(OPTIONAL)",
    publisher = "(OPTIONAL)",
    series    = "(OPTIONAL)",
    year      = "(TODO)",
}
`,
	}, {
		"@misc{mizutani2021misc,\ntitle={{Title of the Resource}},\nauthor=\"Mizutani, Iori\",\n}",
		`@misc{mizutani2021misc,
    title       = {{Title of the Resource}},
    author      = "Mizutani, Iori",
    url         = "(TODO)",
    institution = "(OPTIONAL)",
    metanote    = "(OPTIONAL)",
    note        = "(TODO)",
    year        = "(TODO)",
}
`,
	}, {
		"@techreport{mizutani2021techreport,\ntitle={{Title of the Technical Document}},\nauthor=\"Mizutani, Iori\",\n}",
		`@techreport{mizutani2021techreport,
    title       = {{Title of the Technical Document}},
    author      = "Mizutani, Iori",
    url         = "(OPTIONAL)",
    institution = "(TODO)",
    metanote    = "(OPTIONAL)",
    series      = "(OPTIONAL)",
    version     = "(OPTIONAL)",
    year        = "(TODO)",
}
`,
	},
}

// BibEntryEqual compares 2 bib entries
func BibEntryEqual(t *testing.T, from, to string) bool {
	fromScanner := bufio.NewScanner(strings.NewReader(from))
	toScanner := bufio.NewScanner(strings.NewReader(to))
	for fromScanner.Scan() || toScanner.Scan() {
		from := fromScanner.Text()
		to := toScanner.Text()
		t.Logf("%v\n%v", from, to)
		if from != to {
			t.Errorf("%v\nmust be\n%v", from, to)
			return false
		}
	}
	return true
}

func TestFiltersUpdate(t *testing.T) {
	cwd, _ := os.Getwd()
	viper.SetConfigName("bibfuse")
	viper.AddConfigPath(cwd)
	if err := viper.ReadInConfig(); err != nil { // handle errors reading the config file
		t.Errorf("Fatal error config file: %s \n", err)
	}

	// load the filters
	filters := make(Filters)
	for _, key := range viper.AllKeys() {
		keys := strings.Split(key, ".")
		citeType, filterType := keys[0], keys[1]
		if !filters.HasFilter(citeType) {
			filters[citeType] = NewFilter()
		}
		filters[citeType][filterType] = viper.GetStringSlice(key)
	}

	for _, tt := range bibtests {
		parsed, err := bibtex.Parse(strings.NewReader(tt.in))
		if err != nil {
			t.Error(err)
		}
		entry := parsed.Entries[0]
		filters.Update(entry)
		bt := bibtex.NewBibTex()
		bt.AddEntry(entry)
		result := bt.PrettyString()
		if result != tt.out {
			t.Errorf("BibEntryTemplate => \n%v, want \n%v", result, tt.out)
			BibEntryEqual(t, tt.out, result)
		}
	}
}
