package bridge

import (
	"example/win-app/common"
)

type BridgeFunc func(params string, interop common.Interop) (interface{}, error)

var Bridges = map[string]BridgeFunc{
	"hello": func(params string, interop common.Interop) (interface{}, error) {
		return "world", nil
	},
}
