package help

import (
	"flag"
	"fmt"
	"os"
)

// G4755 is a number of a train, as soon it will send me to city of nanchang where is near to my homeland.
var train = flag.String("train", "G4775", "You need to set train number that you will be ride")
var date = flag.String("date", "2017-01-22", "Special start date")

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [command] [option]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
}
