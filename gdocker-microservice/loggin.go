package main

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

type loggingService struct {
	//loggingService implements the PriceFetcher interface which allows it to get access to the price value and log it.
	//next is a way to wrap or decorate PriceFether and just indicates that this is a composition
	next PriceFetcher //implements the PriceFetcher interface which has the FetchPrice behavior that we need.
}

func NewLoggingService(next PriceFetcher) PriceFetcher {
	return &loggingService{
		next: next,
	}
}

// request ids will come through the context argument in this function

//this is a logger for FetchPrice, it uses FetchPrice to find the data with the ticker, but ultimately only to log them.

func (s *loggingService) FetchPrice(ctx context.Context, ticker string) (price float64, err error) {
	//time.Now captures the time at the exact time the defer statement is executed, as the begin time
	//the defer function calculates the duration of the operation by subtracting the begin time from the end time
	//happens after the return, because we need the return to get the data
	defer func(begin time.Time) {
		//the logrus package is used logging, and even provides log fields.
		logrus.WithFields(logrus.Fields{ //the logrus.Fields{} constructs a map-like structure that hold key-value pairs representing the fields to be logged
			//on your 'data dock' or 'elastic search' ? (google) you can see the search ID so you can trace the error and see what happened.
			"requestID": ctx.Value("requestID"),
			"took":      time.Since(begin), //we find out how long the return process took, with this time method.
			"err":       err,
			"price":     price,
		}).Info("fetchPrice") //info fetchprice also logs the message fetchPrice at the logInfo level
	}(time.Now()) //time.Now starts before the return s.next.FetchPrice to see how long it will take, then stores it in begin time.Time and logs the data.
	//the main purpose of the above code is to log information about the duration and result of the price fetching operation

	return s.next.FetchPrice(ctx, ticker)
}
