package help

import (
	"flag"
	"fmt"
	"os"
)

// G4755 is a number of a train, as soon it will send me to city of nanchang where is near to my homeland.
var configfile = flag.String("train", "G4775", "You need to set train number that you will be ride")

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [command]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
}
