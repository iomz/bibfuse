package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/iomz/bibfuse"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nickng/bibtex"
	"github.com/spf13/viper"
)

const (
	defaultConfigFile = "bibfuse.toml"
	defaultDBFile     = "bib.db"
	defaultOutFile    = "out.bib"
)

var (
	optionalLineRE = regexp.MustCompile("(?m)[\r\n]+^.*(OPTIONAL).*$")
	todoLineRE     = regexp.MustCompile("(?m)[\r\n]+^.*(TODO).*$")
	emptyLineRE    = regexp.MustCompile("(?m)[\r\n]+^.*\"\".*$")
)

const (
	createTableSQL = `CREATE TABLE IF NOT EXISTS entries(
            id INTEGER PRIMARY KEY,
            cite_name TEXT UNIQUE NOT NULL,
            cite_type TEXT NOT NULL,
            author TEXT DEFAULT "",
            title TEXT DEFAULT "",
            booktitle TEXT DEFAULT "",
            doi TEXT DEFAULT "",
            edition TEXT DEFAULT "",
            isbn TEXT DEFAULT "",
            issn TEXT DEFAULT "",
            institution TEXT DEFAULT "",
            journal TEXT DEFAULT "",
            metanote TEXT DEFAULT "",
            note TEXT DEFAULT "",
            number TEXT DEFAULT "",
            numpages TEXT DEFAULT "",
            pages TEXT DEFAULT "",
            publisher TEXT DEFAULT "",
            school TEXT DEFAULT "",
            series TEXT DEFAULT "",
            type TEXT DEFAULT "",
            url TEXT DEFAULT "",
            version TEXT DEFAULT "",
            volume TEXT DEFAULT "",
            year TEXT
        );`
	insertEntrySQL = `INSERT OR IGNORE INTO entries (
            cite_name, cite_type, title, author, booktitle, doi, edition, isbn, issn,
            institution, journal, metanote, note, number, numpages, pages, publisher,
            school, series, url, type, version, volume, year
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
)

type options struct {
	config           string
	useDefaultConfig bool
	dbFile           string
	outFile          string
	noOptional       bool
	noTodo           bool
	showEmpty        bool
	smart            bool
	verbose          bool
	showVersion      bool
}

func main() {
	opts, files := parseFlags()
	if opts.showVersion {
		printVersion()
		return
	}
	if err := run(opts, files); err != nil {
		log.Fatal(err)
	}
}

func parseFlags() (options, []string) {
	opts := options{}

	conf := flag.String("config", defaultConfigFile, "The bibfuse.[toml|yml] defining the filters.")
	dbFile := flag.String("db", defaultDBFile, "The SQLite file to read/write.")
	noOption := flag.Bool("no-optional", false, "Suppress \"OPTIONAL\" fields in the resulting bibtex.")
	noTodo := flag.Bool("no-todo", false, "Suppress \"TODO\" fields in the resulting bibtex.")
	outFile := flag.String("out", defaultOutFile, "The resulting bibtex to write (it overrides if exists).")
	showEmpty := flag.Bool("show-empty", false, "Do not hide empty fields in the resulting bibtex.")
	smart := flag.Bool("smart", false, "Use oneof selectively filters when importing bibtex.")
	verbose := flag.Bool("verbose", false, "Print verbose messages.")
	version := flag.Bool("version", false, "Print version.")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s: [options] [.bib ... .bib]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	opts.config = *conf
	opts.useDefaultConfig = *conf == defaultConfigFile
	opts.dbFile = *dbFile
	opts.outFile = *outFile
	opts.noOptional = *noOption
	opts.noTodo = *noTodo
	opts.showEmpty = *showEmpty
	opts.smart = *smart
	opts.verbose = *verbose
	opts.showVersion = *version

	return opts, flag.Args()
}

func printVersion() {
	bi, _ := debug.ReadBuildInfo()
	fmt.Printf("%v\n", bi.Main.Version)
}

func run(opts options, files []string) error {
	if err := configureViper(opts); err != nil {
		return err
	}

	filters, oneofs, err := loadRules()
	if err != nil {
		return err
	}

	dbPath := filepath.Join(".", opts.dbFile)
	db, err := createDB(dbPath)
	if err != nil {
		return fmt.Errorf("table creation failed: %w", err)
	}
	defer db.Close()

	newItemCount, err := importBibFiles(db, filters, oneofs, opts, files)
	if err != nil {
		return err
	}
	log.Printf("+%d new entries", newItemCount)

	content, entryCount, err := exportBibliography(db, opts)
	if err != nil {
		return err
	}
	log.Printf("%s contains %d entries", dbPath, entryCount)

	outPath := filepath.Join(".", opts.outFile)
	if err := os.WriteFile(outPath, []byte(content), 0o644); err != nil {
		return err
	}
	log.Printf("%d entries written to %s", entryCount, outPath)

	return nil
}

func configureViper(opts options) error {
	if !opts.useDefaultConfig {
		configPath, err := filepath.Abs(opts.config)
		if err != nil {
			return err
		}
		viper.SetConfigFile(configPath)
		return nil
	}

	viper.SetConfigName("bibfuse")
	viper.AddConfigPath(".")

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("no caller information")
	}
	viper.AddConfigPath(filepath.Join(filepath.Dir(filename), "../../"))
	return nil
}

func loadRules() (bibfuse.Filters, bibfuse.Oneofs, error) {
	if err := viper.ReadInConfig(); err != nil {
		return nil, nil, fmt.Errorf("config: %w", err)
	}

	filters := make(bibfuse.Filters)
	oneofs := make(bibfuse.Oneofs)

	for _, key := range viper.AllKeys() {
		keys := strings.Split(key, ".")
		if len(keys) < 2 {
			continue
		}
		citeType, filterType := keys[0], keys[1]
		switch {
		case filterType == "todos" || filterType == "optionals":
			if !filters.HasFilter(citeType) {
				filters[citeType] = bibfuse.NewFilter()
			}
			filters[citeType][filterType] = viper.GetStringSlice(key)
		case strings.HasPrefix(filterType, "oneof_"):
			if !oneofs.HasOneof(citeType) {
				oneofs[citeType] = bibfuse.NewOneof()
			}
			oneofs[citeType].AddOneof(viper.GetStringSlice(key))
		}
	}
	return filters, oneofs, nil
}

func createDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(createTableSQL); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func importBibFiles(db *sql.DB, filters bibfuse.Filters, oneofs bibfuse.Oneofs, opts options, files []string) (int, error) {
	newItemCount := 0
	for _, fileName := range files {
		filePath := filepath.Join(".", fileName)
		log.Printf("parsing %s", filePath)

		reader, err := os.Open(filePath)
		if err != nil {
			return newItemCount, err
		}

		parsed, err := bibtex.Parse(reader)
		reader.Close()
		if err != nil {
			return newItemCount, err
		}

		for _, entry := range parsed.Entries {
			bi, err := filters.BuildBibItem(entry, opts.smart, oneofs)
			if err != nil {
				log.Println(err)
				continue
			}

			added, err := insertEntry(db, bi)
			if err != nil {
				return newItemCount, fmt.Errorf("[%s] %w", entry.CiteName, err)
			}

			if added {
				newItemCount++
				if opts.verbose {
					log.Printf("added %s", entry.CiteName)
				}
			} else if opts.verbose {
				log.Printf("[%s] duplicate entry", entry.CiteName)
			}
		}
	}
	return newItemCount, nil
}

func insertEntry(db *sql.DB, bi bibfuse.BibItem) (bool, error) {
	res, err := db.Exec(
		insertEntrySQL,
		bi.CiteName,
		bi.CiteType,
		bi.Title,
		bi.Author,
		bi.Booktitle,
		bi.DOI,
		bi.Edition,
		bi.ISBN,
		bi.ISSN,
		bi.Institution,
		bi.Journal,
		bi.Metanote,
		bi.Note,
		bi.Number,
		bi.Numpages,
		bi.Pages,
		bi.Publisher,
		bi.School,
		bi.Series,
		bi.URL,
		bi.TechreportType,
		bi.Version,
		bi.Volume,
		bi.Year,
	)
	if err != nil {
		return false, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return affected > 0, nil
}

func exportBibliography(db *sql.DB, opts options) (string, int, error) {
	rows, err := db.Query(`SELECT cite_name, cite_type, title, author, booktitle, doi, edition, isbn, issn, institution, journal, metanote, note, number, numpages, pages, publisher, school, series, type, url, version, volume, year FROM entries ORDER BY cite_name ASC`)
	if err != nil {
		return "", 0, err
	}
	defer rows.Close()

	bib := bibtex.NewBibTex()

	for rows.Next() {
		row := bibfuse.NewBibItem()
		if err := rows.Scan(
			&row.CiteName,
			&row.CiteType,
			&row.Title,
			&row.Author,
			&row.Booktitle,
			&row.DOI,
			&row.Edition,
			&row.ISBN,
			&row.ISSN,
			&row.Institution,
			&row.Journal,
			&row.Metanote,
			&row.Note,
			&row.Number,
			&row.Numpages,
			&row.Pages,
			&row.Publisher,
			&row.School,
			&row.Series,
			&row.TechreportType,
			&row.URL,
			&row.Version,
			&row.Volume,
			&row.Year,
		); err != nil {
			return "", 0, err
		}
		entry := row.ToBibEntry()
		bib.AddEntry(entry)
	}

	if err := rows.Err(); err != nil {
		return "", 0, err
	}

	outString := applyOutputFilters(bib.PrettyString(), opts)
	outString = bibfuse.BackslashCleaner(outString)
	return outString, len(bib.Entries), nil
}

func applyOutputFilters(input string, opts options) string {
	out := input
	if opts.noOptional {
		out = optionalLineRE.ReplaceAllString(out, "")
	}
	if opts.noTodo {
		out = todoLineRE.ReplaceAllString(out, "")
	}
	if !opts.showEmpty {
		out = emptyLineRE.ReplaceAllString(out, "")
	}
	return out
}
