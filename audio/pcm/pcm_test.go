package pcm

import (
	"io/ioutil"
	"strings"
	"testing"
)

func TestConv(t *testing.T) {
	p := "/Users/cangxiaoze/user/git/pcm_data/00f85543a7c093e2a1bb6efddf41f1c7/2022-07-21/"
	files, _ := ioutil.ReadDir(p)
	for _, f := range files {
		if len(f.Name()) > 4 && f.Name()[len(f.Name())-4:] == ".pcm" {
			name := p + "/" + f.Name()
			out := name + ".wav"
			println(name, out)
			Conv(name, out)
		}
	}
	//file := "/Users/cangxiaoze/user/git/pcm_data/00f85543a7c093e2a1bb6efddf41f1c7/2022-07-21/1658336010314-f58767fe9b114f1fa0463debd7d0921a-user.pcm"
	//out := file + ".wav"
	//Conv(file, out)
}

func TestAudioFaceLen(t *testing.T) {
	file := "/Users/cangxiaoze/user/git/audio_face_data/audio2face.csv"
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	lines := len(strings.Split(string(bytes), "\n")) - 1
	audioMs := float64(lines) * float64(1000) / 600 * 24
	//println("min", audioMs/60000)
	println("min", int(audioMs/1000)/60, "second", int(audioMs/1000)%60)
}
