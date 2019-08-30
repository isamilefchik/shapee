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
	frameSize := *(flag.Int("stft_len", 1024, "STFT frame size in number of samples."))
	frameShift := *(flag.Int("stft_shift", 1000, "STFT frame shift in number of samples."))
	flag.Parse()

	fmt.Println("Importing freq ref audio...")
	freqWav, freqFormat := shapee.ImportAudio(*inFreqPath)
	//freqWav, freqBits, freqSR, freqNumChannels := shapee.ImportAudio(*inFreqPath)
	//fmt.Println("Importing amp ref audio...")
	//ampW, ampBits, ampSR, ampNumChannels := shapee.ImportAudio(*inAmpPath)
	ampW, ampFormat := shapee.ImportAudio(*inAmpPath)

	var freqMag [][]float64
	var freqPhase [][]float64
	//var ampMag [][]float64
	//var ampPhase [][]float64

	if freqNumChannels == 1 {
		freqSTFT := shapee.ComputeSTFT(freqWav, frameShift, frameSize)
		freqMags, freqPhases := shapee.ComplexToPolar(freqSTFT)
	} else {
		freqSTFT := shapee.ComputeSTFT(freqWav, frameShift, frameSize)
		freqMags, freqPhases := shapee.ComplexToPolar(freqSTFT)
	}

	//if ampNumChannels == 1 {
	//ampMag, ampPhase := shapee.ComputeSTFT(ampW, frameShift, frameSize)
	//} else {
	//ampMag, ampPhase := shapee.ComputeSTFT(ampW, frameShift, frameSize)
	//}

	//fmt.Println(len(freqMag[0]))

	// So Go doesn't get mad at me:
	//fmt.Println(freqWav[0])
	//fmt.Println(outPath)
	//fmt.Println(frameSize)
	//fmt.Println(frameShift)
	//fmt.Println(freqBits)
	//fmt.Println(freqSR)
	//fmt.Println(freqNumChannels)
	//fmt.Println(ampW[0])
	//fmt.Println(ampBits)
	//fmt.Println(ampSR)
	//fmt.Println(ampNumChannels)

}
