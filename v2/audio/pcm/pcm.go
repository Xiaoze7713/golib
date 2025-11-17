package pcm

import (
	"encoding/binary"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

func Conv(pcmFile string, outFile string) {
	in, err := os.Open(pcmFile)
	if err != nil {
		log.Fatal(err)
	}

	out, err := os.Create(outFile)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// 8 kHz, 16 bit, 1 channel, WAV.
	cnt, err := ioutil.ReadFile(pcmFile)
	println("len = ", float64(len(cnt))*8/(16000*16))
	e := wav.NewEncoder(out, 16000, 16, 1, 1)

	// Create new audio.IntBuffer.
	audioBuf, err := newAudioIntBuffer(in)
	if err != nil {
		log.Fatal(err)
	}
	// Write buffer to output file. This writes a RIFF header and the PCM chunks from the audio.IntBuffer.
	if err := e.Write(audioBuf); err != nil {
		log.Fatal(err)
	}
	if err := e.Close(); err != nil {
		log.Fatal(err)
	}
	println(float64(e.WrittenBytes) * 8 / (16000 * 16))
}

func newAudioIntBuffer(r io.Reader) (*audio.IntBuffer, error) {
	buf := audio.IntBuffer{
		Format: &audio.Format{
			NumChannels: 1,
			SampleRate:  16000,
		},
	}
	for {
		var sample int16
		err := binary.Read(r, binary.LittleEndian, &sample)
		switch {
		case err == io.EOF:
			return &buf, nil
		case err != nil:
			return nil, err
		}
		buf.Data = append(buf.Data, int(sample))
	}
}
