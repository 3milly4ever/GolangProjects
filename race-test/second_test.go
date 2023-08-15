package main

import (
	"sync/atomic"
)

type Server struct {
	//isRunning could potentially be accessed in a graceful shutdown
	//isRunning
	gameRound atomic.Int32
}

// func TestDataRaceConditions(t *testing.T) {
// 	n := 10
// 	s := Server{
// 		gameRound: atomic.LoadInt32(&n, 1),
// 	}
// }
