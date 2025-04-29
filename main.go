package main

import (
	"os"

	"github.com/lunagic/hera/hera"
)

func main() {
	hera.Start(os.Args[1:]...)
}
