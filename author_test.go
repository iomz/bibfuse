package bibfuse

import (
	"errors"
	"reflect"
	"testing"
)

var authortests = []struct {
	in  map[string]string
	err error
	out *Author
}{
	{
		map[string]string{"first_name": "Iori", "last_name": "Mizutani"},
		nil,
		&Author{"Iori", "Mizutani"},
	},
	{
		map[string]string{"first_name": "Salvador Domingo Felipe Jacinto", "last_name": "Dal{\\'i} i Dom{\\`e}nech, 1st Marquess of Dal{\\'i} of P{\\'u}bol"},
		nil,
		&Author{"Salvador Domingo Felipe Jacinto", "Dal{\\'i} i Dom{\\`e}nech, 1st Marquess of Dal{\\'i} of P{\\'u}bol"},
	},
	{
		map[string]string{"first_name": "William B.", "last_name": "Pitt"},
		nil,
		&Author{"William B.", "Pitt"},
	},
	{
		map[string]string{"first_name": "William B.", "last_name": "Pitt"},
		nil,
		&Author{"William B.", "Pitt"},
	},
	{
		map[string]string{"first_name": "William B", "last_name": "Pitt"},
		errors.New("no dot in abbreviation"),
		nil,
	},
	{
		map[string]string{"first_name": "Iori", "last_name": "M."},
		errors.New("last name should not be abbreviated"),
		nil,
	},
}

func TestNewAuthor(t *testing.T) {
	for _, tt := range authortests {
		a, err := NewAuthor(tt.in["first_name"], tt.in["last_name"])
		if err != nil && err.Error() != tt.err.Error() {
			t.Errorf("NewAuthor(%v, %v) err => %v, want %v", tt.in["first_name"], tt.in["last_name"], err, tt.err)
		}
		if !reflect.DeepEqual(a, tt.out) {
			t.Errorf("NewAuthor(%v, %v) => %v, want %v", tt.in["first_name"], tt.in["last_name"], a, tt.out)
		}
	}
}

var authorsstringtests = []struct {
	in  Authors
	out string
}{
	{
		Authors{
			&Author{"Iori", "Mizutani"},
			&Author{"Ganesh", "Ramanathan"},
			&Author{"Simon", "Mayer"},
		},
		"Mizutani, Iori and Ramanathan, Ganesh and Mayer, Simon",
	},
}

func TestAuthorsString(t *testing.T) {
	for _, tt := range authorsstringtests {
		result := tt.in.String()
		if result != tt.out {
			t.Errorf("Authors.String() => %v, want %v", result, tt.out)
		}
	}
}

var authorstests = []struct {
	in  string
	err error
	out string
}{
	{
		"Mizutani, Iori and Ramanathan, Ganesh and Mayer, Simon",
		nil,
		"Mizutani, Iori and Ramanathan, Ganesh and Mayer, Simon",
	},
	{
		"Internet Engineering Task Force",
		nil,
		"Internet Engineering Task Force",
	},
	{
		"Dal{\\'i} i Dom{\\`e}nech, 1st Marquess of Dal{\\'i} of P{\\'u}bol, Salvador Domingo Felipe Jacinto",
		errors.New("too many comma"),
		"",
	},
	{
		"Mizutani, Iori and Pitt, William B.",
		errors.New(""),
		"Mizutani, Iori and Pitt, William B.",
	},
	{
		"Pitt, William B",
		errors.New("no dot in abbreviation"),
		"",
	},
	{
		"Mizutani, Iori, Dr.sc.",
		errors.New("too many comma"),
		"",
	},
}

func TestAuthors(t *testing.T) {
	for _, tt := range authorstests {
		authors, err := NewAuthors(tt.in)
		if err != nil && err.Error() != tt.err.Error() {
			t.Errorf("NewAuthors(%v) err => %v, want %v", tt.in, err, tt.err)
		}
		if authors.String() != tt.out {
			t.Errorf("Authors.String() => %v, want %v", authors.String(), tt.out)
		}
	}
}
