package primitiveWriter

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/fogleman/primitive/primitive"
	"github.com/nfnt/resize"
)

var (
	// Input      string
	Outputs    flagArray
	Background string
	Configs    shapeConfigArray
	Alpha      int
	inputMax   int
	OutputSize int
	Mode       int
	workers    int
	Nth        int
	Repeat     int
	V, VV      bool
	logLevel   int
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

type flagArray []string

func (i *flagArray) String() string {
	return strings.Join(*i, ", ")
}

func (i *flagArray) Set(value string) error {
	*i = append(*i, value)
	return nil
}

type shapeConfig struct {
	Count  int
	Mode   int
	Alpha  int
	Repeat int
}

type writeParams struct {
	Configs shapeConfigArray
}

type Options struct {
	Input      string `json:input`
	ShapeCount int    `json:shapecount`
	Mode       string `json:mode`
	Background string `json:background`
	Alpha      int    `json:alpha`
	Repeat     int    `json:repeat`
	modeInt    int
}

// type writeOptions struct {
// 	Background string
// 	Alpha      int
// 	InputSize  int
// 	OutputSize int
// 	Mode       int
// 	Workers    int
// 	Nth        int
// 	Repeat     int
// 	V, VV      bool
// }

func init() {
	// assign defaults
	inputMax = 100
	logLevel = 1
	// set workers
	workers = runtime.NumCPU()
}

func NewOptions(in *Options) (*Options, error) {

	// set defaults
	out := &Options{"", 0, "triangle", "", 128, 0, 0}

	out.Input = in.Input
	out.ShapeCount = in.ShapeCount
	// TODO: account for combo, if supported
	if in.Mode != "" {
		out.Mode = in.Mode
	}
	if in.Background != "" {
		out.Background = in.Background
	}
	if in.Alpha != 0 {
		out.Alpha = in.Alpha
	}
	if in.Repeat > 0 {
		out.Repeat = in.Repeat
	}

	// assign modeInt (used by the algorithm)
	out.modeInt = -1
	for i := 0; i < len(modes); i++ {
		if modes[i] == out.Mode {
			out.modeInt = i
			break
		}
	}

	err := out.validate()
	return out, err
}

func (wo *Options) validate() error {
	if wo.Input == "" {
		return fmt.Errorf("input param required")
	}
	if wo.ShapeCount == 0 {
		return fmt.Errorf("shape_count param required")
	}
	if wo.modeInt < 0 {
		return fmt.Errorf(`%s is not a supported mode. Must be one of: %v`, wo.Mode, modes)
	}

	// TODO: What is this? Is it needed?
	// for _, config := range Configs {
	// 	if config.Count < 1 {
	// 		ok = errorMessage("ERROR: number argument must be > 0")
	// 	}
	// }
	return nil
}

type shapeConfigArray []shapeConfig

func (i *shapeConfigArray) String() string {
	return ""
}

func (i *shapeConfigArray) Set(value int) error {
	*i = append(*i, shapeConfig{value, Mode, Alpha, Repeat})
	return nil
}

func errorMessage(message string) bool {
	fmt.Fprintln(os.Stderr, message)
	return false
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Write(o *Options) (string, error) {

	// set configs
	configs := shapeConfigArray{}
	configs.Set(o.ShapeCount)

	configs[0].Mode = o.modeInt
	configs[0].Alpha = o.Alpha
	configs[0].Repeat = o.Repeat

	// set log level
	primitive.LogLevel = logLevel

	// seed random number generator
	// TODO: remove?
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
	model := primitive.NewModel(input, bg, o.modeInt, workers)
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
