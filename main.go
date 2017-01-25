package main

import (
	"go-home/call"
	"os"
)

func main() {
	commands := call.Commands
	args := os.Args

	for _, cmd := range commands {
		if cmd.Run != nil {
			cmd.Flag.Parse(args[1:])
			args = cmd.Flag.Args()
			switch cmd.UsageLine {
			case "trainschedule":
				train := call.Train
				date1 := call.Date1
				if train == "" || date1 == "" {
					cmd.Flag.PrintDefaults()
					os.Exit(1)
				}
				args1 := []string{train, date1}
				os.Exit(cmd.Run(cmd, args1))
			case "lefttriicket":
				from := call.From
				to := call.To
				date2 := call.Date2
				if from == "" || to == "" || date2 == "" {
					cmd.Flag.PrintDefaults()
					os.Exit(1)
				}
				args2 := []string{from, to, date2}
				os.Exit(cmd.Run(cmd, args2))
			}
		}
	}
}
