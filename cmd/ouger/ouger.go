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
		bytes, err := codec.Decode(stdin)
		if err != nil {
			panic(fmt.Errorf("unable to decode data: %v", err))
		}
		os.Stdout.Write(bytes)
	} else if os.Args[1] == "encode" {

		bytes, err := codec.Encode(stdin)
		if err != nil {
			panic(fmt.Errorf("unable to encode data: %v", err))
		}
		os.Stdout.Write(bytes)
	} else {
		panic(fmt.Errorf("invalid argument: %v", os.Args[1]))
	}
}
