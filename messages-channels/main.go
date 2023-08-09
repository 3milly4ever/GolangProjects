package main

import (
	"fmt"
	"time"
)

// message communication with for select(daemons)
type Message struct {
	From    string
	Payload string //normally a slice of bytes
}

type Server struct {
	msgch  chan Message  //declare channel names with ch at the end so everyone knows its a channel
	quitch chan struct{} //a quit channel, which quits the server. empty struct because it will be 0 bytes and take less than a bool which is 2 bytes and an int which is one
}

// a server is not always something that listens to a port, it serves
// server has a channel, and it's listening and waiting
func (s *Server) StartAndListen() {
	//this is the main logic of this code that is used for communication between servers, structs and daemons.
	//you can name your for loop
free:
	for {
		//we can make a select statement to prevent deadlocks
		//a select statement is like a switch, but a concurrent switch
		select {
		//block here until someone is sending a message to the channel
		//if someone is sending a message, execute this, but if no one is sending a message we need a default function
		case msg := <-s.msgch: //receiving message through message channel
			fmt.Printf("received message from: %s payload %s\n", msg.From, msg.Payload)

		//if we get a signal to quit we execute this logic
		case <-s.quitch:
			fmt.Println("The server is doing graceful shutdown")
			//logic for the graceful shutdown
			break free
		//if none of the other cases are ready to communicate, this default case will be executed so then the loop will continue without waiting
		//if the program should be doing something even while the channels arent ready to communicate such as checking for new messages or events,
		//the default case can allow you to handle such situations without getting stuck
		default:

		}
	}
	fmt.Println("the server is shut down")

}

// we send to the server's channel
func sendMessageToServer(msgch chan Message, payload string) {
	msg := Message{
		From:    "Sophie",
		Payload: payload,
	}
	//msg is what we created above, sends the message to created go channel
	msgch <- msg

	fmt.Println("Sending message")
}

func gracefulQuitServer(quitch chan struct{}) {
	close(quitch)
}

func main() {
	s := &Server{
		msgch:  make(chan Message),
		quitch: make(chan struct{}),
	}

	go s.StartAndListen()

	done := make(chan struct{})

	go func() {
		time.Sleep(2 * time.Second)
		sendMessageToServer(s.msgch, "Hello!")
		done <- struct{}{}
	}()

	go func() {
		time.Sleep(4 * time.Second)
		gracefulQuitServer(s.quitch)
		done <- struct{}{}
	}()

	go sendMessageToServer(s.msgch, "Love you")
	//ensures that the main function doesn't exit and keeps the goroutines running.
	<-done
	<-done
}
