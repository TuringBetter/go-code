package main

import "fmt"

//go:generate stringer -type=Pill -linecomment
type Pill int

const (
	Placebo Pill = iota
	Aspirin
	Ibuprofen
	Paracetamol
)

func main() {
	fmt.Println("Pill types:", Placebo, Aspirin, Ibuprofen, Paracetamol)
}
