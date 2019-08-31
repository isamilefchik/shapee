package main

import (
	"flag"
	"fmt"
	"github.com/isamilefchik/shapee"
)

func main() {
	inFreqPath := flag.String("i_f", "./music/objekt.wav", "Filepath to frequency reference audio.")
	//inAmpPath := flag.String("i_a", "./music/olixl.wav", "Filepath to amplitude reference audio.")
	outPath := flag.String("o", "./music/result.wav", "Output audio filepath.")
	winLen := *(flag.Int("stft_len", 1024, "STFT window size in number of samples."))
	winShift := *(flag.Int("stft_shift", 100, "STFT window shift in number of samples."))
	flag.Parse()

	fmt.Println("Importing freq ref audio...")

	freqWav, freqFormat := shapee.ImportWavAudio(*inFreqPath)

	iSTFTResult := make([][]float64, len(freqWav))
	for channel := range freqWav {
		freqSTFT := shapee.ComputeSTFT(freqWav[channel], winShift, winLen)
		freqMags, freqPhases := shapee.ComplexToPolar(freqSTFT)
		freqSTFT = shapee.PolarToComplex(freqMags, freqPhases)
		iSTFTResult[channel] = shapee.ComputeISTFT(freqSTFT, winShift)
	}
	shapee.ExportWavAudio(iSTFTResult, freqFormat, *outPath)
}
