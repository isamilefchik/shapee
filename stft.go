package shapee

import (
	"fmt"
	//"github.com/r9y9/gossp/stft"
	//"github.com/r9y9/gossp/window"
	"github.com/mjibson/go-dsp/fft"
	"github.com/mjibson/go-dsp/window"
	"math"
	"math/cmplx"
)

// ComputeSTFT computes the STFT of the input waveform and returns
// it in complex form. It uses a Hamming window.
func ComputeSTFT(wave []float64, winShift int, winLen int) [][]complex128 {
	fmt.Printf("Performing STFT...\n")

	num_windows := int(math.Ceil(float64(len(wave)) / float64(winShift)))

	// Check and add zero-padding if it is needed at end of audio to make room
	// for final STFT window.
	waveEndLen := len(wave[(num_windows-1)*winShift : len(wave)])
	if waveEndLen < winLen {
		zeros := make([]float64, winLen-waveEndLen)
		wave = append(wave, zeros...)
	}

	stft := make([][]complex128, num_windows)

	for i := 0; i < num_windows; i++ {
		for j := i * winShift; j < (i+1)*winShift; j++ {
			waveWindow := make([]float64, winLen)
			copy(waveWindow, wave[j:j+winLen])
			window.Apply(waveWindow, window.Hamming)
			stft[i] = fft.FFTReal(waveWindow)
		}
	}

	return stft
}

// ComputeISTFT recovers the time-series from an STFT in complex form.
func ComputeISTFT(stft [][]complex128, winShift int) []float64 {
	wave := make([]float64, winShift*len(stft))
	realiFFTs := make([][]float64, len(stft))

	for i := range stft {
		iFFT := fft.IFFT(stft[i])
		realiFFTs[i] = make([]float64, len(iFFT))
		for j := range iFFT {
			realiFFTs[i][j] = cmplx.Abs(iFFT[j])
		}
		window.Apply(realiFFTs[i], window.Hamming)
	}

	// Overlap add
	winIndex := 0
	for i := range realiFFTs {
		for j := range realiFFTs[i] {
			wave[winIndex+j] += realiFFTs[i][j]
		}
		winIndex += winShift
	}

	return wave
}

// ComplexToPolar calculates the polar form of a complex STFT.
func ComplexToPolar(stft [][]complex128) ([][]float64, [][]float64) {
	stftMag := make([][]float64, len(stft))
	stftAng := make([][]float64, len(stft))
	for i := range stft {
		stftMag[i] = make([]float64, len(stft[i])/2)
		stftAng[i] = make([]float64, len(stft[i])/2)
		for j := range stftMag[i] {
			stftMag[i][j], stftAng[i][j] = cmplx.Polar(stft[i][j])
		}
	}

	return stftMag, stftAng
}

// PolarToComplex calculates the complex form of a polar STFT.
func PolarToComplex(stftMag [][]float64, stftAng [][]float64) [][]complex128 {
	stft := make([][]complex128, len(stftMag))
	for i := range stft {
		stft[i] = make([]complex128, len(stftMag[i]))
		for j := range stftMag[i] {
			stft[i][j] = cmplx.Rect(stftMag[i][j], stftAng[i][j])
		}
	}

	return stft
}
