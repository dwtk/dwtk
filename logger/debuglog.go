package logger

import (
	"io/ioutil"
	"log"
	"os"
)

var (
	Debug = log.New(ioutil.Discard, "D: ", 0)
)

func EnableDebug() {
	Debug.SetOutput(os.Stderr)
}
