package shapee

import (
	"fmt"
	"github.com/mjibson/go-dsp/fft"
	"github.com/mjibson/go-dsp/window"
	"math"
	"math/cmplx"
	"sync"
)

// ComputeSTFT computes the STFT of the input waveform and returns
// it in complex form. It uses a Hamming window.
func ComputeSTFT(wave []float64, winShift int, winLen int) [][]complex128 {
	fmt.Printf("      Performing STFT... ")

	num_windows := int(math.Floor(float64(len(wave)) / float64(winShift)))

	// Check and add zero-padding if it is needed at end of audio to make room
	// for final STFT window.
	waveEndLen := len(wave[num_windows*winShift:])
	if waveEndLen < winLen {
		zeros := make([]float64, winLen-waveEndLen)
		wave = append(wave, zeros...)
	}

	stft := make([][]complex128, num_windows)

	var waitGroup sync.WaitGroup
	waitGroup.Add(num_windows)

	for i := 0; i < num_windows; i++ {
		iGo := i
		go func() {
			defer waitGroup.Done()
			waveWindow := make([]float64, winLen)
			copy(waveWindow, wave[iGo*winShift:iGo*winShift+winLen])
			window.Apply(waveWindow, window.Hamming)
			cmplxWindow := make([]complex128, winLen)
			for j := range cmplxWindow {
				cmplxWindow[j] = complex(waveWindow[j], 0.0)
			}
			stft[iGo] = fft.FFT(cmplxWindow)
		}()
	}

	waitGroup.Wait()

	fmt.Printf("                 Done.\n")
	return stft
}

// ComputeISTFT recovers the time-series from an STFT in complex form.
func ComputeISTFT(stft [][]complex128, winShift int) []float64 {
	fmt.Printf("   Recovering time-series (iSTFT)...")
	wave := make([]float64, winShift*len(stft)+(2*len(stft[0])))
	realiFFTs := make([][]float64, len(stft))

	// IFFTs:
	for i := range stft {
		iFFT := fft.IFFT(stft[i])
		realiFFTs[i] = make([]float64, len(iFFT))
		for j := range iFFT {
			realiFFTs[i][j] = real(iFFT[j])
		}
		window.Apply(realiFFTs[i], window.Hamming)
	}

	// Overlap add:
	winIndex := 0
	for i := range realiFFTs {
		for j := range realiFFTs[i] {
			wave[winIndex+j] += realiFFTs[i][j]
		}
		winIndex += winShift
	}

	// Normalization:
	max := 0.
	for i := range wave {
		if math.Abs(wave[i]) > max {
			max = math.Abs(wave[i])
		}
	}
	for i := range wave {
		if math.Abs(wave[i]) < max {
			wave[i] = 0.5 * (wave[i] / max)
		} else {
			wave[i] = 0.5
		}
	}

	fmt.Printf("      Done.\n\n")

	return wave
}

// ComplexToPolar calculates the polar form of a complex STFT.
func ComplexToPolar(stft [][]complex128) ([][]float64, [][]float64) {
	stftMag := make([][]float64, len(stft))
	stftAng := make([][]float64, len(stft))
	for i := range stft {
		stftMag[i] = make([]float64, len(stft[i]))
		stftAng[i] = make([]float64, len(stft[i]))
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
