package main

import (
	"fmt"

	// This example uses our fork for darwin_arm64 support, but the issue
	// should surface using the upstream master, too:
	// v8 "rogchap.com/v8go"
	v8 "github.com/airplanedev/v8go"
)

var code = `
	{
		"foo": "bar"
	}
`

func main() {
	iso := v8.NewIsolate()
	defer iso.Dispose()

	v8ctx := v8.NewContext(iso)
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
