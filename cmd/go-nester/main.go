package main

import "github.com/sboy99/go-nester/pkg/cmd"

func main() {
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
