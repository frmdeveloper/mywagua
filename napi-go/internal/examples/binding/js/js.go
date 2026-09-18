package main

import (
	"frm/napi-go"
	entry "frm/napi-go/module/binding"
)

func init() {
	entry.Register(func(env napi.EnvType, export *napi.Object) {
		str, _ := napi.CreateString(env, "from golang")
		export.Set("from", str)
	})
}

func main() {}
