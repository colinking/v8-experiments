package main

import (
	"fmt"

	// This example uses our fork for darwin_arm64 support, but the issue
	// should surface using the upstream master, too:
	// v8 "rogchap.com/v8go"
	v8 "github.com/airplanedev/v8go"
)

func main() {
	// This will return an `undefined` value, rather than an empty object literal:
	run(`{}`)

	// If you wrap the object literal in a variable declarion, they work fine:
	run(`const foo = {}; foo`)
	run(`const foo = {
		"foo": "bar"
	}; foo`)

	// However, if you have a key-value in the object literal then it fails with:
	//
	//   panic: running script: SyntaxError: Unexpected token ':'
	//
	//   goroutine 1 [running]:
	//   main.main()
	//   	/Users/colin/dev/colinking/v8-experiments/object-literals/main.go:27 +0x264
	//   exit status 2
	//
	run(`{
		"foo": "bar"
	}`)
}

func run(code string) {
	fmt.Printf("> Code:\n%s\n", code)

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
	fmt.Printf("> Output:\n%s\n\n", s)
}
