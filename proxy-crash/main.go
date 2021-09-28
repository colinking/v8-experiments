package main

import (
	"fmt"

	v8 "github.com/airplanedev/v8go"
)

var code = `
	new Proxy({}, {
		get: function () {
			foobar()
			return "hello!"
		},
	})
`

func main() {
	iso := v8.NewIsolate()
	defer iso.Dispose()

	global := v8.NewObjectTemplate(iso)
	err := global.Set("foobar", v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		v8v, _ := v8.NewValue(iso, "callback")
		return v8v
	}))
	if err != nil {
		panic(fmt.Errorf("setting global: %w", err))
	}

	v8ctx := v8.NewContext(iso, global)
	defer v8ctx.Close()

	v8v, err := v8ctx.RunScript(code, "main.js")
	if err != nil {
		panic(fmt.Errorf("running script: %w", err))
	}

	s, err := v8.JSONStringify(v8ctx, v8v)
	if err != nil {
		panic(fmt.Errorf("stringifying value: %w", err))
	}
	fmt.Println(s)
}
