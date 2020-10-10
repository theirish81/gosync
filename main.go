package main

import (
	"flag"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)
	mode := flag.String("mode", "client", "Mode (server|get|*client)")
	flag.Parse()
	config := LoadConfig()
	switch *mode {
	case "server":
		StartServer(config)
	case "client":
		StartClient(config)
	}
}
