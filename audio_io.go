package shapee

import (
	"fmt"
	"github.com/youpy/go-wav"
	"math"
	"os"
	"sync"
)

// ImportAudio imports a .wav file and returns the audio waveform as well
// as the bitrate, samplerate, and number of channels.
func ImportAudio(inputPath string) ([][]float64, *wav.WavFormat) {
	numSamples := uint32(math.MaxUint32)

	//blockAlign := 4
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

	//if int(wavFormat.BlockAlign) != blockAlign {
	//fmt.Println("Block align is invalid ", wavFormat.BlockAlign)
	//}

	samples, err := reader.ReadSamples(numSamples)
	wavSamples := make([][]float64, 2)

	if wavFormat.NumChannels == 1 {
		wavSamples[0] = make([]float64, len(samples))
		for i, curr_sample := range samples {
			wavSamples[0][i] = reader.FloatValue(curr_sample, 0)
		}
	} else {
		wavSamples[0] = make([]float64, len(samples))
		wavSamples[1] = make([]float64, len(samples))
		for i, curr_sample := range samples {
			wavSamples[0][i] = reader.FloatValue(curr_sample, 0)
			wavSamples[1][i] = reader.FloatValue(curr_sample, 1)
		}
	}

	//return wavSamples, wavFormat.BitsPerSample, wavFormat.SampleRate, wavFormat.NumChannels
	return wavSamples, wavFormat
}

// ExportAudio exports a .wav file from a mono/stereo floating point waveform
func ExportAudio(wave [][]float64, wavFormat *wav.WavFormat, exportPath string) {
	f, err := os.Create(exportPath)
	if err != nil {
		panic(err)
	}

	fmt.Println("Writing to file...")
	wavWriter := wav.NewWriter(f, uint32(len(wave[0])), wavFormat.NumChannels, wavFormat.SampleRate, wavFormat.BitsPerSample)
	samples := make([]wav.Sample, len(wave[0]))

	var waitGroup sync.WaitGroup

	for i := 0; i < len(wave[0]); i += 1000 {
		waitGroup.Add(1)
		iGo := i
		go func() {
			defer waitGroup.Done()

			var waveSlice [][]float64
			if iGo+1000 >= len(wave[0]) {
				waveSlice = make([][]float64, 2)
				waveSlice[0] = make([]float64, len(wave[0])-iGo)
				waveSlice[1] = make([]float64, len(wave[1])-iGo)
				copy(waveSlice[0], wave[0][iGo:])
				copy(waveSlice[1], wave[1][iGo:])
			} else {
				waveSlice = make([][]float64, 2)
				waveSlice[0] = make([]float64, 1000)
				waveSlice[1] = make([]float64, 1000)
				copy(waveSlice[0], wave[0][iGo:iGo+1000])
				copy(waveSlice[1], wave[1][iGo:iGo+1000])
			}

			for j := range waveSlice[0] {
				values := [2]int{0, 0}
				for k := 0; k < int(wavFormat.NumChannels); k++ {
					values[k] = int(waveSlice[k][j] * (math.Pow(2, float64(wavFormat.BitsPerSample))))
				}
				samples[iGo+j] = wav.Sample{values}
			}
		}()
	}

	waitGroup.Wait()

	fmt.Println("WavWriter working...")

	wavWriter.WriteSamples(samples)
}

//
//
//
//
//
