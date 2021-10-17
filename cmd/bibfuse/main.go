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

// createDB creates the db if not exists
func createDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS entries(
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
            keyword TEXT DEFAULT "",
            location TEXT DEFAULT "",
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
        );`,
	)
	return db, err
}

// writeToDB write the BibEntry to the sqlite3 database
func writeToDB(db *sql.DB, bi bibfuse.BibItem) (*sql.Stmt, sql.Result, error) {
	var stmt *sql.Stmt
	var res sql.Result

	tx, err := db.Begin()
	if err != nil {
		return nil, nil, err
	}
	defer tx.Commit()

	stmt, err = tx.Prepare("INSERT INTO entries (cite_name, cite_type, title, author, booktitle, doi, edition, isbn, issn, institution, journal, keyword, location, metanote, note, number, numpages, pages, publisher, school, series, url, type, volume, year) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	res, err = stmt.Exec(bi.CiteName, bi.CiteType, bi.Title, bi.Author,
		bi.Booktitle, bi.DOI, bi.Edition, bi.ISBN, bi.ISSN, bi.Institution, bi.Journal, bi.Keyword, bi.Location, bi.Metanote, bi.Note, bi.Number, bi.Numpages, bi.Pages, bi.Publisher, bi.School, bi.Series, bi.URL, bi.TechreportType, bi.Volume, bi.Year)
	return stmt, res, err
}

func main() {
	conf := flag.String("config", "bibfuse.toml", "The bibfuse.[toml|yml] defining the filters.")
	dbFile := flag.String("db", "bib.db", "The SQLite file to read/write.")
	noOption := flag.Bool("no-optional", false, "Suppress \"OPTIONAL\" fields in the resulting bibtex.")
	noTodo := flag.Bool("no-todo", false, "Suppress \"TODO\" fields in the resulting bibtex.")
	outFile := flag.String("out", "out.bib", "The resulting bibtex to write (it overrides if exists).")
	showEmpty := flag.Bool("show-empty", false, "Suppress empty fields in the resulting bibtex.")
	verbose := flag.Bool("verbose", false, "Print verbose messages.")
	version := flag.Bool("version", false, "Print version.")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s: [options] [.bib ... .bib]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	files := flag.Args()

	// print version
	if *version {
		bi, _ := debug.ReadBuildInfo()
		fmt.Printf("%v\n", bi.Main.Version)
		os.Exit(0)
	}

	// load config
	if *conf != "bibfuse.toml" {
		configPath, err := filepath.Abs(*conf)
		if err != nil {
			panic(err)
		}
		viper.SetConfigFile(configPath)
	} else {
		viper.SetConfigName("bibfuse")
		viper.AddConfigPath(".")
		// add the path to the default config
		_, filename, _, ok := runtime.Caller(0)
		if !ok {
			panic("No caller information")
		}
		viper.AddConfigPath(filepath.Join(filepath.Dir(filename), "../../"))
	}

	// read the config file
	if err := viper.ReadInConfig(); err != nil { // handle errors reading the config file
		log.Fatalf("Fatal error config file: %s \n", err)
	}

	// load the filters
	filters := make(bibfuse.Filters)
	for _, key := range viper.AllKeys() {
		keys := strings.Split(key, ".")
		citeType, filterType := keys[0], keys[1]
		if !filters.HasFilter(citeType) {
			filters[citeType] = bibfuse.NewFilter()
		}
		filters[citeType][filterType] = viper.GetStringSlice(key)
	}

	// create the db
	dbPath := filepath.Join(".", *dbFile)
	db, err := createDB(dbPath)
	defer db.Close()
	if err != nil {
		log.Fatalf("Table creation failed: %q", err)
	}

	// iterate the given files
	newItemCount := 0
	for _, f := range files {
		filePath := filepath.Join(".", f)
		log.Printf("Parsing %v", filePath)
		reader, err := os.Open(filePath)
		if err != nil {
			log.Fatal(err)
		}
		defer reader.Close()
		parsed, err := bibtex.Parse(reader)
		if err != nil {
			log.Fatal(err)
		}

		// inject each entry to the DB
		for _, entry := range parsed.Entries {
			bi := filters.BuildBibItem(entry)
			stmt, res, err := writeToDB(db, bi)
			if stmt != nil {
				defer stmt.Close()
			}
			if err != nil {
				if err.Error() == "UNIQUE constraint failed: entries.cite_name" {
					if *verbose {
						log.Printf("[%s] %q", entry.CiteName, err)
					}
				} else {
					log.Fatalf("[%s] %q", entry.CiteName, err)
				}
			}
			if res != nil {
				newItemCount++
				if *verbose {
					log.Printf("Added %s", entry.CiteName)
				}
			}
		}
	}
	log.Printf("+%v new entries", newItemCount)

	// create a new BibTex to print
	bib := bibtex.NewBibTex()
	rows, err := db.Query("SELECT cite_name, cite_type, title, author, booktitle, doi, edition, isbn, issn, institution, journal, keyword, location, metanote, note, number, numpages, pages, publisher, school, series, type, url, version, volume, year FROM entries ORDER BY cite_name ASC")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		row := bibfuse.NewBibItem()
		err = rows.Scan(&row.CiteName, &row.CiteType, &row.Title, &row.Author, &row.Booktitle, &row.DOI, &row.Edition, &row.ISBN, &row.ISSN, &row.Institution, &row.Journal, &row.Keyword, &row.Location, &row.Metanote, &row.Note, &row.Number, &row.Numpages, &row.Pages, &row.Publisher, &row.School, &row.Series, &row.TechreportType, &row.URL, &row.Version, &row.Volume, &row.Year)
		if err != nil {
			log.Fatal(err)
		}
		entry := row.ToBibEntry()
		bib.AddEntry(entry)
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Printf("%v contains %v entries", dbPath, len(bib.Entries))

	// leave out (OPTIONAL) and (TODO) if the options are given
	outString := bib.PrettyString()
	if *noOption {
		re := regexp.MustCompile("(?m)[\r\n]+^.*(OPTIONAL).*$")
		outString = re.ReplaceAllString(outString, "")
	}
	if *noTodo {
		re := regexp.MustCompile("(?m)[\r\n]+^.*(TODO).*$")
		outString = re.ReplaceAllString(outString, "")
	}
	if !*showEmpty {
		re := regexp.MustCompile("(?m)[\r\n]+^.*\"\".*$")
		outString = re.ReplaceAllString(outString, "")
	}

	// write to a file
	outPath := filepath.Join(".", *outFile)
	writer, err := os.OpenFile(outPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer writer.Close()
	fmt.Fprintf(writer, outString)
}
