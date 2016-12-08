package main

import (
	"flag"
	"os"

	"github.com/b1naryth1ef/wiretap/lib"
)

func main() {
	flag.Parse()

	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	rec := wiretap.NewRecording(flag.Args()[0], wiretap.MODE_REPLAY)
	rec.Open()
	rec.Play(os.Stdout)
}
