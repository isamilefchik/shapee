package main

import (
	"flag"
	"fmt"
	"github.com/isamilefchik/shapee"
)

func main() {
	inFreqPath := flag.String("i_f", "./music/objekt.wav", "Filepath to frequency reference audio.")
	inAmpPath := flag.String("i_a", "./music/olixl.wav", "Filepath to amplitude reference audio.")
	//outPath := flag.String("o", "./music/result.wav", "Output audio filepath.")
	frameSize := *(flag.Int("stft_len", 1024, "STFT frame size in number of samples."))
	frameShift := *(flag.Int("stft_shift", 1000, "STFT frame shift in number of samples."))
	flag.Parse()

	fmt.Println("Importing freq ref audio...")
	freqW, _, _, freqNumChannels := shapee.ImportAudio(*inFreqPath)
	//freqW, freqBits, freqSR, freqNumChannels := shapee.ImportAudio(*inFreqPath)
	fmt.Println("Importing amp ref audio...")
	ampW, _, _, ampNumChannels := shapee.ImportAudio(*inAmpPath)
	//ampW, ampBits, ampSR, ampNumChannels := shapee.ImportAudio(*inAmpPath)

	var freqMag [][]float64
	//var freqPhase [][]float64
	//var ampMag [][]float64
	//var ampPhase [][]float64

	if freqNumChannels == 1 {
		freqMag, _ = shapee.ComputeSTFT(freqW, frameShift, frameSize)
		//freqMag, freqPhase := shapee.ComputeSTFT(freqW, frameShift, frameSize)
	} else {
		freqMag, _ = shapee.ComputeSTFT(freqW, frameShift, frameSize)
		//freqMag, freqPhase := shapee.ComputeSTFT(freqW, frameShift, frameSize)
	}

	if ampNumChannels == 1 {
		_, _ = shapee.ComputeSTFT(ampW, frameShift, frameSize)
		//ampMag, ampPhase := shapee.ComputeSTFT(ampW, frameShift, frameSize)
	} else {
		_, _ = shapee.ComputeSTFT(ampW, frameShift, frameSize)
		//ampMag, ampPhase := shapee.ComputeSTFT(ampW, frameShift, frameSize)
	}

	fmt.Println(freqMag)

	// So Go doesn't get mad at me:
	//fmt.Println(freqW[0])
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
