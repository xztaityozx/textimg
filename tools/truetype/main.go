// https://codeday.me/jp/qa/20190224/308772.html

package main

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"os"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

var (
	dpi      = flag.Float64("dpi", 72, "screen resolution in Dots Per Inch")
	fontfile = flag.String("fontfile", "/usr/share/fonts/truetype/ricty-diminished/RictyDiminished-Regular.ttf", "filename of the ttf font")
	hinting  = flag.String("hinting", "none", "none | full")
	size     = flag.Float64("size", 125, "font size in points")
	spacing  = flag.Float64("spacing", 1.5, "line spacing (e.g. 2 means double spaced)")
	wonb     = flag.Bool("whiteonblack", false, "white text on a black background")
	// text     = string("JOJOあいあい▄▀")
	text = string("▄▀▄▀▄▀▄▀▄▀▄▀▄▀▄▀▄▀▄▀▄▀")
)

func main() {
	flag.Parse()
	fmt.Printf("Loading fontfile %q\n", *fontfile)
	b, err := ioutil.ReadFile(*fontfile)
	if err != nil {
		log.Println(err)
		return
	}
	f, err := truetype.Parse(b)
	if err != nil {
		log.Println(err)
		return
	}

	width := len([]rune(text)) * int(*size) / 2
	// width := runewidth.StringWidth(text) * int(*size)
	height := 1 * int(*size)

	fmt.Printf("width = %d\n", width)
	fmt.Printf("height = %d\n", height)
	fmt.Println("-------------")

	// Freetype context
	fg, bg := image.Black, image.White
	rgba := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(*dpi)
	c.SetFont(f)
	c.SetFontSize(*size)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)
	switch *hinting {
	default:
		c.SetHinting(font.HintingNone)
	case "full":
		c.SetHinting(font.HintingFull)
	}

	// Make some background

	// Draw the guidelines.
	ruler := color.RGBA{0xdd, 0xdd, 0xdd, 0xff}
	for rcount := 0; rcount < 4; rcount++ {
		for i := 0; i < 200; i++ {
			rgba.Set(250*rcount, i, ruler)
		}
	}

	// Truetype stuff
	opts := truetype.Options{}
	opts.Size = *size
	face := truetype.NewFace(f, &opts)

	// Calculate the widths and print to image
	for i, x := range []rune(text) {
		bounds, advance, ok := face.GlyphBounds(x)
		if !ok {
			return
		}
		fmt.Printf("Bounds: %v\n", bounds)
		fmt.Printf("Advance: %v\n", advance)
		fmt.Printf("Metrics height: %v\n", face.Metrics().Height)
		fmt.Printf("Metrics xheight: %v\n", face.Metrics().XHeight)
		awidth, ok := face.GlyphAdvance(x)
		if ok != true {
			log.Println(err)
			return
		}
		iwidthf := int(float64(awidth) / (*size / 2))
		// if runewidth.RuneWidth(x) == 1 {
		// 	iwidthf /= 2
		// }

		fmt.Printf("iwidthf = %+v\n", iwidthf)

		pt := freetype.Pt(i*iwidthf, int(*size))
		c.DrawString(string(x), pt)
		fmt.Printf("awidth = %+v\n", awidth)
	}

	// Save that RGBA image to disk.
	outFile, err := os.Create("out.png")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer outFile.Close()
	bf := bufio.NewWriter(outFile)
	err = png.Encode(bf, rgba)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	err = bf.Flush()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	fmt.Println("Wrote out.png OK.")

}
