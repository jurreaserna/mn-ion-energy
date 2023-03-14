package main

import (
	"fmt"
)

func main() {
	length, width := 30, 30
	height := 5
	ionLattice := lattice{ions: createLattice(length, width, height)}
	i := ionLattice.ions[1][10][2]
	fmt.Printf("%+v\n", i)
	ionLattice.updateEnergy()
	i = ionLattice.ions[1][10][2]
	fmt.Printf("%+v\n", i)
}
