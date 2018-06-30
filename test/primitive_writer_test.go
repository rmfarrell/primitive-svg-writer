package primitiveWriter

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	w "github.com/rmfarrell/svg-server/primitive_writer"
)

func TestWrite(t *testing.T) {
	input := []byte(`{
				"input": "./this-guy.jpg",
				"shapecount": 2,
				"mode":"beziers"
			}`)

	optionsInput := w.Options{}
	err := json.Unmarshal(input, &optionsInput)
	if err != nil {
		t.Error(err)
	}
	primitiveSvg, err := w.NewPrimtitiveSvg(&optionsInput)
	if err != nil {
		t.Error(err)
	}
	svg, err := primitiveSvg.Write()
	if err != nil {
		t.Error(err)
	}
	if !(strings.HasPrefix(svg, "<svg")) {
		t.Errorf("Expected: %s to have prefix <svg", svg)
	}
}

func TestNewOptions(t *testing.T) {
	fixtures := []struct {
		input    *w.Options
		expected error
	}{
		// valid
		{
			&w.Options{
				Input:      "./path/to/file",
				ShapeCount: 1,
			},
			nil,
		},
		// no input
		{
			&w.Options{
				Input:      "",
				ShapeCount: 1,
			},
			fmt.Errorf("input param required"),
		},
		// no shape_count
		{
			&w.Options{
				Input:      "./path/to/file",
				ShapeCount: 0,
			},
			fmt.Errorf("shape_count param required"),
		},
		// unsupported mode
		{
			&w.Options{
				Input:      "./path/to/file",
				ShapeCount: 1,
				Mode:       "your mom",
			},
			fmt.Errorf("your mom is not a supported mode. Must be one of: [combo triangle rect ellipse circle rotatedrect beziers rotatedellipse polygon]"),
		},
	}
	for _, fx := range fixtures {
		_, err := w.NewPrimtitiveSvg(fx.input)
		if err == nil && fx.expected == nil {
			continue
		}
		if err.Error() != fx.expected.Error() {
			t.Errorf(`Expected
					%v
					but received
					%v`, fx.expected, err)
		}
	}
}
