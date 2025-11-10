package bibfuse

import (
	"bufio"
	"errors"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/nickng/bibtex"
	"github.com/spf13/viper"
)

var itemtests = []struct {
	in  BibItem
	opt BibItemOption
	out map[string]string
}{
	{
		BibItem{"mizutani2021article", "article", "{Title of the Article}", "Mizutani, Iori", "", "(OPTIONAL)", "", "(OPTIONAL)", "(OPTIONAL)", "", "(TODO)", "(OPTIONAL)", "", "(OPTIONAL)", "(OPTIONAL)", "(OPTIONAL)", "(OPTIONAL)", "", "", "", "(OPTIONAL)", "", "(OPTIONAL)", "2021"},
		ByBibTexName,
		map[string]string{"author": "Mizutani, Iori", "booktitle": "", "cite_name": "mizutani2021article", "cite_type": "article", "doi": "(OPTIONAL)", "edition": "", "institution": "", "isbn": "(OPTIONAL)", "issn": "(OPTIONAL)", "journal": "(TODO)", "metanote": "(OPTIONAL)", "note": "", "number": "(OPTIONAL)", "numpages": "(OPTIONAL)", "pages": "(OPTIONAL)", "publisher": "(OPTIONAL)", "school": "", "series": "", "title": "{Title of the Article}", "type": "", "url": "(OPTIONAL)", "version": "", "volume": "(OPTIONAL)", "year": "2021"},
	},
}

func TestAllFields(t *testing.T) {
	for _, tt := range itemtests {
		result := tt.in.AllFields(tt.opt)
		if !reflect.DeepEqual(result, tt.out) {
			t.Errorf("BibItem.AllFields(%v) => \n%v, want \n%v", tt.opt, result, tt.out)
		}
	}
}

