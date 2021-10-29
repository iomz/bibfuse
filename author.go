package bibfuse

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Author holds an author information
type Author struct {
	FirstName string `default:"" sqlite3:"fist_name"`
	LastName  string `default:"" sqlite3:"last_name"`
}

// FirstNameCleaner returns a cleaned-up first name
func FirstNameCleaner(firstName string) (string, error) {
	match, _ := regexp.MatchString(`\b[A-ZÀ-Ú]{1}([^\.A-Za-zÀ-ÖØ-öø-ÿ]|\z)`, firstName)
	fmt.Println(match)
	if match {
		return "", errors.New("no dot in abbreviation")
	}
	return firstName, nil
}

// LastNameCleaner returns a cleaned-up first name
func LastNameCleaner(lastName string) (string, error) {
	match, _ := regexp.MatchString(`\A[A-ZÀ-Ú]{1}\.\z`, lastName)
	fmt.Println(match)
	if match {
		return "", errors.New("last name should not be abbreviated")
	}
	return lastName, nil
}

// NewAuthor returns a new Author
func NewAuthor(firstName, lastName string) (*Author, error) {
	a := new(Author)
	fn, err := FirstNameCleaner(firstName)
	if err != nil {
		return nil, err
	}
	a.FirstName = fn
	ln, err := LastNameCleaner(lastName)
	if err != nil {
		return nil, err
	}
	a.LastName = ln
	return a, nil
}

// Authors are a slice of multiple Author
type Authors []*Author

// String returns a string for the author field
func (as Authors) String() string {
	var sb strings.Builder
	for _, a := range as {
		if sb.Len() != 0 { // if it's not the first author
			sb.WriteString(" and ")
		}
		sb.WriteString(fmt.Sprintf("%s, ", a.LastName))
		sb.WriteString(a.FirstName)
	}
	return sb.String()
}
