package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"

	"golang.org/x/image/draw"
	"golang.org/x/term"
)

func printUsage() {
	fmt.Println(`imgprint - a tool for printing a image to terminal using 24 bit colors
	Usage: imgpring [PATH TO IMAGE FOR PRINTING]`)
}

func getImageFile(userPath string) *os.File {
	var imgPath string

	if filepath.IsAbs(userPath) {
		imgPath = userPath
	} else if filepath.IsLocal(userPath) {
		wd, err := os.Getwd()
		if err != nil {
			fmt.Println("Could not obtain currecnt working directory")
		}

		imgPath = filepath.Join(wd, userPath)
	}

	if _, err := os.Stat(imgPath); os.IsNotExist(err) {
		fmt.Printf("file %s does not exist", imgPath)
		return nil
	}

	file, err := os.Open(imgPath)
	if err != nil {
		fmt.Println("Unknown error when opening file")
		return nil
	}
	return file
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	imagePath := os.Args[1]
	imageFile := getImageFile(imagePath)

	if imageFile == nil {
		fmt.Println("Could not open provided file")
		fmt.Println()
		printUsage()
		os.Exit(1)
	}

	img, _, err := image.Decode(imageFile)
	if err != nil {
		panic(err)
	}

	bounds := img.Bounds()

	width := 70
	height := 70

	if term.IsTerminal(0) {
		width, height, err = term.GetSize(0)
		if err != nil {
			panic(err)
		}
	}

	if bounds.Dx() >= bounds.Dy() {
		ratio := float64(bounds.Dx()) / float64(bounds.Dy())
		width = int(ratio * float64(height))
	} else {
		ratio := float64(bounds.Dy()) / float64(bounds.Dx())
		height = int(ratio * float64(width))
	}

	scaled := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.ApproxBiLinear.Scale(scaled, scaled.Rect, img, img.Bounds(), draw.Over, nil)

	printImage(scaled)
}
