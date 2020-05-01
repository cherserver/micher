package main

import (
	"fmt"
	"github.com/cherserver/micher/service/application"
	"github.com/cherserver/micher/service/environment"
	"os"
)

func main() {
	var err error
	var env = environment.New()
	if err = env.Init(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Failed to initialize environment:", err)
		os.Exit(1)
	}

	app := application.New(env)
	if err = app.Init(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Failed to initialize application:", err)
		os.Exit(1)
	}
}
