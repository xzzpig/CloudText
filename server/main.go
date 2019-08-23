package main

import (
	"flag"
	"github.com/alex023/eventbus"
	"log"
	"net/http"
)

var (
	username = flag.String("u", "cloudtext", "username")
	password = flag.String("p", GetRandomString(8), "password")
	host     = flag.String("h", "0.0.0.0:23451", "The http server host on.")
	httpPath = flag.String("t", "/cloudtext/text", "")
	wsPath   = flag.String("w", "/cloudtext/ws", "")
)

const (
	EVENT_SET = "EVENT_SET"
)

func main() {
	flag.Parse()
	cloudText := CloudText{bus: eventbus.Default()}
	http.HandleFunc(*httpPath, cloudText.handleText)
	log.Println("Add handler for http on ", *httpPath)
	//http.Handle(*wsPath, websocket.Handler(cloudText.handleWebsocketUpgrade))
	http.HandleFunc(*wsPath, cloudText.handleWebsocketUpgrade)
	log.Println("Add handler for ws   on ", *wsPath)
	//wsHandler := websocket.Handler(cloudText.handleWebsocketUpgrade)
	//http.HandleFunc(*wsPath, func(writer http.ResponseWriter, request *http.Request) {
	//	if request.Header.Get("Origin")==""{
	//		request.Header.Set("Origin","http://localhost")
	//	}
	//	wsHandler.ServeHTTP(writer, request)
	//})
	log.Println("Starting http server on", *host)
	log.Println("username:", *username, ",password:", *password)
	_, _ = cloudText.bus.Subscribe(func(message interface{}) {

	}, EVENT_SET)
	err := http.ListenAndServe(*host, nil)
	if err != nil {
		log.Panicln("Failed to start http server\n", err)
	}
}
