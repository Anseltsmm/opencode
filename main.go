package main

import (
	"github.com/Anseltsmm/azkia/cmd"
	"github.com/Anseltsmm/azkia/internal/logging"
)

func main() {
	defer logging.RecoverPanic("main", func() {
		logging.ErrorPersist("Application terminated due to unhandled panic")
	})

	cmd.Execute()
}
