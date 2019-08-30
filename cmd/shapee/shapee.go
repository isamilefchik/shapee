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
	winShift := *(flag.Int("stft_shift", 1000, "STFT window shift in number of samples."))
	flag.Parse()

	fmt.Println("Importing freq ref audio...")
	freqWav, freqFormat := shapee.ImportAudio(*inFreqPath)
	//freqWav, freqBits, freqSR, freqNumChannels := shapee.ImportAudio(*inFreqPath)
	//fmt.Println("Importing amp ref audio...")
	//ampW, ampBits, ampSR, ampNumChannels := shapee.ImportAudio(*inAmpPath)
	//ampW, ampFormat := shapee.ImportAudio(*inAmpPath)

	var freqSTFT [][]complex128
	var freqMags [][]float64
	var freqPhases [][]float64

	//var ampSTFT [][]complex128
	//var ampMag [][]float64
	//var ampPhase [][]float64

	if freqFormat.NumChannels == 1 {
		freqSTFT = shapee.ComputeSTFT(freqWav, winShift, winLen)
		freqMags, freqPhases = shapee.ComplexToPolar(freqSTFT)
	} else {
		freqSTFT = shapee.ComputeSTFT(freqWav, winShift, winLen)
		freqMags, freqPhases = shapee.ComplexToPolar(freqSTFT)
	}

	freqSTFT = shapee.PolarToComplex(freqMags, freqPhases)

	//if ampFormat.NumChannels == 1 {
	//ampMag, ampPhase := shapee.ComputeSTFT(ampW, winShift, winLen)
	//} else {
	//ampMag, ampPhase := shapee.ComputeSTFT(ampW, winShift, winLen)
	//}

	iSTFTResult := shapee.ComputeISTFT(freqSTFT, winShift)
	shapee.ExportAudio(iSTFTResult, freqFormat, *outPath)

	// So Go doesn't get mad at me:
	//fmt.Println(freqWav[0])
	//fmt.Println(outPath)
	//fmt.Println(winLen)
	//fmt.Println(winShift)
	//fmt.Println(freqBits)
	//fmt.Println(freqSR)
	//fmt.Println(freqNumChannels)
	//fmt.Println(ampW[0])
	//fmt.Println(ampBits)
	//fmt.Println(ampSR)
	//fmt.Println(ampNumChannels)

}
