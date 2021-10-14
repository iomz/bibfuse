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

	_ "github.com/mattn/go-sqlite3"
	"github.com/nickng/bibtex"
)

// getCitationTypeFilter returns tods and optionals []string
func getCitationTypeFilter(citeType string) ([]string, []string) {
	switch citeType {
	case "article":
		return []string{ // TODOs
				"author",
				"title",
				"journal",
				"year",
			}, []string{ // OPTIONALs
				"doi",
				"isbn",
				"issn",
				"keyword",
				"metanote",
				"number",
				"numpages",
				"pages",
				"publisher",
				"volume",
				"url",
			}
	case "book":
		return []string{ // TODOs
				"author",
				"title",
				"publisher",
				"year",
			}, []string{ // OPTIONALs
				"doi",
				"edition",
				"isbn",
				"issn",
				"metanote",
				"url",
			}
	case "incollection":
		return []string{ // TODOs
				"author",
				"title",
				"booktitle",
				"publisher",
				"year",
			}, []string{ // OPTIONALs
				"doi",
				"keyword",
				"isbn",
				"issn",
				"metanote",
				"numpages",
				"pages",
				"url",
			}
	case "inproceedings":
		return []string{ // TODOs
				"author",
				"title",
				"booktitle",
				"year",
			}, []string{ // OPTIONALs
				"doi",
				"isbn",
				"issn",
				"keyword",
				"location",
				"metanote",
				"numpages",
				"pages",
				"publisher",
				"series",
				"url",
				"year",
			}
	case "misc":
		return []string{ // TODOs
				"author",
				"title",
				"note",
				"url",
				"year",
			}, []string{ // OPTIONALs
				"institution",
				"metanote",
			}
	case "techreport":
		return []string{ // TODOs
				"author",
				"title",
				"institution",
				"year",
			}, []string{ // OPTIONALs
				"metanote",
				"series",
				"url",
				"version",
			}
	default:
		return []string{"Title"}, []string{}
	}
}

// execInsertStatement returns *driver.Stmt of a sqlite3 INSERT query for the specific citation type
func execInsertStatement(db *sql.DB, entry *bibtex.BibEntry) (*sql.Stmt, sql.Result, error) {
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
}

// newBibEntry create a new bibtex.BibEntry and return the pointer
func newBibEntry(citeName, citeType, title, author, booktitle, doi, edition, keyword, location, isbn, issn, institution, journal, metanote, note, number, numpages, pages, publisher, series, techreportType, url, version, volume, year string) *bibtex.BibEntry {
	entry := bibtex.NewBibEntry(citeType, citeName)
	fieldMap := map[string]string{
		"title":       title,
		"author":      author,
		"booktitle":   booktitle,
		"doi":         doi,
		"edition":     edition,
		"keyword":     keyword,
		"location":    location,
		"isbn":        isbn,
		"issn":        issn,
		"institution": institution,
		"journal":     journal,
		"metanote":    metanote,
		"note":        note,
		"number":      number,
		"numpages":    numpages,
		"pages":       pages,
		"publisher":   publisher,
		"series":      series,
		"type":        techreportType,
		"url":         url,
		"version":     version,
		"volume":      volume,
		"year":        year,
	}
	for k, v := range fieldMap {
		if v != "" {
			entry.AddField(k, bibtex.NewBibConst(v))
		}
	}
	return entry
}

// updateBibTexEntryWithFilter replace the fields with (TODO) and (OPTIONAL)
func updateBibTexEntryWithFilter(entry *bibtex.BibEntry, todos, optionals []string) {
	for _, f := range [...]string{"author", "title", "journal", "year"} {
		if _, ok := entry.Fields[f]; !ok {
			entry.AddField(f, bibtex.NewBibConst("(TODO)"))
		}
	}
	for _, f := range [...]string{"doi", "isbn", "issn", "keyword", "metanote", "number", "numpages", "pages", "publisher", "volume", "url"} {
		if _, ok := entry.Fields[f]; !ok {
			entry.AddField(f, bibtex.NewBibConst("(OPTIONAL)"))
		}
	}
}

// writeToDB write the BibEntry to the sqlite3 database
func writeToDB(db *sql.DB, entry *bibtex.BibEntry) (*sql.Stmt, sql.Result, error) {
	todos, options := getCitationTypeFilter(entry.Type)
	for _, f := range todos {
		if _, ok := entry.Fields[f]; !ok {
			entry.AddField(f, bibtex.NewBibConst("(TODO)"))
		}
	}
	for _, f := range options {
		if _, ok := entry.Fields[f]; !ok {
			entry.AddField(f, bibtex.NewBibConst("(OPTIONAL)"))
		}
	}
	stmt, res, err := execInsertStatement(db, entry)
	return stmt, res, err
}

func main() {
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

	// create the db
	dbPath := filepath.Join(".", *dbFile)
	db, err := sql.Open("sqlite3", dbPath)
	defer db.Close()
	if err != nil {
		panic(err)
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
			stmt, res, err := writeToDB(db, entry)
			if stmt != nil {
				defer stmt.Close()
			}
			if err.Error() == "UNIQUE constraint failed: entries.cite_name" {
				if *verbose {
					log.Printf("[%s] %s", entry.CiteName, err)
				}
			} else if err != nil {
				log.Fatalf("[%s] %s", entry.CiteName, err)
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
		var citeName string
		var citeType string
		var title string
		var author string
		var booktitle string
		var doi string
		var edition string
		var keyword string
		var location string
		var isbn string
		var issn string
		var institution string
		var journal string
		var metanote string
		var note string
		var number string
		var numpages string
		var pages string
		var publisher string
		var series string
		var techreportType string
		var url string
		var version string
		var volume string
		var year string
		err = rows.Scan(&citeName, &citeType, &title, &author, &booktitle, &doi, &edition, &keyword, &location, &isbn, &issn, &institution, &journal, &metanote, &note, &number, &numpages, &pages, &publisher, &series, &techreportType, &url, &version, &volume, &year)
		if err != nil {
			log.Fatal(err)
		}
		entry := newBibEntry(citeName, citeType, title, author, booktitle, doi, edition, keyword, location, isbn, issn, institution, journal, metanote, note, number, numpages, pages, publisher, series, techreportType, url, version, volume, year)
		bib.AddEntry(entry)
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}
	}

	// leave out (OPTION) and (TODO) if the options are given
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
