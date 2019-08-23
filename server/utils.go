package main

import (
	"encoding/json"
	"fmt"
	"github.com/alex023/eventbus"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const (
	ACTION_AUTH = "Auth"
	ACTION_SET  = "Set"
	ACTION_GET  = "Get"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	HandshakeTimeout: 5 * time.Second,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

type CloudText struct {
	value string
	bus   *eventbus.Bus
}

type CloudTextResult struct {
	Value string
	Err   string
}

type CloudTextPackage struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Action   string `json:"action"`
	Data     string `json:"data"`
}

func (cloudText *CloudText) setText(value string) (err error) {
	if cloudText.value==value {
		return
	}
	cloudText.value = value
	err = cloudText.bus.Push(EVENT_SET, value)
	return
}

func (cloudText *CloudText) getText() (value string, err error) {
	value = cloudText.value
	return
}

func (cloudText *CloudText) handleText(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Content-Type", "application/json")
	if request.Header.Get("username") != *username || request.Header.Get("password") != *password {
		writer.WriteHeader(http.StatusNetworkAuthenticationRequired)
		return
	}
	switch request.Method {
	case http.MethodGet:
		value, err := cloudText.getText()
		result := CloudTextResult{Value: value}
		if err != nil {
			log.Println(err)
			result.Err = fmt.Sprint(err)
		} else {
			log.Println(request.RemoteAddr, "\tget")
			result.Err = "null"
		}
		b, err := json.Marshal(result)
		if err != nil {
			log.Println(err)
			return
		}
		_, err = writer.Write(b)
		if err != nil {
			log.Println(err)
			return
		}
	case http.MethodPost:
		value := request.FormValue("value")
		err := cloudText.setText(value)
		result := CloudTextResult{}
		if err != nil {
			log.Println(err)
			result.Err = fmt.Sprint(err)
		} else {
			log.Println(request.RemoteAddr, "\tset\t", value)
			result.Err = "null"
		}
		b, err := json.Marshal(result)
		if err != nil {
			log.Println(err)
			return
		}
		_, err = writer.Write(b)
		if err != nil {
			log.Println(err)
			return
		}
	default:
		writer.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (cloudText *CloudText) handleWebsocketUpgrade(writer http.ResponseWriter, request *http.Request) {
	auth := false
	ws, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	log.Println("A new websocket connection:", ws.RemoteAddr().String())
	go func() {
		time.Sleep(3 * time.Second)
		if !auth {
			auth = true
			log.Println("websocket auth failed:", ws.RemoteAddr())
			if err := ws.Close(); err != nil {
				log.Println("Fail to close websocket ", ws.RemoteAddr().String(), "\n", err)
			}
		}
	}()
	go func() {
		for auth != true {
			var pack CloudTextPackage
			if err := ws.ReadJSON(&pack); err != nil {
				log.Println("Fail decode json from websocket ", ws.RemoteAddr().String(), "\n", err)
				break
			}
			if pack.Action == ACTION_AUTH {
				response := CloudTextPackage{Action: ACTION_AUTH, ID: pack.ID}
				if pack.Username == *username && pack.Password == *password {
					auth = true
					response.Data = "success"
					if _, err := cloudText.bus.Subscribe(func(message interface{}) {
						resp := CloudTextPackage{Action: ACTION_SET, Data: message.(string)}
						if err := ws.WriteJSON(resp); err != nil {
							log.Println("Fail send response to websocket ", ws.RemoteAddr().String(), "\n", err)
							return
						}
					}, EVENT_SET); err != nil {
						log.Println("Fail add subscribe(", EVENT_SET, ") to websocket ", ws.RemoteAddr().String(), "\n", err)
					} else {
						log.Println("Add subscribe(", EVENT_SET, ") to websocket ", ws.RemoteAddr().String())
					}
					go func() {
						for {
							req := CloudTextPackage{}
							if err := ws.ReadJSON(&req); err != nil {
								_ = ws.Close()
								break
							}
							log.Println("Read message from ", ws.RemoteAddr().String(), ":\n", req)
							switch req.Action {
							case ACTION_SET:
								_ = cloudText.setText(req.Data)
							case ACTION_GET:
								req.Data=cloudText.value
								_ = ws.WriteJSON(req)
							}
						}
					}()
				} else {
					response.Data = "failed"
				}
				if err := ws.WriteJSON(response); err != nil {
					log.Println("Fail send response to websocket ", ws.RemoteAddr().String(), "\n", err)
					continue
				} else {
					log.Println("Send response to websocket ", ws.RemoteAddr().String())
				}
			}
		}
	}()
}
