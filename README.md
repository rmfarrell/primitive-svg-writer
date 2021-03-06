# Primitive SVG Writer
Fork of [Fogleman's fabulous Primitives library](https://github.com/fogleman/primitive) which outputs an SVG string instead of outputting to another file.
This library scans a bitmap and uses Primitives to output an SVG string.
## Usage
```
import primitive primitiveSVGWriter

primitive, err := NewPrimtitiveSvg(&primitive.Options{
  Input: "./path/to/file.jpg", // required
  ShapeCount: 100, // required
  Mode: "beziers", // combo | triangle | rect | ellipse | circle | rotatedrect | beziers | rotatedellipse | polygon
  Background: "#000000",
  Alpha: 200,
  Repeat: 1,
})

if err != nil {
  fmt.Println(err)
}

svg, err := primitive.Write()
if err != nil {
  fmt.Println(err)
}

fmt.Println(svg) // sweet SVG goodness

```
See https://github.com/fogleman/primitive for more info.

## Options
- **Input**      {String} local path to bitmap file
- **ShapeCount** {Int} number of shapes to render 
- **Mode**       {String} (combo | triangle | rect | ellipse | circle | rotatedrect | beziers | rotatedellipse | polygon) rendering algorithm used for shape
- **Alpha**      {Int} transparency of rendered shapes
- **Repeat**     {Int} add N extra shapes per iteration with reduced search

## Y Tho?
The original library is great for local file creation, but this has no opinion on the output. You can write the SVG over HTTP, for example.



