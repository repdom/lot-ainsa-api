package config

import (
	"log"
	"os"
)

func SetupLogger() {

	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
