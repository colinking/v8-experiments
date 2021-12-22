package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	v8 "rogchap.com/v8go"
)

func main() {
	iso := v8.NewIsolate()

	global := v8.NewObjectTemplate(iso)
	cb := v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		fmt.Fprintf(os.Stderr, "callback called\n")
		// This shows that Go callbacks block termination, until they return.
		time.Sleep(5 * time.Second)
		fmt.Fprintf(os.Stderr, "callback awake\n")
		return v8.Null(info.Context().Isolate())
	})
	if err := global.Set("callback", cb); err != nil {
		panic(err)
	}

	v8ctx := v8.NewContext(iso, global)
	defer func() {
		fmt.Printf("v8 ctx closing...\n")
		v8ctx.Close()
		fmt.Printf("v8 ctx closed\n")
	}()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Fprintf(os.Stderr, "starting script...\n")
		v, err := v8ctx.RunScript(`
			callback()
		`, "main.js")
		fmt.Fprintf(os.Stderr, "script returned: %+v\n", v)
		if err != nil {
			panic(fmt.Errorf("running script: %w", err))
		}
	}()

	time.After(100 * time.Millisecond)
	iso.TerminateExecution()
	fmt.Printf("termination started...\n")

	wg.Wait()
}