var bibtests = []struct {
	in     string
	smart  bool
	oneofs Oneofs
	err    error
	out    string
}{
	{
		"@article{mizutani2021article,\ntitle={{Title of the Article}},\nauthor=\"Mizutani, Iori\",\n}",
		false,
		Oneofs{},
		nil,
		`@article{mizutani2021article,
    title       = {{Title of the Article}},
    author      = "Mizutani, Iori",
    url         = "(OPTIONAL)",
    booktitle   = "",
    doi         = "(OPTIONAL)",
    edition     = "",
    institution = "",
    isbn        = "(OPTIONAL)",
    issn        = "(OPTIONAL)",
    journal     = "(TODO)",
    metanote    = "(OPTIONAL)",
    note        = "",
    number      = "(OPTIONAL)",
    numpages    = "(OPTIONAL)",
    pages       = "(OPTIONAL)",
    publisher   = "(OPTIONAL)",
    school      = "",
    series      = "",
    type        = "",
    version     = "",
    volume      = "(OPTIONAL)",
    year        = "(TODO)",
}
`,
	}, {
		"@book{mizutani2021book,\ntitle={{Title of the Book}},\nauthor=\"Mizutani, Iori\",\n}",
		false,
		Oneofs{},
		nil,
		`@book{mizutani2021book,
    title       = {{Title of the Book}},
    author      = "Mizutani, Iori",
    url         = "(OPTIONAL)",
    booktitle   = "",
    doi         = "(OPTIONAL)",
    edition     = "(OPTIONAL)",
    institution = "",
    isbn        = "(OPTIONAL)",
    issn        = "(OPTIONAL)",
    journal     = "",
    metanote    = "(OPTIONAL)",
    note        = "",
    number      = "",
    numpages    = "",
    pages       = "",
    publisher   = "(TODO)",
    school      = "",
    series      = "",
    type        = "",
    version     = "",
    volume      = "",
    year        = "(TODO)",
}
`,
	}, {
		"@incollection{mizutani2021incollection,\ntitle={{Title of the Book Chapter}},\nauthor=\"Mizutani, Iori\",\n}",
		false,
		Oneofs{},
		nil,
		`@incollection{mizutani2021incollection,
    title       = {{Title of the Book Chapter}},
    author      = "Mizutani, Iori",
    url         = "(OPTIONAL)",
    booktitle   = "(TODO)",
    doi         = "(OPTIONAL)",
    edition     = "",
    institution = "",
    isbn        = "(OPTIONAL)",
    issn        = "(OPTIONAL)",
    journal     = "",
    metanote    = "(OPTIONAL)",
    note        = "",
    number      = "",
    numpages    = "(OPTIONAL)",
    pages       = "(OPTIONAL)",
    publisher   = "(TODO)",
    school      = "",
    series      = "",
    type        = "",
    version     = "",
    volume      = "",
    year        = "(TODO)",
}
`,
	}, {
		"@inproceedings{mizutani2021inproceedings,\ntitle={{Title of the Conference Paper}},\nauthor=\"Mizutani, Iori\",\n}",
		false,
		Oneofs{},
		nil,
		`@inproceedings{mizutani2021inproceedings,
    title       = {{Title of the Conference Paper}},
    author      = "Mizutani, Iori",
    url         = "(OPTIONAL)",
    booktitle   = "(TODO)",
    doi         = "(OPTIONAL)",
    edition     = "",
    institution = "",
    isbn        = "(OPTIONAL)",
    issn        = "(OPTIONAL)",
    journal     = "",
    metanote    = "(OPTIONAL)",
    note        = "",
    number      = "",
    numpages    = "(OPTIONAL)",
    pages       = "(OPTIONAL)",
    publisher   = "(OPTIONAL)",
    school      = "",
    series      = "(OPTIONAL)",
    type        = "",
    version     = "",
    volume      = "",
    year        = "(TODO)",
}
`,
	}, {
		"@mastersthesis{mizutani2021mastersthesis,\ntitle={{Title of the Master's Thesis}},\n}",
		false,
		Oneofs{},
		nil,
		`@mastersthesis{mizutani2021mastersthesis,
    title       = {{Title of the Master's Thesis}},
    author      = "(TODO)",
    url         = "(OPTIONAL)",
    booktitle   = "",
    doi         = "",
    edition     = "",
    institution = "",
    isbn        = "",
    issn        = "",
    journal     = "",
    metanote    = "(OPTIONAL)",
    note        = "",
    number      = "",
    numpages    = "",
    pages       = "",
    publisher   = "",
    school      = "(TODO)",
    series      = "",
    type        = "",
    version     = "",
    volume      = "",
    year        = "(TODO)",
}
`,
	}, {
		"@misc{mizutani2021misc,\ntitle={{Title of the Resource}},\nauthor=\"Mizutani, Iori\",\n}",
		false,
		Oneofs{},
		nil,
		`@misc{mizutani2021misc,
    title       = {{Title of the Resource}},
    author      = "Mizutani, Iori",
    url         = "(TODO)",
    booktitle   = "",
    doi         = "",
    edition     = "",
    institution = "(OPTIONAL)",
    isbn        = "",
    issn        = "",
    journal     = "",
    metanote    = "(OPTIONAL)",
    note        = "(TODO)",
    number      = "",
    numpages    = "",
    pages       = "",
    publisher   = "",
    school      = "",
    series      = "",
    type        = "",
    version     = "",
    volume      = "",
    year        = "(TODO)",
}
`,
	}, {
		"@phdthesis{mizutani2021phdthesis,\ntitle={{Title of the Ph.D. Thesis}},\n}",
		false,
		Oneofs{},
		nil,
		`@phdthesis{mizutani2021phdthesis,
    title       = {{Title of the Ph.D. Thesis}},
    author      = "(TODO)",
    url         = "(OPTIONAL)",
    booktitle   = "",
    doi         = "",
    edition     = "",
    institution = "",
    isbn        = "",
    issn        = "",
    journal     = "",
    metanote    = "(OPTIONAL)",
    note        = "",
    number      = "",
    numpages    = "",
    pages       = "",
    publisher   = "",
    school      = "(TODO)",
    series      = "",
    type        = "",
    version     = "",
    volume      = "",
    year        = "(TODO)",
}
`,
	}, {
		"@techreport{mizutani2021techreport,\ntitle={{Title of the Technical Document}},\nauthor=\"Mizutani, Iori\",\n}",
		false,
		Oneofs{},
		nil,
		`@techreport{mizutani2021techreport,
    title       = {{Title of the Technical Document}},
    author      = "Mizutani, Iori",
    url         = "(OPTIONAL)",
    booktitle   = "",
    doi         = "",
    edition     = "",
    institution = "(TODO)",
    isbn        = "",
    issn        = "",
    journal     = "",
    metanote    = "(OPTIONAL)",
    note        = "",
    number      = "",
    numpages    = "",
    pages       = "",
    publisher   = "",
    school      = "",
    series      = "(OPTIONAL)",
    type        = "",
    version     = "(OPTIONAL)",
    volume      = "",
    year        = "(TODO)",
}
`,
	}, {
		"@unpublished{mizutani2021unpublished,\ntitle={{Title of the Unpublished Work}},\nauthor=\"Mizutani, Iori\",\n}",
		false,
		Oneofs{},
		nil,
		`@unpublished{mizutani2021unpublished,
    title       = {{Title of the Unpublished Work}},
    author      = "Mizutani, Iori",
    url         = "(TODO)",
    booktitle   = "",
    doi         = "",
    edition     = "",
    institution = "",
    isbn        = "",
    issn        = "",
    journal     = "",
    metanote    = "(OPTIONAL)",
    note        = "(TODO)",
    number      = "",
    numpages    = "",
    pages       = "",
    publisher   = "",
    school      = "",
    series      = "",
    type        = "",
    version     = "",
    volume      = "",
    year        = "",
}
`,
	}, {
		"@article{mizutani2021article,\ntitle={{Title of the Article}},\nauthor=\"Mizutani, Iori\",\ndoi={xxxxx/xxxxx.xxx.xx},\nisbn={978-3-16-148410-0},\nissn={2049-3630},\nnumber=1,\nnumpages=350,\npages={1--350},\nvolume=1,\n}",
		true,
		Oneofs{"article": &Oneof{
			[]string{"doi", "pages", "numpages"},
			[]string{"doi", "isbn"},
			[]string{"doi", "issn"},
			[]string{"doi", "number"},
			[]string{"doi", "publisher"},
			[]string{"doi", "volume"},
			[]string{"doi", "url"},
		}},
		nil,
		`@article{mizutani2021article,
    title       = {{Title of the Article}},
    author      = "Mizutani, Iori",
    url         = "",
    booktitle   = "",
    doi         = "xxxxx/xxxxx.xxx.xx",
    edition     = "",
    institution = "",
    isbn        = "",
    issn        = "",
    journal     = "(TODO)",
    metanote    = "(OPTIONAL)",
    note        = "",
    number      = "",
    numpages    = "",
    pages       = "",
    publisher   = "",
    school      = "",
    series      = "",
    type        = "",
    version     = "",
    volume      = "",
    year        = "(TODO)",
}
`,
	}, {
		"@article{mizutani2021article,\ntitle={{Title of the Article}},\nauthor=\"M., Iori\",\n}",
		false,
		Oneofs{},
		errors.New("last name should not be abbreviated"),
		`
`,
	},
}

