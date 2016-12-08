package main

import (
	"flag"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/b1naryth1ef/wiretap/lib"
	"github.com/kr/pty"
)

var shell = flag.String("shell", "/bin/bash", "Shell to use")

func main() {
	flag.Parse()

	shellSession := exec.Command(*shell)
	ptySession, _ := pty.Start(shellSession)

	rec := wiretap.NewRecording("test.rec", wiretap.MODE_RECORD)
	tap := wiretap.NewWiretapper(rec, ptySession, os.Stdout)
	tap.Run()

	go func() {
		io.Copy(ptySession, os.Stdin)
	}()

	<-tap.Complete
	_, err := shellSession.Process.Wait()
	log.Printf("Failed to close shell: %v", err)
}
