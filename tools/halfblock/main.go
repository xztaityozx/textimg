package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"os"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/math/fixed"
)

var (
	black    = color.RGBA{0, 0, 0, 255}
	red      = color.RGBA{255, 0, 0, 255}
	green    = color.RGBA{0, 255, 0, 255}
	yellow   = color.RGBA{255, 255, 0, 255}
	blue     = color.RGBA{0, 0, 255, 255}
	magenta  = color.RGBA{255, 0, 255, 255}
	cyan     = color.RGBA{0, 255, 255, 255}
	white    = color.RGBA{211, 211, 211, 255}
	darkGray = color.RGBA{169, 169, 169, 255}
)

func main() {
	const width = 200
	const height = 200
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	face := readFace("/usr/share/fonts/truetype/ricty-diminished/RictyDiminished-Regular.ttf", 50.0)

	drawBackgroundAll(img, black)
	args := os.Args
	for i, c := range []rune(args[1]) {
		drawLabel(img, i*50, 1*50, c, white, face)
	}
	drawLabel(img, 100, 50, []rune("▀")[0], white, face)
	drawLabel(img, 100, 50, []rune("▄")[0], red, face)
	drawLabel(img, 123, 49, []rune("▀")[0], blue, face)
	drawLabel(img, 146, 50, []rune("▄")[0], yellow, face)
	drawLabel(img, 169, 49, []rune("▀")[0], blue, face)
	// drawLabel(img, 146, 50, []rune("▀")[0], green, face)

	w, err := os.Create("out.png")
	if err != nil {
		panic(err)
	}
	defer w.Close()
	png.Encode(w, img)
}

// drawBackgroundAll はimgにbgを背景色として描画する。
func drawBackgroundAll(img *image.RGBA, bg color.RGBA) {
	var (
		bounds = img.Bounds().Max
		width  = bounds.X
		height = bounds.Y
	)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, bg)
		}
	}
}

// drawLabel はimgにラベルを描画する。
func drawLabel(img *image.RGBA, x, y int, r rune, col color.RGBA, face font.Face) {
	point := fixed.Point26_6{fixed.Int26_6(x * 64), fixed.Int26_6(y * 64)}
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: face,
		Dot:  point,
	}
	d.DrawString(string(r))
}

// readFace はfontPathのフォントファイルからfaceを返す。
func readFace(fontPath string, fontSize float64) font.Face {
	var fontData []byte

	// ファイルが存在しなければビルトインのフォントをデフォルトとして使う
	_, err := os.Stat(fontPath)
	if err == nil {
		fontData, err = ioutil.ReadFile(fontPath)
		if err != nil {
			panic(err)
		}
	} else {
		msg := fmt.Sprintf("[WARN] %s is not found. please set font path with `-f` option", fontPath)
		fmt.Fprintln(os.Stderr, msg)
		fontData = gomono.TTF
	}

	ft, err := truetype.Parse(fontData)
	if err != nil {
		panic(err)
	}
	opt := truetype.Options{
		Size:              fontSize,
		DPI:               0,
		Hinting:           0,
		GlyphCacheEntries: 0,
		SubPixelsX:        0,
		SubPixelsY:        0,
	}
	face := truetype.NewFace(ft, &opt)
	return face
}
