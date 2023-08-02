package main

import (
	"flag"
	"fmt"
	"log"
	"structure/api"
	"structure/store"
)

func main() {
	listenAddr := flag.String("listenaddr", ":9090", "the server address")
	flag.Parse()

	store := store.NewMemoryStore()
	server := api.NewServer(*listenAddr, store)
	fmt.Println("server running on port: ", *listenAddr)
	log.Fatal(server.Start())

}
