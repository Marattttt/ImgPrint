package main

import "fmt"

const (
	esc = "\033"
)

func formatRGB(r, g, b int, text string) string {
	formatted := fmt.Sprintf(esc+"[38;2;%d;%d;%dm%s"+esc+"[0m", r, g, b, text)
	return formatted
}

func main() {
	fmt.Println(formatRGB(255, 255, 125, "Hello world!"))
}
