package main

import (
	"encoding/json"
	"example/win-app/bridge"
	"example/win-app/common"
	"log"

	"github.com/canpacis/go-webview2"
)

var preload = `
class WebView2Bridge extends EventTarget {
	#handlers = {};
  constructor() {
		super()

    window.addEventListener("message", (event) => {
      if (event.data.composer === "main-thread") {
        const handler = this.#handlers[event.data.id];
				this.dispatchEvent(new CustomEvent("message", { detail: event }))

        if (handler) {
          if(event.data.error) {
            handler.reject(event.data.payload);
          }else {
            handler.resolve(event.data.payload);
          }
        }

				delete this.#handlers[event.data.id]
      }
    });
  }

  async request(request, params = []) {
    return new Promise((resolve, reject) => {
      if ("chrome" in window && "webview" in window.chrome) {
        const uuidv4 = () => {
          return ([1e7] + -1e3 + -4e3 + -8e3 + -1e11).replace(/[018]/g, (c) =>
            (c ^ (crypto.getRandomValues(new Uint8Array(1))[0] & (15 >> (c / 4)))).toString(16)
          );
        };
        const id = uuidv4();

        window.chrome.webview.postMessage(JSON.stringify({ id, request, params: JSON.stringify(params) }));
        this.#handlers[id] = { resolve, reject };
      } else {
        console.warn("There is no webview context");
      }
    });
  }
}

`

func msgcb(msg string, w webview2.WebView) {
	interop := common.Interop{Window: w}
	question := common.Question{}
	err := json.Unmarshal([]byte(msg), &question)

	answer := common.Answer{
		Composer: "main-thread",
		Payload:  "",
	}

	if err != nil {
		answer.Payload = "bad request"
		answer.Error = true
	}
	answer.ID = question.ID

	request := bridge.Bridges[question.Request]

	if request == nil {
		answer.Payload = "unknown request"
		answer.Error = true
	} else {
		result, err := request(question.Params, interop)

		answer.Error = err != nil
		if answer.Error {
			answer.Payload = err.Error()
		} else {
			answer.Payload = result
		}
	}

	interop.Send(answer)
}

func main() {
	w := webview2.NewWithOptions(webview2.WebViewOptions{
		Debug:     true,
		AutoFocus: true,
		WindowOptions: webview2.WindowOptions{
			Title: "Minimal webview example",
		},
		MessageCallback: msgcb,
	})
	if w == nil {
		log.Fatalln("Failed to load webview.")
	}
	defer w.Destroy()
	w.SetSize(800, 600, webview2.HintNone)
	w.Navigate("https://en.m.wikipedia.org/wiki/Main_Page")
	w.Init(preload)
	w.Run()
}
