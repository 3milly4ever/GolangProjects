package main

//the purpose of this metrics file, and what metrics are used for is to monitor the performance of services, track user behavior,
//identify bottlenecks and detect anomalies. the data found by metrics are used for debugging, performance optimization, capacity planning,
//and overall system analysis.

import (
	"context"
	"fmt"
)

// a metricService here is implemented as a wrapper around an existing PriceFetcher service.
// this metricService adds the functionality of pushing metrics to a metrics storage(one of which is Prometheus), and then delegates the actual
// fetching of the price to the original PriceFetcher. we are adding a metric collection functionality to the FetchPrice method.
// by adding this functionality, we don't modify the the core functionality of the original service(PriceFetcher) and this helps to keep
// the price fetching and the metric collection separate, which promotes a modular and extensible design.
type metricService struct {
	next PriceFetcher
}

// this is a constructor that creates a new metricService. takes in an existing PriceFetcher as an argument, and returns a new PriceFetcher
// with metric collection capabilities.
func NewMetricService(next PriceFetcher) PriceFetcher {
	//this returned PriceFetcher will be a pointer to the metricService struct.
	return &metricService{
		next: next,
	}
}

// the purpose of this metricService is to collect metrics(performance of the services, user behavior)
// metrics can include measurements of the time taken for the operation, the number of requests made, and any errors encountered.
// the storage of metrics in Prometheus can involve additional logic and configurations specific to the metrics storage being used.
func (s *metricService) FetchPrice(ctx context.Context, ticker string) (price float64, err error) {
	//we push the metrics to Prometheus
	fmt.Println("pushing metrics to prometheus")
	//after we push the metrics we delegate the price fetching to the original PriceFetcher by calling it through the s.next field inside the metricService
	return s.next.FetchPrice(ctx, ticker)
}
