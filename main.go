package main

import (
	"os"

	"go-home/call"
	_ "go-home/help"
	"go-home/util"
)

func main() {
	args := os.Args
	date := util.FomatNowDate()

	if len(args) >= 2 {
		if len(args) == 3 {
			call.Start(args[1], args[2])
		} else {
			call.Start(args[1], date)
		}
	}
}
