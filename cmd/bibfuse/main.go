package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
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
            keyword TEXT DEFAULT "",
            location TEXT DEFAULT "",
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
func writeToDB(db *sql.DB, entry *bibtex.BibEntry) (*sql.Stmt, sql.Result, error) {
	var stmt *sql.Stmt
	var res sql.Result

	tx, err := db.Begin()
	if err != nil {
		return nil, nil, err
	}
	defer tx.Commit()

	switch entry.Type {
	case "article":
		stmt, err = tx.Prepare("INSERT INTO entries (cite_name, cite_type, title, author, doi, isbn, issn, journal, keyword, metanote, number, numpages, pages, publisher, url, volume, year) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		res, err = stmt.Exec(entry.CiteName,
			entry.Type,
			entry.Fields["title"],
			entry.Fields["author"],
			entry.Fields["doi"],
			entry.Fields["isbn"],
			entry.Fields["issn"],
			entry.Fields["journal"],
			entry.Fields["keyword"],
			entry.Fields["metanote"],
			entry.Fields["number"],
			entry.Fields["numpages"],
			entry.Fields["pages"],
			entry.Fields["publisher"],
			entry.Fields["url"],
			entry.Fields["volume"],
			entry.Fields["year"],
		)
	case "book":
		stmt, err = tx.Prepare("INSERT INTO entries (cite_name, cite_type, title, author, doi, edition, isbn, issn, metanote, publisher, url, year) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		res, err = stmt.Exec(entry.CiteName,
			entry.Type,
			entry.Fields["title"],
			entry.Fields["author"],
			entry.Fields["doi"],
			entry.Fields["edition"],
			entry.Fields["isbn"],
			entry.Fields["issn"],
			entry.Fields["metanote"],
			entry.Fields["publisher"],
			entry.Fields["url"],
			entry.Fields["year"],
		)
	case "incollection":
		stmt, err = tx.Prepare("INSERT INTO entries (cite_name, cite_type, title, author, booktitle, doi, isbn, issn, keyword, metanote, numpages, pages, publisher, url, year) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		res, err = stmt.Exec(entry.CiteName,
			entry.Type,
			entry.Fields["title"],
			entry.Fields["author"],
			entry.Fields["booktitle"],
			entry.Fields["doi"],
			entry.Fields["isbn"],
			entry.Fields["issn"],
			entry.Fields["keyword"],
			entry.Fields["metanote"],
			entry.Fields["numpages"],
			entry.Fields["pages"],
			entry.Fields["publisher"],
			entry.Fields["url"],
			entry.Fields["year"],
		)
	case "inproceedings":
		stmt, err = tx.Prepare("INSERT INTO entries (cite_name, cite_type, title, author, booktitle, doi, isbn, issn, keyword, location, metanote, numpages, pages, publisher, series, url, year) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		res, err = stmt.Exec(entry.CiteName,
			entry.Type,
			entry.Fields["title"],
			entry.Fields["author"],
			entry.Fields["booktitle"],
			entry.Fields["doi"],
			entry.Fields["isbn"],
			entry.Fields["issn"],
			entry.Fields["keyword"],
			entry.Fields["location"],
			entry.Fields["metanote"],
			entry.Fields["numpages"],
			entry.Fields["pages"],
			entry.Fields["publisher"],
			entry.Fields["series"],
			entry.Fields["url"],
			entry.Fields["year"],
		)
	case "misc":
		stmt, err = tx.Prepare("INSERT INTO entries (cite_name, cite_type, title, author, institution, metanote, note, url, year) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		res, err = stmt.Exec(entry.CiteName,
			entry.Type,
			entry.Fields["title"],
			entry.Fields["author"],
			entry.Fields["institution"],
			entry.Fields["metanote"],
			entry.Fields["note"],
			entry.Fields["url"],
			entry.Fields["year"],
		)
	case "techreport":
		stmt, err = tx.Prepare("INSERT INTO entries (cite_name, cite_type, title, author, institution, metanote, series, url, version, year) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		res, err = stmt.Exec(entry.CiteName,
			entry.Type,
			entry.Fields["title"],
			entry.Fields["author"],
			entry.Fields["institution"],
			entry.Fields["metanote"],
			entry.Fields["series"],
			entry.Fields["url"],
			entry.Fields["version"],
			entry.Fields["year"],
		)
	}
	return stmt, res, err

	return stmt, res, err
}

func main() {
	conf := flag.String("c", "bibfuse.toml", "The config.[toml|yml] to use.")
	dbFile := flag.String("db", "bib.db", "The SQLite file to read/write.")
	noOption := flag.Bool("no-optional", false, "Suppress \"OPTIONAL\" fields in the resulting bibtex.")
	noTodo := flag.Bool("no-todo", false, "Suppress \"TODO\" fields in the resulting bibtex.")
	outFile := flag.String("out", "out.bib", "The resulting bibtex to write (it overrides if exists).")
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
		cwd, _ := os.Getwd()
		viper.AddConfigPath(cwd)
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
			filters.Update(entry)
			stmt, res, err := writeToDB(db, entry)
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
				log.Printf("Added %s", entry.CiteName)
			}
		}

	}

	// create a new BibTex to print
	bib := bibtex.NewBibTex()
	rows, err := db.Query("SELECT cite_name, cite_type, title, author, booktitle, doi, edition, keyword, location, isbn, issn, institution, journal, metanote, note, number, numpages, pages, publisher, series, type, url, version, volume, year FROM entries ORDER BY cite_name ASC")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var row bibfuse.BibItem
		err = rows.Scan(&row.CiteName, &row.CiteType, &row.Title, &row.Author, &row.Booktitle, &row.DOI, &row.Edition, &row.Keyword, &row.Location, &row.ISBN, &row.ISSN, &row.Institution, &row.Journal, &row.Metanote, &row.Note, &row.Number, &row.Numpages, &row.Pages, &row.Publisher, &row.Series, &row.TechreportType, &row.URL, &row.Version, &row.Volume, &row.Year)
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

	// write to a file
	outPath := filepath.Join(".", *outFile)
	writer, err := os.OpenFile(outPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer writer.Close()
	fmt.Fprintf(writer, outString)
}
