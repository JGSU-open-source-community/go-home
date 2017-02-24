package main

import (
	"os"
)

func main() {
	latitudeAndLongitude()
	commands := Commands
	args := os.Args

	for _, cmd := range commands {
		if cmd.Run != nil && cmd.Name() == args[1] {
			cmd.Flag.Parse(args[2:])
			args = cmd.Flag.Args()
			os.Exit(cmd.Run(cmd, args))
		}
	}
}
