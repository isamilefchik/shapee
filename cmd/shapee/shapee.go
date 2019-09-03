package main

import (
	"flag"
	"fmt"
	"github.com/gobuffalo/packr"
	"github.com/isamilefchik/shapee"
	"os"
)

func main() {

	inAmpPath := flag.String("ia",
		"./music/lanskyTimbre.wav", "Filepath to amplitude reference audio.")
	inFreqPath := flag.String("if",
		"./music/lanskyFreq.wav", "Filepath to frequency reference audio.")
	outPath := flag.String("o",
		"./music/result.wav", "Output audio filepath.")
	winLenP := flag.Int("stft_len", 2048,
		"STFT window size in number of samples.")
	winShiftP := flag.Int("stft_shift", 100,
		"STFT window shift in number of samples.")
	wP := flag.Int("w", 5,
		"Number of DFT filters for frequnecy shaping.")
	flag.Parse()

	winLen := *winLenP
	winShift := *winShiftP
	w := *wP

	// Embed greeting into binary when building:
	greetBox := packr.NewBox("../../packrbox")
	greeting, _ := greetBox.Find("greeting.txt")
	fmt.Printf("\n%s\n", greeting)

	freqWav, freqFormat := shapee.ImportWavAudio(*inFreqPath)
	ampWav, ampFormat := shapee.ImportWavAudio(*inAmpPath)

	if freqFormat.NumChannels != ampFormat.NumChannels {
		err := "error: number of channels do not match."
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	if freqFormat.SampleRate != ampFormat.SampleRate {
		err := "error: sample rates do not match."
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Let's get shaping!\n\n")

	result := make([][]float64, len(freqWav))
	for channel := range freqWav {
		fmt.Printf("Channel %d:\n", channel+1)

		fmt.Printf("   Frequency reference audio prep:\n")
		freqSTFT := shapee.ComputeSTFT(freqWav[channel], winShift, winLen)
		freqMag, freqPhase := shapee.ComplexToPolar(freqSTFT)

		fmt.Printf("   Amplitude reference audio prep:\n")
		ampSTFT := shapee.ComputeSTFT(ampWav[channel], winShift, winLen)
		ampMag, _ := shapee.ComplexToPolar(ampSTFT)

		shaper := &shapee.Shaper{ampMag, freqMag, freqPhase, w}
		resultMag, resultPhase := shaper.FreqShape()
		resultSTFT := shapee.PolarToComplex(resultMag, resultPhase)

		result[channel] = shapee.ComputeISTFT(resultSTFT, winShift)
	}
	shapee.ExportWavAudio(result, freqFormat, *outPath)
}
