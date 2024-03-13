package main

import (
	"cmp"
	"fmt"
	"image"
	"image/color"
	"slices"
	"strings"
	"sync"
)

type printData struct {
	id   int
	data []byte
}

const (
	esc             = "\033"
	setCustomColor  = "[48;2;"
	setDefaultColor = "[0m"
)

func printImage(img image.Image) {
	const maxRoutines = 100
	routines := min(img.Bounds().Dy()/2, maxRoutines)

	printDatas := make([]printData, routines)

	wg := sync.WaitGroup{}

	step := img.Bounds().Dy() / routines

	remainder := img.Bounds().Dy()
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
		fmt.Print(string(data.data), esc+setDefaultColor)
	}
}

func formatImagePart(img image.Image, startY, endY int) []byte {
	line := strings.Builder{}
	bounds := img.Bounds()

	endY = min(startY, bounds.Max.Y)

	for y := startY; y <= endY; y += 2 {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			col := img.At(x, y)
			ch := ' '
			line.WriteString(formatRGB(col, ch))
		}
		line.WriteString(esc + setDefaultColor + "\n")
	}

	return []byte(line.String())
}

// The returned text does not terminate the color change
func formatRGB(col color.Color, char rune) string {
	r, g, b, a := col.RGBA()

	to256 := func(col, a uint32) uint32 {
		return uint32(float64(col) / float64(a) * 0xff)
	}

	formatted := fmt.Sprintf(esc+setCustomColor+"%d;%d;%dm%c", to256(r, a), to256(g, a), to256(b, a), char)
	return formatted
}
