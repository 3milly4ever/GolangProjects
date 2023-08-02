package main

import (
	"cmd/api"
	"flag"
)

func main() {
	listenAddr := flag.String("listenaddr", ":9090", "the server address")
	flag.Parse()

	server := api.NewServer(*listenAddr)

}
