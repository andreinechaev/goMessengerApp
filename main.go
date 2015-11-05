package main

import (
	"log"
	"net/http"
	"golang.org/x/net/websocket"
	"encoding/json"
)

const server = ":6969"

var active = make(map[string]*websocket.Conn)

type JSONRequest struct {
	Msg string `json:"message"`
	Name string `json:"name"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var JSON JSONRequest
	if err := decoder.Decode(&JSON); err != nil {
		resp := &JSONRequest{
			Name: "Error",
			Msg: "Couldn't decode JSON",
		}
		js, _ := json.Marshal(resp)
		w.Write(js)
	}
	active[JSON.Name] = nil
	resp := &JSONRequest{
			Name: "Message",
			Msg: "Regestered",
		}
	js, _ := json.Marshal(resp)
	w.Write(js)
}

func Echo(ws *websocket.Conn) {
	var reqJSON JSONRequest
	defer ws.Close()

	for {
		if err := websocket.JSON.Receive(ws, &reqJSON); err != nil {
			panic(err)
			resp := &JSONRequest{
				Msg: "Cannot parse request",
				Name: "Message error",
			}
			websocket.JSON.Send(ws, resp)
			return
		}
		out, _ := json.Marshal(reqJSON)
		log.Println(string(out))
		active[reqJSON.Name] = ws
		resp := &JSONRequest {
			Msg: reqJSON.Msg,
			Name: reqJSON.Name,
		}

		for _, v := range(active) {
			if err := websocket.JSON.Send(v, resp); err != nil {
				log.Println(err.Error())
			return
		}
		}
	}
}

func main() {
	log.Println("Server started on", server)
	http.Handle("/", websocket.Handler(Echo))
	http.HandleFunc("/register", Register)
	if err := http.ListenAndServe(server, nil); err != nil {
		panic(err)
	}
}
