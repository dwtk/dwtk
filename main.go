package main

import (
	"os"

	"github.com/dwtk/dwtk/internal/cmd"
)

func main() {
	err1 := cmd.RootCmd.Execute()
	err2 := cmd.Close()
	if err2 != nil {
		cmd.RootCmd.Println("Error:", err2.Error())
	}
	if err1 != nil || err2 != nil {
		os.Exit(1)
	}
}
