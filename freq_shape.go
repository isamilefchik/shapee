package shapee

import (
	"fmt"
)

type ShaperPack struct {
	AmpMag    [][]float64
	FreqMag   [][]float64
	FreqPhase [][]float64
	W         int
}

// FreqShape completes the audio transformation as proposed in
// Christopher Penrose's paper "Frequency Shaping of Audio Signals".
func FreqShape(pack *ShaperPack) ([][]float64, [][]float64) {
	fmt.Printf("\U0001f608 Frequency shaping... ")

	aMag := pack.AmpMag
	fMag := pack.FreqMag
	fPhase := pack.FreqPhase
	w := pack.W

	bigN := len(fMag[0])
	s := make([][]float64, len(fMag))

	// The length of the result is the length of the shorter of the
	// two input audio files.
	var rLen int
	if len(fMag) > len(aMag) {
		rLen = len(aMag)
	} else {
		rLen = len(fMag)
	}
	rMag := make([][]float64, rLen)

	for i := 0; i < rLen; i++ {
		s[i] = make([]float64, ((bigN/2-1)/w)+1)
		for j := range s[i] {
			aSum := 0.0
			fSum := 0.0
			for n := 0; n <= w; n++ {
				aSum += aMag[i][j*w+n]
				fSum += fMag[i][j*w+n]
			}
			s[i][j] = aSum / fSum
		}

		rMag[i] = make([]float64, len(fMag[i]))
		for k := 0; k <= bigN/2-1; k++ {
			kPrime := k / w
			rMag[i][k] = fMag[i][k] * s[i][kPrime]

			// Reflect on upper-half in preperation for iSTFT
			rMag[i][len(rMag[i])-1-k] = rMag[i][k]
		}
	}

	// For the sake of staying consistent to the paper:
	rPhase := fPhase

	fmt.Printf("                  Done.\n")

	return rMag, rPhase
}
