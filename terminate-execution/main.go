package main

import (
	"fmt"
	"time"

	// This example uses our fork for darwin_arm64 support, but the issue
	// should surface using the upstream master, too:
	// v8 "rogchap.com/v8go"
	v8 "github.com/airplanedev/v8go"
)

func main() {
	iso := v8.NewIsolate()
	defer func() {
		if !iso.IsExecutionTerminating() {
			iso.Dispose()
		}
	}()

	v8ctx := v8.NewContext(iso)
	defer v8ctx.Close()

	go func() {
		_, err := v8ctx.RunScript(`
		while (true) {}
	`, "main.js")
		if err != nil {
			panic(fmt.Errorf("running script: %w", err))
		}
	}()

	time.After(10 * time.Millisecond)
	iso.TerminateExecution()
}
