package main

import (
	"github.com/sboy99/go-nester/pkg/cmd"
	"github.com/sboy99/go-nester/pkg/config"
)

func main() {
	if _, err := config.LoadConfig(); err != nil {
		panic(err)
	}
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
