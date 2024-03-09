package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strings"
)

func formatRGB(col color.Color, text string) string {
	const (
		esc      = "\033"
		rgbBegin = "[48;2;"
		rgbEnd   = "[0m"
	)

	r, g, b, a := col.RGBA()
	fmt.Printf("%d:%d:%d\n", (r*a)/0xffffff, (g*a)/0xffffff, (b*a)/0xffffff)

	deMultiply := func(col, a uint32) uint32 {
		return uint32(float64(col) / float64(a) * 0xff)
	}

	formatted := fmt.Sprintf(esc+rgbBegin+"%d;%d;%dm%s"+esc+rgbEnd, deMultiply(r, a), deMultiply(g, a), deMultiply(b, a), text)
	return formatted
}

// get s a char based on saturation
func getChar(col color.Color) rune {
	const asciiChars = `$@B%8&WM#*oahkbdpqwmZO0QLCJUYXzcvunxrjft/|()1{}[]?^-_+~i!lI;:,.`

	const maxPixelValue = 0xffff
	r, g, b, _ := col.RGBA()

	// values are alpha-premultipliexd, so there is no need to account for the alpha
	// brightness := float64(max(r, g, b)) / maxPixelValue

	saturation := float64(max(r, g, b)-min(r, g, b)) / maxPixelValue
	charIndex := int(float64(len(asciiChars)-1) * saturation)

	fmt.Printf("saturation: %f; index: %d\n", saturation, charIndex)

	return rune(asciiChars[charIndex])
}

func printUsage() {
	fmt.Print(`imgprint - a tool for printing a image to terminal using ascii symbols and 24 bit colors
	Usage: imgpring [PATH TO IMAGE FOR PRINTING]
`)
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	imagePath := os.Args[1]
	println(imagePath)

	for info, err := os.Stat(imagePath); err != nil || info.IsDir(); {
		fmt.Println(imagePath + " is not an image")
		fmt.Scanln(&imagePath)
	}

	imageFile, err := os.Open(imagePath)
	if err != nil {
		panic(err)
	}

	img, _, err := image.Decode(imageFile)
	if err != nil {
		panic(err)
	}

	fmt.Println(formatRGB(color.RGBA{255, 0, 255, 255}, "Hello world!"))

	printImage(&img)
}

func printImage(img *image.Image) {
	bounds := (*img).Bounds()
	fmt.Printf("%d by %d\n", bounds.Min.X, bounds.Min.Y)
	line := strings.Builder{}

	minAlpha := uint32(0xffff)
	for y := bounds.Min.Y; y < bounds.Max.Y; y += 2 {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			col := (*img).At(x, y)
			_, _, _, a := col.RGBA()
			minAlpha = min(minAlpha, a)
			// ch := getChar(col)
			ch := ' '
			line.WriteString(formatRGB(col, string(ch)))
		}
		line.WriteRune('\n')
	}

	fmt.Printf("Min alpha encountered: %d\n", minAlpha)

	fmt.Print(line.String())
}
