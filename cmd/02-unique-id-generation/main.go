package main

import (
	"github.com/jtribble/fly-io-dist-sys/pkg/maelstrom"
)

func main() {
	node := maelstrom.NewNode(&GenerateHandler{})
	node.RunUntilInterrupted()
}
