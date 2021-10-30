package bibfuse

import (
	"fmt"
	"regexp"
	"strings"
)

// Author holds an author information
type Author struct {
	FirstName string `default:"" sqlite3:"fist_name"` // this can be empty
	LastName  string `default:"" sqlite3:"last_name"`
}

// BackslashCleaner minimizes sequences of backslashes
func BackslashCleaner(line string) string {
	multiBackSlash := regexp.MustCompile(`\\+`)
	return multiBackSlash.ReplaceAllString(line, `\`)
}

// FirstNameCleaner returns a cleaned-up first name
func FirstNameCleaner(firstName string) (string, error) {
	match, _ := regexp.MatchString(`\b[A-ZÀ-Ú]{1}([^\.A-Za-zÀ-ÖØ-öø-ÿ{}()]|\z)`, firstName)
	if match {
		return "", fmt.Errorf("no dot in abbreviation: %v", firstName)
	}
	return BackslashCleaner(firstName), nil
}

// LastNameCleaner returns a cleaned-up first name
func LastNameCleaner(lastName string) (string, error) {
	match, _ := regexp.MatchString(`\A[A-ZÀ-Ú]{1}\.\z`, lastName)
	if match {
		return "", fmt.Errorf("last name should not be abbreviated: %v", lastName)
	}
	return BackslashCleaner(lastName), nil
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

// NewAuthors return a new Authors
func NewAuthors(authorFieldValue string) (Authors, error) {
	var authors Authors

	rawAuthorsStringSlice := strings.Split(authorFieldValue, " and ")
	for _, rawAuthorString := range rawAuthorsStringSlice {
		authorNames := strings.Split(rawAuthorString, ",")
		a := new(Author)
		var err error
		switch len(authorNames) {
		case 1:
			// put the name to LastName
			a, err = NewAuthor("", strings.TrimSpace(authorNames[0]))
		case 2:
			a, err = NewAuthor(
				strings.TrimSpace(authorNames[1]),
				strings.TrimSpace(authorNames[0]),
			)
		default:
			err = fmt.Errorf("too many comma")
		}
		if err != nil {
			return authors, fmt.Errorf("%w %v", err, rawAuthorString)
		}

		authors = append(authors, a)
	}
	return authors, nil
}

// String returns a string for the author field
func (as Authors) String() string {
	var sb strings.Builder
	for _, a := range as {
		if sb.Len() != 0 { // if it's not the first author
			sb.WriteString(" and ")
		}
		if a.FirstName == "" {
			sb.WriteString(a.LastName)
		} else {
			sb.WriteString(fmt.Sprintf("%s, ", a.LastName))
			sb.WriteString(a.FirstName)
		}
	}
	return sb.String()
}
