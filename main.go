package main

import (
	"github.com/go-mods/tagsvar/cmd"
	"os"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
