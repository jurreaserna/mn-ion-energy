package main

import (
	"math"
	"math/rand"
	"time"
)

type ion struct {
	i_type      string
	spin        []float64
	s_magnitude float64
}

type lattice struct {
	ions [][][]ion
}

type interactionKey struct {
	ionA string
	ionB string
}

var interactions = map[interactionKey]float64{
	interactionKey{"mn4", "mn4"}:       0.0,
	interactionKey{"mn4", "mn3ver"}:    1.35,
	interactionKey{"mn4", "mnhor"}:     7.77,
	interactionKey{"mn3hor", "mn4"}:    7.77,
	interactionKey{"mn3ver", "mn4"}:    1.35,
	interactionKey{"mn3hor", "mn3ver"}: 4.65,
	interactionKey{"mn3ver", "mn3hor"}: 4.65,
}

func generateRand() float64 {
	seed := time.Now().UnixNano()
	source := rand.NewSource(seed)
	r := rand.New(source)
	return r.Float64()
}

func (i *ion) setupSpin() {
	magnitude := (*i).s_magnitude
	rNum := generateRand()
	sz := (2*rNum - 1) * magnitude
	sxy := math.Sqrt(math.Pow(magnitude, 2) - math.Pow(sz, 2))

	rNum = generateRand()
	phi := 2 * math.Pi * rNum
	sx := sxy * math.Cos(phi)
	sy := sxy * math.Sin(phi)
	(*i).spin = []float64{sx, sy, sz}
}

func createIon(n int) ion {
	if n%3 == 0 {
		return ion{i_type: "mn4", s_magnitude: 2.0}
	} else if n%3 == 1 {
		return ion{i_type: "mn3ver", s_magnitude: 1.5}
	} else {
		return ion{i_type: "mn3hor", s_magnitude: 1.5}
	}
}

func createLattice(length int, width int, height int) [][][]ion {
	var ions [][][]ion

	for i := 0; i < length; i++ {
		var ionMatrix [][]ion
		for j := 0; j < width; j++ {
			var ionRow []ion
			for k := 0; k < height; k++ {
				sIon := createIon(i + j + k)
				sIon.setupSpin()
				ionRow = append(ionRow, sIon)
			}
			ionMatrix = append(ionMatrix, ionRow)
		}
		ions = append(ions, ionMatrix)
	}
	return ions
}

func neighborEnergy(interactions map[interactionKey]float64, ionA ion, ionB ion) float64 {
	interaction := interactions[interactionKey{ionA: ionA.i_type, ionB: ionB.i_type}]
	return (ionA.spin[0]*ionB.spin[0] + ionA.spin[1]*ionB.spin[1] + ionA.spin[2]*ionB.spin[2]) * interaction
}

func energyDiff(latticeIon ion, ionTemp ion, l lattice, i int, j int, k int) float64 {
	e1, e2 := 0.0, 0.0

	neighbours := []ion{
		l.ions[i+1][j][k],
		l.ions[i-1][j][k],
		l.ions[i][j+1][k],
		l.ions[i][j-1][k],
		l.ions[i][j][k+1],
		l.ions[i][j][k-1],
	}

	for _, neighbour := range neighbours {
		e1 -= neighborEnergy(interactions, latticeIon, neighbour)
		e2 -= neighborEnergy(interactions, ionTemp, neighbour)
	}
	return e2 - e1
}

func (l *lattice) updateEnergy() {
	for i := 1; i < len((*l).ions)-1; i++ {
		for j := 1; j < len((*l).ions[0])-1; j++ {
			for k := 1; k < len((*l).ions[0][0])-1; k++ {

				latticeIon := (*l).ions[i][j][k]
				ionTemp := latticeIon
				ionTemp.setupSpin()

				eDiff := energyDiff(latticeIon, ionTemp, *l, i, j, k)

				if eDiff < 0 {
					(*l).ions[i][j][k] = ionTemp
				}
			}
		}
	}
}
