package main

import (
	"os"

	"go-home/call"
	_ "go-home/help"
)

func main() {
	args := os.Args

	if len(args) >= 2 {
		call.Start(args[1])
	}
}
