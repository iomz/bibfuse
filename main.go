package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nickng/bibtex"
)

var (
	infile = flag.String("in", "", "Input file (default: stdin)")
	reader = os.Stdin
)

// writeToDB write the BibEntry to the sqlite3 database
func writeToDB(db *sql.DB, entry *bibtex.BibEntry) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit()

	switch entry.Type {
	case "article":
		stmt, err := tx.Prepare("INSERT INTO entries (cite_name, cite_type, author, title, doi, isbn, issn, journal, keyword, metanote, number, numpages, pages, publisher, url, volume, year) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		defer stmt.Close()
		if err != nil {
			return err
		}
		for _, f := range [...]string{"author", "title", "journal", "number", "publisher", "volume", "year"} {
			if _, ok := entry.Fields[f]; !ok {
				entry.AddField(f, bibtex.NewBibConst("(TODO)"))
			}
		}
		for _, f := range [...]string{"doi", "keyword", "isbn", "issn", "metanote", "numpages", "pages", "url"} {
			if _, ok := entry.Fields[f]; !ok {
				entry.AddField(f, bibtex.NewBibConst("(OPTIONAL)"))
			}
		}
		_, err = stmt.Exec(entry.CiteName,
			entry.Type,
			entry.Fields["author"],
			entry.Fields["title"],
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
		stmt, err := tx.Prepare("INSERT INTO entries (cite_name, cite_type, author, title, address, edition, doi, isbn, issn, metanote, publisher, url, year) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		defer stmt.Close()
		if err != nil {
			return err
		}
		for _, f := range [...]string{"author", "title", "address", "edition", "publisher", "year"} {
			if _, ok := entry.Fields[f]; !ok {
				entry.AddField(f, bibtex.NewBibConst("(TODO)"))
			}
		}
		for _, f := range [...]string{"doi", "isbn", "issn", "metanote", "url"} {
			if _, ok := entry.Fields[f]; !ok {
				entry.AddField(f, bibtex.NewBibConst("(OPTIONAL)"))
			}
		}
		_, err = stmt.Exec()
		_, err = stmt.Exec(entry.CiteName,
			entry.Type,
			entry.Fields["author"],
			entry.Fields["title"],
			entry.Fields["address"],
			entry.Fields["edition"],
			entry.Fields["doi"],
			entry.Fields["isbn"],
			entry.Fields["issn"],
			entry.Fields["metanote"],
			entry.Fields["publisher"],
			entry.Fields["url"],
			entry.Fields["year"],
		)
	case "inproceedings":
		stmt, err := tx.Prepare("INSERT INTO entries (cite_name, cite_type, author, title, booktitle, doi, isbn, issn, keyword, location, metanote, numpages, pages, publisher, series, url, year) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		defer stmt.Close()
		if err != nil {
			return err
		}
		for _, f := range [...]string{"author", "title", "booktitle", "publisher", "year"} {
			if _, ok := entry.Fields[f]; !ok {
				entry.AddField(f, bibtex.NewBibConst("(TODO)"))
			}
		}
		for _, f := range [...]string{"doi", "isbn", "issn", "keyword", "location", "metanote", "numpages", "pages", "series", "url", "year"} {
			if _, ok := entry.Fields[f]; !ok {
				entry.AddField(f, bibtex.NewBibConst("(OPTIONAL)"))
			}
		}
		_, err = stmt.Exec(entry.CiteName,
			entry.Type,
			entry.Fields["author"],
			entry.Fields["title"],
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
		stmt, err := tx.Prepare("INSERT INTO entries (cite_name, cite_type, author, title, institution, metanote, month, note, url, year) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		defer stmt.Close()
		if err != nil {
			return err
		}
		for _, f := range [...]string{"author", "title", "institution", "note", "url", "year"} {
			if _, ok := entry.Fields[f]; !ok {
				entry.AddField(f, bibtex.NewBibConst("(TODO)"))
			}
		}
		for _, f := range [...]string{"metanote", "month"} {
			if _, ok := entry.Fields[f]; !ok {
				entry.AddField(f, bibtex.NewBibConst("(OPTIONAL)"))
			}
		}
		_, err = stmt.Exec(entry.CiteName,
			entry.Type,
			entry.Fields["author"],
			entry.Fields["title"],
			entry.Fields["institution"],
			entry.Fields["metanote"],
			entry.Fields["month"],
			entry.Fields["note"],
			entry.Fields["url"],
			entry.Fields["year"],
		)
	case "techreport":
		stmt, err := tx.Prepare("INSERT INTO entries (cite_name, cite_type, author, title, institution, metanote, month, version, series, url, year) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		defer stmt.Close()
		if err != nil {
			return err
		}
		for _, f := range [...]string{"author", "title", "institution", "url", "year"} {
			if _, ok := entry.Fields[f]; !ok {
				entry.AddField(f, bibtex.NewBibConst("(TODO)"))
			}
		}
		for _, f := range [...]string{"metanote", "month", "version", "series"} {
			if _, ok := entry.Fields[f]; !ok {
				entry.AddField(f, bibtex.NewBibConst("(OPTIONAL)"))
			}
		}
		_, err = stmt.Exec(entry.CiteName,
			entry.Type,
			entry.Fields["author"],
			entry.Fields["title"],
			entry.Fields["institution"],
			entry.Fields["metanote"],
			entry.Fields["month"],
			entry.Fields["version"],
			entry.Fields["series"],
			entry.Fields["url"],
			entry.Fields["year"],
		)
	}
	if err != nil {
		return err
	}
	return nil
}

func main() {
	flag.Parse()
	if *infile != "" || len(flag.Args()) > 0 {
		if len(flag.Args()) > 0 {
			*infile = flag.Arg(0)
		}
		rdFile, err := os.Open(*infile)
		if err != nil {
			log.Fatal(err)
		}
		defer rdFile.Close()
		reader = rdFile
	}

	// create the db
	dbPath := filepath.Join(".", "bib.db")
	//log.Printf("Preparing the database: %s", dbPath)
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
            address TEXT DEFAULT "",
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
            month TEXT DEFAULT "",
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

	parsed, err := bibtex.Parse(reader)
	if err != nil {
		log.Fatal(err)
	}
	/*
		*config != "" {
			 	var conf Config
			 	if _, err := toml.DecodeFile(*config, &conf); err != nil {
			 		log.Fatalf("Cannot read config: %s", err)
			 	}
			 	filter(parsed, &conf)
			 }
			fmt.Fprintf(writer, parsed.PrettyString())
	*/

	for _, entry := range parsed.Entries {
		if err = writeToDB(db, entry); err != nil {
			log.Fatalf("[%s] DB writing failed: %s", entry.CiteName, err)
		}
		//fmt.Println(entry.CiteName)
	}

	rows, err := db.Query("SELECT * FROM entries ORDER BY cite_name ASC")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	bib := bibtex.NewBibTex()
	for rows.Next() {
		var id int
		var citeName string
		var citeType string
		var author string
		var title string
		var address string
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
		var month string
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

		err = rows.Scan(&id, &citeName, &citeType, &author, &title, &address, &booktitle, &doi, &edition, &keyword, &location, &isbn, &issn, &institution, &journal, &metanote, &month, &note, &number, &numpages, &pages, &publisher, &series, &techreportType, &url, &version, &volume, &year)
		if err != nil {
			log.Fatal(err)
		}
		entry := bibtex.NewBibEntry(citeType, citeName)
		fieldMap := map[string]string{
			"author":      author,
			"title":       title,
			"address":     address,
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
			"month":       month,
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
		bib.AddEntry(entry)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(os.Stdout, bib.PrettyString())
}
