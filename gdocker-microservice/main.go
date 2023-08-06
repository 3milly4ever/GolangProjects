package main

import (
	"flag"
)

// we will fetch prices of crypto currencies
func main() {
	listenAddr := flag.String("listenAddr", ":3000", "listen address the service is running")
	flag.Parse()
	svc := NewLoggingService(NewMetricService(&priceFetcher{}))

	server := NewJSONAPIServer(*listenAddr, svc)
	server.Run()
	//fmt.Println(price)

}