func TestBuildBibItem(t *testing.T) {
	filters := loadConfig(t)
	for _, tt := range bibtests {
		parsed, err := bibtex.Parse(strings.NewReader(tt.in))
		if err != nil {
			t.Error(err)
		}
		entry := parsed.Entries[0]
		bi, err := filters.BuildBibItem(entry, tt.smart, tt.oneofs)
		if err != nil {
			if !strings.Contains(err.Error(), tt.err.Error()) {
				t.Errorf("BuildBibItem() err => %v\n, want %v", err, tt.err)
			}
			continue
		}
		entry = bi.ToBibEntry()
		bt := bibtex.NewBibTex()
		bt.AddEntry(entry)
		result := bt.PrettyString()
		if result != tt.out {
			t.Errorf("bt.PrettyString() => \n%v, want \n%v", result, tt.out)
			bibEntryEqual(t, tt.out, result)
		}
	}
}

func TestNewBibItem(t *testing.T) {
	bi := NewBibItem()
	for k, v := range bi.AllFields(ByFieldName) {
		if v != "" {
			t.Errorf("BibItem.%v => %v, want %v", k, v, "")
		}
	}
}

var fieldtests = []struct {
	in   BibItem
	args map[string]string
	out  BibItem
	err  error
}{
	{
		NewBibItem(),
		map[string]string{"name": "title", "value": "Title"},
		BibItem{"", "", "Title", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		nil,
	},
}

func TestSetFieldByBibTexName(t *testing.T) {
	for _, tt := range fieldtests {
		err := tt.in.SetFieldByBibTexName(tt.args["name"], tt.args["value"])
		if err != tt.err {
			t.Errorf("err => \n%v, want \n%v", err, tt.err)
		} else if !reflect.DeepEqual(tt.in, tt.out) {
			t.Errorf("BibItem => \n%v, want \n%v", tt.in, tt.out)
		}
	}
}

var bibitemtests = []struct {
	in  BibItem
	out string
}{
	{
		BibItem{"mizutani2021article", "article", "{Title of the Article}", "Mizutani, Iori", "", "(OPTIONAL)", "", "(OPTIONAL)", "(OPTIONAL)", "", "(TODO)", "(OPTIONAL)", "", "(OPTIONAL)", "(OPTIONAL)", "(OPTIONAL)", "(OPTIONAL)", "", "", "", "(OPTIONAL)", "", "(OPTIONAL)", "2021"},
		`@article{mizutani2021article,
    title       = {{Title of the Article}},
    author      = "Mizutani, Iori",
    url         = "(OPTIONAL)",
    booktitle   = "",
    doi         = "(OPTIONAL)",
    edition     = "",
    institution = "",
    isbn        = "(OPTIONAL)",
    issn        = "(OPTIONAL)",
    journal     = "(TODO)",
    metanote    = "(OPTIONAL)",
    note        = "",
    number      = "(OPTIONAL)",
    numpages    = "(OPTIONAL)",
    pages       = "(OPTIONAL)",
    publisher   = "(OPTIONAL)",
    school      = "",
    series      = "",
    type        = "",
    version     = "",
    volume      = "(OPTIONAL)",
    year        = 2021,
}
`,
	},
}

func TestToBibEntry(t *testing.T) {
	for _, tt := range bibitemtests {
		entry := tt.in.ToBibEntry()
		bt := bibtex.NewBibTex()
		bt.AddEntry(entry)
		result := bt.PrettyString()
		if result != tt.out {
			t.Errorf("bt.PrettyString => \n%v, want \n%v", result, tt.out)
		}
	}
}

func TestFilterHasField(t *testing.T) {
	f := NewFilter()
	ok := f.HasField("article", "todos")
	if ok {
		t.Errorf("f.HasFilter => %v, want false", ok)
	}
}

func TestOneof(t *testing.T) {
	of := NewOneof()
	of.AddOneof([]string{"doi", "url"})
	test := &Oneof{[]string{"doi", "url"}}
	if !reflect.DeepEqual(of, test) {
		t.Errorf("Oneof => %v, want %v", of, test)
	}
}

func TestFieldValueByBibTexName(t *testing.T) {
	bi := NewBibItem()
	if err := bi.SetFieldByBibTexName("title", "Example Title"); err != nil {
		t.Fatalf("SetFieldByBibTexName() err => %v, want nil", err)
	}

	val, ok := bi.FieldValueByBibTexName("title")
	if !ok {
		t.Fatalf("FieldValueByBibTexName() ok => %v, want true", ok)
	}
	if val != "Example Title" {
		t.Fatalf("FieldValueByBibTexName() => %v, want %v", val, "Example Title")
	}

	if _, ok := bi.FieldValueByBibTexName("unknown"); ok {
		t.Fatalf("FieldValueByBibTexName() ok => %v, want false", ok)
	}
}

func TestFiltersFilterFor(t *testing.T) {
	defaultFilter := Filter{"todos": []string{"title"}}
	customFilter := Filter{"todos": []string{"author"}}
	filters := Filters{
		"default": defaultFilter,
		"article": customFilter,
	}

	if got := filters.filterFor("article"); !reflect.DeepEqual(got, customFilter) {
		t.Fatalf("filterFor(\"article\") => %v, want %v", got, customFilter)
	}
	if got := filters.filterFor("book"); !reflect.DeepEqual(got, defaultFilter) {
		t.Fatalf("filterFor(\"book\") => %v, want %v", got, defaultFilter)
	}
	empty := (Filters{}).filterFor("unknown")
	if empty == nil || len(empty) != 0 {
		t.Fatalf("filterFor on empty filters => %v, want empty Filter", empty)
	}
}

// bibEntryEqual compares 2 bib entries
func bibEntryEqual(t *testing.T, from, to string) bool {
	fromScanner := bufio.NewScanner(strings.NewReader(from))
	toScanner := bufio.NewScanner(strings.NewReader(to))
	for fromScanner.Scan() || toScanner.Scan() {
		from := fromScanner.Text()
		to := toScanner.Text()
		//t.Logf("%v\n%v", from, to)
		if from != to {
			//t.Errorf("%v\nmust be\n%v", from, to)
			return false
		}
	}
	return true
}

// generateFilters loads up the default filters
func loadConfig(t *testing.T) Filters {
	viper.SetConfigName("bibfuse")
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}
	viper.AddConfigPath(filepath.Dir(filename))

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

	return filters
}
