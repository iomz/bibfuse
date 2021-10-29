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

var authorstests = []struct {
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
	for _, tt := range authorstests {
		result := tt.in.String()
		if result != tt.out {
			t.Errorf("Authors.String() => %v, want %v", result, tt.out)
		}
	}
}
