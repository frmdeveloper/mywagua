package main

import (
	"fmt"
	"time"
	_ "unsafe"

	_ "frm/napi-go/module"

	"frm/napi-go"
)

var waitTime = time.Second * 3

//go:linkname Register frm/napi-go/module.Register
func Register(env napi.EnvType, export *napi.Object) {
	fn, _ := napi.CreateFunction(env, "", func(ci *napi.CallbackInfo) (napi.ValueType, error) {
		var Test *napi.String
		return napi.CreateAsyncWorker(env,
			func(env napi.EnvType) {
				fmt.Printf("Wait %s\n", waitTime)
				<-time.After(waitTime)
				println("wait time done on exec func")
				Test, _ = napi.CreateString(env, "Test")
			},
			func(env napi.EnvType, Resolve, Reject func(value napi.ValueType)) {
				println("Done, called done function")
				Resolve(Test)
			})
	})
	export.Set("async", fn)
}

func main() {}
