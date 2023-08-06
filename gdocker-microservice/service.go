package main

import (
	"context"
	"fmt"
	"time"
)

// context is also used for handling cancellations and timeouts.
// PriceFetcher is an interface that can fetch a price.
type PriceFetcher interface {
	FetchPrice(context.Context, string) (float64, error) //Context is a golang standard package, a very common practice to use this as the first argument of your functions because a lot of existing tools use context, and it's very adaptable and connectable

}

// implements the PriceFetcher interface
type priceFetcher struct{}

// this is business logic, just fetching a price. json representation is not needed here.
// just let it be clean and only handle the business logic
func (s *priceFetcher) FetchPrice(ctx context.Context, ticker string) (float64, error) {
	return MockPriceFetcher(ctx, ticker)
}

// a mock of where we will be getting our data from.
var priceMocks = map[string]float64{
	"BTC": 20_000.0,
	"ETH": 200.0,
	"GG":  100_000.0,
}

func MockPriceFetcher(ctx context.Context, ticker string) (float64, error) {
	//we mimic the HTTP round trip duration
	time.Sleep(100 * time.Millisecond)

	//ok is a boolean, and if the ticker exists it will return true and price will be set to the value of the value the ticker matches in the map
	price, ok := priceMocks[ticker]
	//error check
	if !ok {
		//if the ticker doesn't match any in our map.
		return price, fmt.Errorf("The given ticker (%s) is not supported", ticker)
	}
	//return the price value and nil for error
	return price, nil
}
