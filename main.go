package main

import (
	"encoding/json"
	"math"
	"net/http"
	"net/url"

	pkg "github.com/DylanMeeus/GoAudio/wave"

	"github.com/gorilla/websocket"

	"github.com/kinmod-ui/musicIshSomethingIdk/util"
)

func dft(arr []pkg.Frame, samplingFreq int) []float64 {

	length := len(arr)

	freqArr := []float64{}

	for i := 1; i < length; i++ {

		realPart := 0.0
		imagPart := 0.0

		for idx, amp := range arr {
			realPart += float64(amp) * (math.Cos(-(2.0 * (math.Pi) * (float64(i) * (float64(idx))) / (float64(length)))))
			imagPart += float64(amp) * (math.Sin(-(2.0 * (math.Pi) * (float64(i) * (float64(idx))) / (float64(length)))))
		}

		freq := math.Hypot(realPart, imagPart)
		freqArr = append(freqArr, freq)
	}

	return freqArr

}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options

var wave pkg.Wave

func main() {

	http.HandleFunc("/ws", sendToQueue)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		util.Mylog.Println("Started on port 8080")
	}

}

func sendToQueue(w http.ResponseWriter, r *http.Request) {

	util.Mylog.Println("I am here started")

	var err error
	wave, err = pkg.ReadWaveFile("./resources/egwav.wav")

	if err != nil {
		panic("Could not parse wave file")
	}

	u := url.URL{Scheme: "ws", Host: "localhost:8081", Path: "/ws"}
	util.Mylog.Println("connecting to ", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		util.Mylog.Println("Error with netdial: ", err)
		return
	}
	defer c.Close()

	util.Mylog.Printf("Read %v samples\n", len(wave.Frames))

	if err != nil {
		util.Mylog.Println("upgrade : ", err)
		return
	}

	start := 0
	duration := len(wave.Frames) / wave.SampleRate
	framesPerSecond := len(wave.Frames) / duration
	// 1052215/12/100
	toSend := framesPerSecond / 50

	for {
		m := makePositive(logarise(dft(wave.Frames[start:start+toSend], wave.WaveFmt.SampleRate)))
		c.WriteJSON(m)

		if len(wave.Frames) >= start+toSend {
			start += toSend

		} else {
			c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))

			msg, err := json.Marshal("Goodbye")
			if err != nil {
				util.Mylog.Println(err)
			}

			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			n, err := w.Write(msg)
			if err != nil {
				util.Mylog.Println(err)
			}
			util.Mylog.Println(n)
			break
		}
	}
}

func makePositive(arr []float64) []float64 {
	mini := 1e9

	for _, amp := range arr {
		mini = math.Min(mini, float64(amp))
	}

	arr2 := []float64{}

	if mini < 0 {
		for _, amp := range arr {
			arr2 = append(arr2, (float64(amp))-mini)
		}

		return arr2
	} else {
		for _, amp := range arr {
			arr2 = append(arr2, (float64(amp)))
		}

		return arr2
	}
}

func normalise(arr []float64) []float64 {

	maxi := 0.0
	for _, amp := range arr {
		maxi = math.Max(maxi, amp)
	}

	arr2 := []float64{}

	for _, amp := range arr {
		arr2 = append(arr2, amp/maxi)
	}

	//util.Mylog.Println(arr2)
	return arr2
}

func logarise(arr []float64) []float64 {

	arr2 := []float64{}

	for _, amp := range arr {
		arr2 = append(arr2, math.Log(amp))
	}

	//util.Mylog.Println(arr2)
	return arr2
}
