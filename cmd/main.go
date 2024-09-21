package main

import (
	"os"

	"github.com/mojganchakeri/whatsapp-manager/internal/app"
)

func main() {
	err := app.Run()
	if err != nil {
		os.Exit(1)
	}
}
