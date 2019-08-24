package shapee

import (
	"fmt"
	"github.com/r9y9/gossp/stft"
	"github.com/r9y9/gossp/window"
	"math/cmplx"
)

// ComputeSTFT computes the STFT of the input waveform and returns
// the result in polar rather than complex form.
func ComputeSTFT(waveform []float64, winShift int, frameLen int) ([][]float64, [][]float64) {
	s := &stft.STFT{
		FrameShift: winShift,
		FrameLen:   frameLen,
		Window:     window.CreateHanning(frameLen),
	}

	fmt.Printf("Performing STFT...\n")

	stftResult := s.STFT(waveform)

	// Complex to polar translation
	stftMag := make([][]float64, len(stftResult))
	stftAng := make([][]float64, len(stftResult))
	for i := 0; i < len(stftResult); i++ {
		stftMag[i] = make([]float64, len(stftResult[i])/2)
		stftAng[i] = make([]float64, len(stftResult[i])/2)
		for j := 0; j < len(stftMag[i]); j++ {
			stftMag[i][j], stftAng[i][j] = cmplx.Polar(stftResult[i][j])
		}
	}

	return stftMag, stftAng
}
