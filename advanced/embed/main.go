package main

import (
	"embed"
	"fmt"
)

//go:embed static/*.txt
var static embed.FS

//go:embed text.txt
var text string

func main() {
	content, _ := static.ReadFile("static/text1.txt")
	fmt.Print(string(content))
	content, _ = static.ReadFile("static/text2.txt")
	fmt.Print(string(content))
	fmt.Println(text)
}
