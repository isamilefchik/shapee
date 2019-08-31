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
	freqWav, freqFormat := shapee.ImportAudio(*inFreqPath)
	//freqWav, freqBits, freqSR, freqNumChannels := shapee.ImportAudio(*inFreqPath)
	//fmt.Println("Importing amp ref audio...")
	//ampW, ampBits, ampSR, ampNumChannels := shapee.ImportAudio(*inAmpPath)
	//ampW, ampFormat := shapee.ImportAudio(*inAmpPath)

	//var freqSTFT [][]complex128
	//var freqMags [][]float64
	//var freqPhases [][]float64

	//var ampSTFT [][]complex128
	//var ampMag [][]float64
	//var ampPhase [][]float64

	//if freqFormat.NumChannels == 1 {
	//freqSTFT = shapee.ComputeSTFT(freqWav[0], winShift, winLen)
	//freqMags, freqPhases = shapee.ComplexToPolar(freqSTFT)
	//} else {
	//lFreqSTFT = shapee.ComputeSTFT(freqWav[0], winShift, winLen)
	//rFreqSTFT = shapee.ComputeSTFT(freqWav[1], winShift, winLen)
	//lFreqMags, lFreqPhases = shapee.ComplexToPolar(lFreqSTFT)
	//rFreqMags, rFreqPhases = shapee.ComplexToPolar(rFreqSTFT)
	//}

	iSTFTResult := make([][]float64, int(freqFormat.NumChannels))
	for channel := 0; channel < int(freqFormat.NumChannels); channel++ {
		freqSTFT := shapee.ComputeSTFT(freqWav[channel], winShift, winLen)
		freqMags, freqPhases := shapee.ComplexToPolar(freqSTFT)
		freqSTFT = shapee.PolarToComplex(freqMags, freqPhases)
		iSTFTResult[channel] = shapee.ComputeISTFT(freqSTFT, winShift)
	}
	shapee.ExportAudio(iSTFTResult, freqFormat, *outPath)

	//if ampFormat.NumChannels == 1 {
	//ampMag, ampPhase := shapee.ComputeSTFT(ampW, winShift, winLen)
	//} else {
	//ampMag, ampPhase := shapee.ComputeSTFT(ampW, winShift, winLen)
	//}

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
