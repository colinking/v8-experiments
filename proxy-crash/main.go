package main

import (
	"fmt"

	// This example uses our fork for darwin_arm64 support, but the issue
	// should surface using the upstream master, too:
	// v8 "rogchap.com/v8go"
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

// It's not an issue with calling JSON.stringify on this proxy, since
// the following runs fine:
//
// var code = `
// 	JSON.stringify(new Proxy({}, {
// 		get: function () {
// 			foobar()
// 			return "hello!"
// 		},
// 	}))
// `

// It also is not an issue with function closures, since the following also
// runs fine:
//
// var code = `
// 	const foo = {
// 		get: function () {
// 			foobar()
// 			return "hello!"
// 		}
// 	};
// 	foo
// `

func main() {
	iso := v8.NewIsolate()
	defer iso.Dispose()

	// You can uncomment this to capture the panic thrown by `JSONStringify`:
	//
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		fmt.Printf("Caught a panic: %v\n", r)
	// 	}
	// }()

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

	s, err := v8.JSONStringify(nil, v8v)
	if err != nil {
		panic(fmt.Errorf("stringifying value: %w", err))
	}
	fmt.Println(s)
}
