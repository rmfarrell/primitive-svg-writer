package primitiveWriter

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/fogleman/primitive/primitive"
	"github.com/nfnt/resize"
)

var (
	inputMax int
	workers  int
	logLevel int
)

var modes = []string{
	"combo",
	"triangle",
	"rect",
	"ellipse",
	"circle",
	"rotatedrect",
	"beziers",
	"rotatedellipse",
	"polygon",
}

type shapeConfig struct {
	Count  int
	Mode   int
	Alpha  int
	Repeat int
}

type shapeConfigArray []shapeConfig

type Options struct {
	Input      string `json:input`
	ShapeCount int    `json:shapecount`
	Mode       string `json:mode`
	Background string `json:background`
	Alpha      int    `json:alpha`
	Repeat     int    `json:repeat`
}

type PrimitiveSvg struct {
	options *Options
	modeInt int
}

func init() {

	// assign defaults
	inputMax = 100
	logLevel = 1

	// set workers
	workers = runtime.NumCPU()
}

func (i *shapeConfigArray) Set(value, mode, alpha, repeat int) error {
	*i = append(*i, shapeConfig{value, mode, alpha, repeat})
	return nil
}

func NewPrimtitiveSvg(in *Options) (*PrimitiveSvg, error) {

	// set defaults
	out := &PrimitiveSvg{}
	options := &Options{"", 0, "triangle", "", 128, 0}

	options.Input = in.Input
	options.ShapeCount = in.ShapeCount

	if in.Mode != "" {
		options.Mode = in.Mode
	}
	if in.Background != "" {
		options.Background = in.Background
	}
	if in.Alpha != 0 {
		options.Alpha = in.Alpha
	}
	if in.Repeat > 0 {
		options.Repeat = in.Repeat
	}

	// assign modeInt (used by the algorithm)
	out.modeInt = -1
	for i := 0; i < len(modes); i++ {
		if modes[i] == options.Mode {
			out.modeInt = i
			break
		}
	}
	out.options = options
	err := out.validate()
	return out, err
}

func (psvg *PrimitiveSvg) validate() error {
	if psvg.options.Input == "" {
		return fmt.Errorf("input param required")
	}
	if psvg.options.ShapeCount == 0 {
		return fmt.Errorf("shape_count param required")
	}
	if psvg.modeInt < 0 {
		return fmt.Errorf(`%s is not a supported mode. Must be one of: %v`, psvg.options.Mode, modes)
	}
	return nil
}

func (psvg *PrimitiveSvg) Write() (string, error) {
	o := psvg.options

	// set configs
	configs := shapeConfigArray{}
	configs.Set(o.ShapeCount, psvg.modeInt, o.Alpha, o.Repeat)

	// set log level
	primitive.LogLevel = logLevel

	// seed random number generator
	rand.Seed(time.Now().UTC().UnixNano())

	// read input image
	primitive.Log(1, "reading %s\n", o.Input)
	input, err := primitive.LoadImage(o.Input)
	if err != nil {
		return "", err
	}

	// scale down image
	size := uint(inputMax)
	input = resize.Thumbnail(size, size, input, resize.Bilinear)

	// determine background color
	var bg primitive.Color
	if o.Background == "" {
		bg = primitive.MakeColor(primitive.AverageImageColor(input))
	} else {
		bg = primitive.MakeHexColor(o.Background)
	}
	// run algorithm
	out := ""
	model := primitive.NewModel(input, bg, psvg.modeInt, workers)
	primitive.Log(1, "%d: t=%.3f, score=%.6f\n", 0, 0.0, model.Score)
	start := time.Now()
	frame := 0
	for _, config := range configs {
		primitive.Log(1, "count=%d, mode=%d, alpha=%d, repeat=%d\n",
			config.Count, config.Mode, config.Alpha, config.Repeat)
		for i := 0; i < config.Count; i++ {
			frame++

			// find optimal shape and add it to the model
			t := time.Now()
			n := model.Step(primitive.ShapeType(config.Mode), config.Alpha, config.Repeat)
			nps := primitive.NumberString(float64(n) / time.Since(t).Seconds())
			elapsed := time.Since(start).Seconds()
			primitive.Log(1, "%d: t=%.3f, score=%.6f, n=%d, n/s=%s\n", frame, elapsed, model.Score, n, nps)
			if i >= config.Count-1 {
				out = model.SVG()
			}
		}
	}
	return out, nil
}
