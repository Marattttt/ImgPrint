package main

import (
	"cmp"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"slices"
	"strings"
	"sync"
)

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

	printImage(&img)
}

type printData struct {
	id   int
	data []byte
}

func printImage(img *image.Image) {
	const maxRoutines = 100
	routines := min((*img).Bounds().Dy()/2, maxRoutines)

	printDatas := make([]printData, routines)

	wg := sync.WaitGroup{}

	step := (*img).Bounds().Dy() / routines

	remainder := (*img).Bounds().Dy()
	for i := 0; i < routines && remainder > 0; i++ {
		remainder -= step

		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			partData := formatImagePart(img, i*step, (i+1)*step)
			printDatas[i] = printData{
				id:   i,
				data: partData,
			}

		}(i)
	}

	wg.Wait()

	slices.SortFunc(printDatas, func(a printData, b printData) int {
		return cmp.Compare(a.id, b.id)
	})

	for _, data := range printDatas {
		fmt.Print(string(data.data), esc+endCustomColor)
	}
}

func formatImagePart(img *image.Image, startY, endY int) []byte {
	line := strings.Builder{}
	bounds := (*img).Bounds()

	endY = min(startY, bounds.Max.Y)

	for y := startY; y <= endY; y += 2 {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			col := (*img).At(x, y)
			ch := ' '
			line.WriteString(formatRGB(col, ch))
		}
		line.WriteString(esc + endCustomColor + "\n")
	}

	return []byte(line.String())
}

const (
	esc            = "\033"
	rgbBegin       = "[48;2;"
	endCustomColor = "[0m"
)

// The returned text does not terminate the color change
func formatRGB(col color.Color, char rune) string {
	r, g, b, a := col.RGBA()

	to256 := func(col, a uint32) uint32 {
		return uint32(float64(col) / float64(a) * 0xff)
	}

	formatted := fmt.Sprintf(esc+rgbBegin+"%d;%d;%dm%c", to256(r, a), to256(g, a), to256(b, a), char)
	return formatted
}
