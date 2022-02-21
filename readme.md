# Windows Webview2 App Template

## Communication between runtimes

WebView2 template exposes a bridge called `WebView2Bridge`. Simply create an instance of the object and request resources you define in bridge.go file.

```go
var Bridges = map[string]BridgeFunc{
  // Defines a hello resource and returns a string as an answer.
  // If there is an error, the javascript promise will reject.
  // The key hello defines the target resource.
  // This function takes a params string as a json object
	"hello": func(params string, interop common.Interop) (interface{}, error) {
		return "world", nil
	},
}
```

```javascript
let bridge = new WebView2Bridge();

await bridge.request("hello", [
  /* params */
]); // returns "world"
```
