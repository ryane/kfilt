package main

import (
	"os"

	"github.com/ryane/kfilt/cmd"
)

func main() {
	cmd.Execute(os.Args[1:])
}
