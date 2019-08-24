package shapee

import (
	"fmt"
	"github.com/youpy/go-wav"
	"math"
	"os"
)

// ImportAudio imports a .wav file and returns the audio waveform as well
// as the bitrate, samplerate, and number of channels.
func ImportAudio(input_path string) ([]float64, uint16, uint32, uint16) {
	number_of_samples := uint32(math.MaxUint32)

	blockAlign := 4
	file, err := os.Open(input_path)
	if err != nil {
		panic(err)
	}

	reader := wav.NewReader(file)
	wavformat, err_rd := reader.Format()
	if err_rd != nil {
		panic(err_rd)
	}

	if wavformat.AudioFormat != wav.AudioFormatPCM {
		panic("Audio format is invalid ")
	}

	if int(wavformat.BlockAlign) != blockAlign {
		fmt.Println("Block align is invalid ", wavformat.BlockAlign)
	}

	samples, err := reader.ReadSamples(number_of_samples)
	wav_samples := make([]float64, 0)

	for _, curr_sample := range samples {
		wav_samples = append(wav_samples, reader.FloatValue(curr_sample, 0))
	}

	return wav_samples, wavformat.BitsPerSample, wavformat.SampleRate, wavformat.NumChannels
}
