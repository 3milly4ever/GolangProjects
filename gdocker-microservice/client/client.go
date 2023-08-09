package client

//the client is responsible for making requests to a remote server's endpoint to fetch price data. it holds the information of how to interact
//with the server and returns the fetched data as a types.PriceResponse.

import (
	"context"
	"encoding/json"
	"fmt"
	"gdocker-microservice/types"
	"net/http"
)

// this struct will represent a client. a client communicates with a remote server(to fetch price data in our case)
// the struct encapsulates the details of how to interact with the server
type Client struct {
	//has a single field, endpoint which holds the URL or address of the remote server's endpoint.
	endpoint string
}

// this is a client constructor, takes in the endpoint as an argument and returns a pointer to the newly created Client
func New(endpoint string) *Client {
	return &Client{
		endpoint: endpoint,
	}
}

// this method is implemented to fetch price data from the remote server using an HTTP get request.
func (c *Client) FetchPrice(ctx context.Context, ticker string) (*types.PriceResponse, error) {
	endpoint := fmt.Sprintf("%s?ticker=%s", c.endpoint, ticker) //https://api.example.com/prices?=BTC" everything before ? is the endpoint and the ticker is BTC, in this example

	//we create a new request, which is a get at the endpoint address provided.
	req, err := http.NewRequest("get", endpoint, nil)
	//error check the creation of the new request
	if err != nil {
		return nil, err
	}
	//now that we have the request object, we use it here to make the request and the response is put into the resp reference below.
	resp, err := http.DefaultClient.Do(req)
	//error check the action of making the request.
	if err != nil {
		return nil, err
	}
	//for any error that doesn't equal an http.StatusOK will get this error
	if resp.StatusCode != http.StatusOK {

		httpErr := map[string]any{}
		if err := json.NewDecoder(resp.Body).Decode(&httpErr); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("error service responded with non OK status code: %s", httpErr["error"])
	}
	//we expect the server to response with JSON data matching types.PriceResponse structure(with ticker and price)
	priceResp := new(types.PriceResponse)
	//we decode the json response using the NewDecoder and we insert it into a PriceResponse structure that will hold the fetched json price data.
	if err := json.NewDecoder(resp.Body).Decode(priceResp); err != nil {
		return nil, err
	}
	//if the request is successful and the response is properly decoded we return priceResp which contains the fetched price data in a PriceResponse structure,
	//and a nil error indicating a successful fetch
	return priceResp, nil
}
