package main

import (
	"fmt"
	"github.com/rh-ecosystem-edge/ouger/pkg/codec"
	"io"
	"os"
)

func main() {
	stdin, err := io.ReadAll(os.Stdin)

	if err != nil {
		panic(fmt.Errorf("unable to read data from stdin: %v", err))
	}

	if os.Args[1] == "decode" {
		codec.Decode(stdin)
	} else if os.Args[1] == "encode" {
		codec.Encode(stdin)
	} else {
		panic(fmt.Errorf("invalid argument: %v", os.Args[1]))
	}
}
