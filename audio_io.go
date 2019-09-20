package shapee

import (
	"fmt"
	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"math"
	"os"
)

// ImportWavAudio imports a .wav file and returns the audio waveform as well
// as its format information.
func ImportWavAudio(inputPath string) ([][]float64, *audio.Format) {
	fmt.Printf("Importing %s\n", inputPath)

	f, err := os.Open(inputPath)
	check(err)

	decoder := wav.NewDecoder(f)

	fmt.Printf("   Loading into buffer...\n")
	buf, err := decoder.FullPCMBuffer()
	check(err)

	wave := make([][]float64, buf.Format.NumChannels)

	fmt.Printf("   Converting to floating point...\n")
	for channel := range wave {
		wave[channel] = make([]float64, buf.NumFrames())
		for j := range wave[channel] {
			wave[channel][j] = float64(buf.Data[j*len(wave)+channel]) /
				math.Pow(2, float64(buf.SourceBitDepth)-1)
		}
	}

	err = f.Close()
	check(err)
	fmt.Printf("   Done.\n\n")

	return wave, buf.Format
}

// ExportWavAudio exports a .wav file from a 2D float64 array.
func ExportWavAudio(wave [][]float64, format *audio.Format, outPath string) {
	fmt.Printf("Exporting result...\n")
	f, err := os.Create(outPath)
	check(err)

	encoder := wav.NewEncoder(f, format.SampleRate, 16, format.NumChannels, 1)

	buf := audio.IntBuffer{
		Format:         format,
		Data:           make([]int, len(wave[0])*len(wave)),
		SourceBitDepth: 16,
	}

	fmt.Printf("   Converting to int...\n")
	for i := range wave[0] {
		for channel := range wave {
			buf.Data[i*len(wave)+channel] = int(wave[channel][i] *
				math.Pow(2, float64(16)))
		}
	}

	fmt.Printf("   Writing to file...\n")
	err = encoder.Write(&buf)
	check(err)

	err = encoder.Close()
	check(err)

	err = f.Close()
	check(err)

	fmt.Printf("   Done.\n\nResult exported to %s\n\n", outPath)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
