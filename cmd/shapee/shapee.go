package main

import (
	"flag"
	"fmt"
	"github.com/isamilefchik/shapee"
	"os"
)

func main() {
	inFreqPath := flag.String("if",
		"./music/objekt.wav", "Filepath to frequency reference audio.")
	inAmpPath := flag.String("ia",
		"./music/olixl.wav", "Filepath to amplitude reference audio.")
	outPath := flag.String("o",
		"./music/result.wav", "Output audio filepath.")
	winLen := *(flag.Int("stft_len", 256,
		"STFT window size in number of samples."))
	winShift := *(flag.Int("stft_shift", 50,
		"STFT window shift in number of samples."))
	w := *(flag.Int("w", 4,
		"Number of DFT filters for frequnecy shaping."))
	flag.Parse()

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

		pack := &shapee.ShaperPack{ampMag, freqMag, freqPhase, w}
		resultMag, resultPhase := shapee.FreqShape(pack)
		resultSTFT := shapee.PolarToComplex(resultMag, resultPhase)

		result[channel] = shapee.ComputeISTFT(resultSTFT, winShift)
	}
	shapee.ExportWavAudio(result, freqFormat, *outPath)
}
