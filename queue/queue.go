package main

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kinmod-ui/musicIshSomethingIdk/util"
)

type queue struct {
	arr [][]float64
}

func (q *queue) Pop() []float64 {
	ret := q.arr[0]

	q.arr = q.arr[1:]

	return ret
}

func (q *queue) Push(x []float64) {
	q.arr = append(q.arr, x)
}

var Datastream = make(map[string]*queue)

func main() {

	util.Mylog.Println("This is the queue speaking hehe")

	http.HandleFunc("/ws-consume", wsSocket)
	http.HandleFunc("/ws", incomingMessages)

	if err := http.ListenAndServe(":8081", nil); err != nil {
		util.Mylog.Println("Started on port 8081")
	}

}

func incomingMessages(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		util.Mylog.Println("Error with netdial: ", err)
		return
	}
	defer c.Close()

	for {
		msg := &[]float64{}
		err := c.ReadJSON(msg)

		if err != nil {
			util.Mylog.Println("Got error ", err)
			if websocket.IsCloseError(err, 1000) {
				util.Mylog.Println("Read all data. Adios")
				Datastream["hehe"].Push([]float64{})
			}
			break
		}

		if _, ok := Datastream["hehe"]; !ok {
			Datastream["hehe"] = &queue{}
		}
		Datastream["hehe"].Push(*msg)
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options

func wsSocket(w http.ResponseWriter, r *http.Request) {

	util.Mylog.Println("I am here started")

	var err error
	w.Header().Set("Content-Type", "application/json")

	c, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		util.Mylog.Println("upgrade : ", err)
		return
	}

	defer c.Close()

	for {

		if _, ok := Datastream["hehe"]; !ok {
			time.Sleep(1 * time.Second)
			continue
		}
		if len(Datastream["hehe"].arr) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}

		m := Datastream["hehe"].Pop()
		//if m == "EOF" {
		//c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		//break
		//}
		if len(m) > 0 {
			c.WriteJSON(m)
			time.Sleep(10 * time.Millisecond)
		} else {
			c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		}
	}
}
