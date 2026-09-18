package main

import (
	"encoding/json"
	_ "unsafe"

	_ "frm/napi-go/module"

	"frm/napi-go"
)

func main() {}

//go:linkname Register frm/napi-go/module.Register
func Register(env napi.EnvType, export *napi.Object) {
	f, _ := napi.GoFuncOf(env, Test)
	export.Set("goFunc", f)
}

// Go promitives from javascript
type Value struct {
	String string            `napi:"value"`
	Int    int               `napi:"int"`
	Float  float64           `napi:"float"`
	Bool   bool              `napi:"bool"`
	Nil    any               `napi:"nil"`
	Mapped map[string]string `napi:"map"`
	Slice  []string          `napi:"array"`
}

func Test(str string, v *Value) {
	d, _ := json.MarshalIndent(v, "", "  ")
	println(string(d))
}
