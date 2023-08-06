package main

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
)

// a function type that represents an api endpoint.
// API endpoints are the entry points to our server or application.
// each one has a unique url, and serves a specific purpose
// when a client makes an http request, it does it to a particular URL/API endpoint
// the server processes the request and sends back a response with the data or performs the action necessary like DELETE for example.
// API endpoints are essential for data exchange between different components of a system
// Key purposes of endpoints are: Data retrieval(user information, product details, other data), Data modification(create, update, delete database records), Action execution(sending notifications, initiating a process, performing calculations), Authentication and authorization(user/password for access), integration(different systems and services to integrate with each other), versioning(can be versioned to support backwards compatibility), security(security mechanisms like access control, rate limiting, and data encryption to protect the server and it's resources), documentation(documented to provide certain developers with clear guidelines on how to use the API, including the available endpoints, their functionalities, and the expected data format)
type APIFunc func(context.Context, http.ResponseWriter, *http.Request) error

// a struct that represents the JSON response to be sent back to the client. has two fields.
type PriceResponse struct {
	Ticker string  `json:"ticker"`
	Price  float64 `json:"price"`
}

type JSONAPIServer struct {
	listenAddr string       //represents the address the server will listen on
	svc        PriceFetcher //a variable of type PriceFetcher, and the PriceFetcher type represents the service responsible for fetching price data
}

// a constructor function. creates and returns a new JSONAPIServer instance.
func NewJSONAPIServer(listenAddr string, svc PriceFetcher) *JSONAPIServer {
	return &JSONAPIServer{ //returns the JSONAPIServer struct instance with the arguments passed assigned to its fields. svc being of type PriceFetcher
		listenAddr: listenAddr,
		svc:        svc,
	}
}

func (s *JSONAPIServer) Run() {
	//this method starts the HTTP server and listens for incoming requests based on the *listenAddr string
	http.HandleFunc("/", makeHTTPHandlerFunc(s.handleFetchPrice)) //registers a request handler for the root path, which calls the handleFetchPrice
	http.ListenAndServe(s.listenAddr, nil)                        //*listens on listenAddr
}

// takes an APIFunc as input and returns an http.HandlerFunc
func makeHTTPHandlerFunc(apiFn APIFunc) http.HandlerFunc {
	ctx := context.Background()                                     //an original context must always be created when using context. only after can .WithValue() and other methods be used.
	ctx = context.WithValue(ctx, "requestID", rand.Intn(100000000)) //we are creating a context, ctx is the parent context, requestID is the key and the random int is the value
	//this following returned function has the above key requestID, which lets other functions that are involved in this process,
	//to log, trace, and perform other operations associated with this specific request.
	return func(w http.ResponseWriter, r *http.Request) {
		//we check if there is an error with the HTTP request, and doing this here is to centralize error handling
		//instead of handling errors within each individual API endpoint function. this way error handling code isn't duplicated in each individual function.
		if err := apiFn(context.Background(), w, r); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		}
	}
}

func (s *JSONAPIServer) handleFetchPrice(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	//we use the ticker as a sort of key to find our value.
	//a ticker is not exactly a key, but a unique symbol or code that represents a specific financial instrument
	//such as a stock or security on an exchange. it is used to IDENTIFY and trade individual assets in financial markets.
	ticker := r.URL.Query().Get("ticker")
	//we fetch the price
	price, err := s.svc.FetchPrice(ctx, ticker) //so our receiver s which is a pointer to JSONAPIServer has a PriceFetcher field as a part of its contract, is used to activate the PriceFetcher field and get the value based on the ticker we just put inside of the ticker reference

	if err != nil {
		//don't handle your errors inside of your request, your return them here and handle them later on in a single place
		return err
	}
	//we create a reference for a PriceResponse struct, and put the price and the ticker we just retrieved above into it.
	priceResp := PriceResponse{
		Price:  price,
		Ticker: ticker,
	}
	//we write the JSON and return it as a response to the http request.
	return writeJSON(w, http.StatusOK, &priceResp)
}

// a created function so we dont have to always duplicate/re-write the two lines.
func writeJSON(w http.ResponseWriter, s int, v any) error {
	w.WriteHeader(s)
	return json.NewEncoder(w).Encode(v)
}
