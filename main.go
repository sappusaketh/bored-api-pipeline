package main

import (
	"os"

	"github.com/sappusaketh/bored-api-pipeline/internal/application"
)

func main() {
	env := os.Getenv("ENV")
	if env == "" {
		env = "local"
	}
	application.New(env).Run()
}
