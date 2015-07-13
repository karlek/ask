// Program ask converts images to ascii art in text format.
package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math"
	"os"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/mewkiz/pkg/errutil"
)

// ascii opens an image file and prints an ascii art image.
func ascii(filename string) (err error) {
	reader, err := os.Open(filename)
	if err != nil {
		return errutil.Err(err)
	}
	defer reader.Close()

	i, _, err := image.Decode(reader)
	if err != nil {
		return errutil.Err(err)
	}

	// Image width and height.
	width, height := i.Bounds().Dx(), i.Bounds().Dy()

	// Different change in x and y values depending on the aspect ratio.
	dx, dy := aspectRatio(width, height)

	// Print each line.
	var line string
	for y := 0; y < height; y += dy {
		line = ""
		// Create a line. Convert the color level of the pixel into a ascii
		// character.
		for x := 0; x < width; x += dx {
			line += level(i.At(x, y).RGBA())
		}
		fmt.Println(line)
	}
	return nil
}

// aspectRatio returns the approximative number of steps on the x- and y-axis.
func aspectRatio(width, height int) (int, int) {
	// Approximation of the relation between font width / font height is 2.
	ratio := float64(width) / float64(height) / 2

	// Number of pixels to ignore.
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
func level(r, g, b, _ uint32) string {
	// Different values ranging from black to white.
	var levels = []rune(" .,_-=:;+!|/$#@")
	var ansiFmt = "\x1b[38;2;%d;%d;%dm%c\x1b[0m"
	// Different values ranging from white to black (string reversed).
	// var levels = []rune("@#$/|!+;:=-_,. ")

	// 3 colors, 256 different values divided by the amount of different
	// characters equals step size.
	var step = float64(len(levels) - 1)

	// Make (0,1).
	c := colorful.Color{float64(r>>8) / 256.0, float64(g>>8) / 256.0, float64(b>>8) / 256.0}
	_, _, v := c.Hsv()

	// From 0 to len(levels); What contrast is the current color?
	l := int(v * step)

	ret := string(levels[l])
	if color {
		ret = fmt.Sprintf(ansiFmt, r>>8, g>>8, b>>8, levels[l])
	}
	// Return the character corresponding to the approximative black and white value.
	return ret
}

// stepFlag is the amount of pixels skipped between each sample.
var (
	stepFlag int
	color    bool
)

func init() {
	flag.Usage = usage
	flag.IntVar(&stepFlag, "s", 10, "skip this many pixels between each character")
	flag.BoolVar(&color, "c", true, "colored ansi output")
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] [FILE],,,\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		flag.Usage()
	}
	// For each file passed as arguments, convert the images to ascii art.
	for _, path := range flag.Args() {
		err := ascii(path)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
