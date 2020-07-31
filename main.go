package main

import (
	"github.com/bejaneps/acroplia/cmd/cli"
	"github.com/rs/zerolog/log"
)

func main() {
	if err := cli.Execute(); err != nil {
		log.Fatal().Msgf("%v", err)
	}
}
