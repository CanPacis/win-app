package common

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/canpacis/go-webview2"
)

type Interop struct {
	Window webview2.WebView
}

type Answer struct {
	ID       string      `json:"id"`
	Composer string      `json:"composer"`
	Payload  interface{} `json:"payload"`
	Error    bool        `json:"error"`
}

type Question struct {
	ID      string `json:"id"`
	Request string `json:"request"`
	Params  string `json:"params"`
}

func (i Interop) Send(message Answer) {
	encoded, err := json.Marshal(message)
	if err != nil {
		log.Fatalln(err.Error())
	}

	i.Window.Dispatch(func() {
		i.Window.Eval(fmt.Sprintf("window.postMessage(%s)", encoded))
	})
}

func (i Interop) SendArbitrary(message interface{}, id string) {
	answer := Answer{
		Error:    false,
		ID:       id,
		Payload:  message,
		Composer: "main-thread",
	}

	i.Send(answer)
}
