package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math"
	"os"

	"github.com/mewkiz/pkg/errutil"
)

func ascii(filename string) (err error) {
	reader, err := os.Open(filename)
	if err != nil {
		return errutil.Err(err)
	}
	defer reader.Close()

	m, _, err := image.Decode(reader)
	if err != nil {
		return errutil.Err(err)
	}

	// Image width and height.
	width, height := m.Bounds().Dx(), m.Bounds().Dy()

	// Different change in x and y values depending on the aspect ratio.
	dx, dy := aspectRatio(width, height)

	for y := 0; y < height; y += dy {
		line := ""
		for x := 0; x < width; x += dx {
			line += level(m.At(x, y))
		}
		fmt.Println(line)
	}
	return nil
}

func aspectRatio(width, height int) (int, int) {
	// Approximation of the relation between font width / font height is 2.
	ratio := float64(width) / float64(height) / 2

	step := float64(stepFlag)
	var dx, dy float64 = step, step

	if ratio > 1 {
		// width > height
		dx = math.Ceil(step / ratio)
	} else {
		// height < width
		dy = math.Ceil(step / ratio)
	}
	return int(dx), int(dy)
}

// level takes a color and translates it into a contrast relative rune.
func level(c color.Color) string {
	// Different values ranging from black to white.
	var levels = []rune(" .,_-=:;+!|/$#@")

	// Different values ranging from white to black (string reversed).
	// var levels = []rune("@#$/|!+;:=-_,. ")

	// 3 colors, 256 different values divided by the amount of different
	// characters equals step size.
	var step = 3 * 256 / (len(levels) - 1)

	// Make colors 8 bit.
	r, g, b, _ := c.RGBA()
	r, g, b = r/256, g/256, b/256

	// From 0 to len(levels); What contrast is the current color?
	l := int(r+g+b) / step

	return string(levels[l])
}

var stepFlag int

func init() {
	flag.Usage = usage
	flag.IntVar(&stepFlag, "s", 10, "skip this many pixels between each character")
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [FILE],,,\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		flag.Usage()
	}
	for _, path := range flag.Args() {
		err := ascii(path)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
