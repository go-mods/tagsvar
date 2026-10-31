package main

import (
	"os"

	"github.com/go-mods/tagsvar/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
