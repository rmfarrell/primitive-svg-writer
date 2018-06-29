package primitiveWriter

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	// Input      string
	Outputs    flagArray
	Background string
	Configs    shapeConfigArray
	Alpha      int
	InputSize  int
	OutputSize int
	Mode       int
	Workers    int
	Nth        int
	Repeat     int
	V, VV      bool
)

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
	Input       string `json:input`
	ShapeCount  int    `json:shape_count`
	Mode        string `json:mode`
	Background  string `json:background`
	Alpha       int    `json:alpha`
	Repetitions int    `json:repetitions`
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

func NewOptions(in *Options) (*Options, error) {

	// set defaults
	out := &Options{"", 0, "triangle", "", 128, 0}

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
	if in.Repetitions > 0 {
		out.Repetitions = in.Repetitions
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

func init() {
	// flag.StringVar(&Input, "i", "", "input image path")
	// flag.Var(&Outputs, "o", "output image path")
	// flag.Var(&Configs, "n", "number of primitives")
	flag.StringVar(&Background, "bg", "", "background color (hex)")
	flag.IntVar(&Alpha, "a", 128, "alpha value")
	flag.IntVar(&InputSize, "r", 256, "resize large input images to this size")
	// flag.IntVar(&OutputSize, "s", 1024, "output image size")
	flag.IntVar(&Mode, "m", 1, "0=combo 1=triangle 2=rect 3=ellipse 4=circle 5=rotatedrect 6=beziers 7=rotatedellipse 8=polygon")
	flag.IntVar(&Workers, "j", 0, "number of parallel workers (default uses all cores)")
	flag.IntVar(&Nth, "nth", 1, "save every Nth frame (put \"%d\" in path)")
	flag.IntVar(&Repeat, "rep", 0, "add N extra shapes per iteration with reduced search")
	flag.BoolVar(&V, "v", false, "verbose")
	flag.BoolVar(&VV, "vv", false, "very verbose")
}

/*
func Write(options *Options) (string, error) {
	// validate options

	// if len(Outputs) == 0 {
	// 	ok = errorMessage("ERROR: output argument required")
	// }

	// set configs
	configs := shapeConfigArray{}
	configs.Set(options.PrimitiveCount)

	if len(Configs) == 1 {
		Configs[0].Mode = Mode
		Configs[0].Alpha = Alpha
		Configs[0].Repeat = Repeat
	}

	// set log level
	primitive.LogLevel = 2

	// seed random number generator
	rand.Seed(time.Now().UTC().UnixNano())

	// determine worker count
	if Workers < 1 {
		Workers = runtime.NumCPU()
	}
	fmt.Println("Workers gathered")

	// read input image
	primitive.Log(1, "reading %s\n", Input)
	input, err := primitive.LoadImage(Input)
	check(err)
	fmt.Println("Read Inage")

	// scale down input image if needed
	size := uint(InputSize)
	if size > 0 {
		input = resize.Thumbnail(size, size, input, resize.Bilinear)
	}
	fmt.Println("Scale")

	// determine background color
	var bg primitive.Color
	if Background == "" {
		bg = primitive.MakeColor(primitive.AverageImageColor(input))
	} else {
		bg = primitive.MakeHexColor(Background)
	}
	fmt.Println("Got Background Color")

	// run algorithm
	model := primitive.NewModel(input, bg, 1, Workers)
	fmt.Println("Got Model")
	primitive.Log(1, "%d: t=%.3f, score=%.6f\n", 0, 0.0, model.Score)
	start := time.Now()
	frame := 0
	for j, config := range Configs {
		fmt.Println("Start Config Loop")
		primitive.Log(1, "count=%d, mode=%d, alpha=%d, repeat=%d\n",
			config.Count, config.Mode, config.Alpha, config.Repeat)

		for i := 0; i < config.Count; i++ {
			fmt.Println("Start config count new loop")
			fmt.Println(i)
			frame++

			// find optimal shape and add it to the model
			t := time.Now()
			fmt.Println(t)
			n := model.Step(primitive.ShapeType(config.Mode), config.Alpha, config.Repeat)
			fmt.Println("n assigned")
			nps := primitive.NumberString(float64(n) / time.Since(t).Seconds())
			fmt.Println("nps assigned")
			elapsed := time.Since(start).Seconds()
			fmt.Println(elapsed)
			primitive.Log(1, "%d: t=%.3f, score=%.6f, n=%d, n/s=%s\n", frame, elapsed, model.Score, n, nps)
			last := j == len(Configs)-1 && i == config.Count-1
			if last {
				fmt.Println("End Config Loop")
				return model.SVG(), nil
			}
		}
	}
	return "", nil
}
*/
