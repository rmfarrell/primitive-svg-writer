package primitiveWriter

import (
	"fmt"
	"reflect"
	"testing"

	w "github.com/rmfarrell/svg-server/primitive_writer"
)

func TestWrite(t *testing.T) {

}

func TestNewOptions(t *testing.T) {

	// valid
	wo, err := w.NewOptions(&w.Options{
		Input:      "./path/to/file",
		ShapeCount: 1,
	})
	if err != nil {
		t.Error(err)
	}

	// defaults set
	expected := &w.Options{
		Input:       "./path/to/file",
		ShapeCount:  1,
		Mode:        "triangle",
		Background:  "",
		Alpha:       128,
		Repetitions: 0,
	}

	if *wo != *expected {
		t.Errorf(`Did not correctly set defaults.
			Expected: 
			%v
			Received:
			%v`, expected, wo)
	}
}

func TestOptionsValidate(t *testing.T) {
	fixtures := []struct {
		input    *w.Options
		expected error
	}{
		{
			&w.Options{
				Input:      "./path/to/file",
				ShapeCount: 1,
			},
			nil,
		},
		{
			&w.Options{
				Input:      "",
				ShapeCount: 1,
			},
			fmt.Errorf("input param required"),
		},
		{
			&w.Options{
				Input:      "./path/to/file",
				ShapeCount: 0,
			},
			fmt.Errorf("shape_count param required"),
		},
	}
	for _, fx := range fixtures {
		_, err := w.NewOptions(fx.input)
		fmt.Println(reflect.TypeOf(fx.expected))
		if err == nil && fx.expected == nil {
			continue
		}
		if err.Error() != fx.expected.Error() {
			fmt.Println(reflect.TypeOf(err))
			t.Errorf(`Expected 
				%v 
				but received 
				%v`, fx.expected, err)
		}
	}
}
