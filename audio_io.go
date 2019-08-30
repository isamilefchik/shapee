package shapee

import (
	"fmt"
	"github.com/youpy/go-wav"
	"math"
	"os"
)

// ImportAudio imports a .wav file and returns the audio waveform as well
// as the bitrate, samplerate, and number of channels.
func ImportAudio(inputPath string) ([]float64, *wav.WavFormat) {
	numSamples := uint32(math.MaxUint32)

	blockAlign := 4
	file, err := os.Open(inputPath)
	if err != nil {
		panic(err)
	}

	reader := wav.NewReader(file)
	wavFormat, errRd := reader.Format()
	if errRd != nil {
		panic(errRd)
	}

	if wavFormat.AudioFormat != wav.AudioFormatPCM {
		panic("Audio format is invalid ")
	}

	if int(wavFormat.BlockAlign) != blockAlign {
		fmt.Println("Block align is invalid ", wavFormat.BlockAlign)
	}

	samples, err := reader.ReadSamples(numSamples)
	wavSamples := make([]float64, 0)

	for _, curr_sample := range samples {
		wavSamples = append(wavSamples, reader.FloatValue(curr_sample, 0))
	}

	//return wavSamples, wavFormat.BitsPerSample, wavFormat.SampleRate, wavFormat.NumChannels
	return wavSamples, wavFormat
}

// ExportAudio exports a .wav file from
func ExportAudio(wave []float64, wavFormat *wav.WavFormat, exportPath string) {
	// TODO: Export audio
}
