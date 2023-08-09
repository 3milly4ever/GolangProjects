package main

import (
	"flag"
)

// we will fetch prices of crypto currencies
func main() {
	//we have command-line flags here so the user can customize the behavior of the microservice without having to modify the code itself.
	//here the user can specify the listen address to use allowing them to change it to a different port if necessary.
	//the microservice becomes more flexible by using command line flags in the sense that it can be deployed on different machines or environments
	//with different configurations. makes it easier for developers and system administrators to interact with the microservice without having to
	//modify the source code.

	//we are going to access our service in a gateway not a browser.
	//we are going to containerize this.

	//kind of simulates docker
	// client := client.New("http://localhost:3000")
	// price, err := client.FetchPrice(context.Background(), "BT")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Printf("%+v\n", price)
	// return

	//we define a command line flag with flag.String that allows the user to specify the address at which the service will listen
	listenAddr := flag.String("listenaddr", ":3000", "listen address the service is running")
	flag.Parse()                                                //we parse the provided flags so the program can read and interpret the values above from the command line.
	svc := NewLoggingService(NewMetricService(&priceFetcher{})) //we have our original priceFetcher functionality, which then gets the metric service
	//functionality and the logging service functionality wrapped around it.
	//we add the functionalities separately so they can be  edited and scaled easily if additions or changes are necessary to make.
	//the benefits of this are:

	//modularity: each component becomes self-contained and can be developed, tested, and maintained independently. this promotes reusability
	//and makes it easier to understand and manage the codebase

	//separation of concerns: with each functionality separated into its own module, developers can focus on specific aspects of the system
	//without worrying about unrelated parts, and also leads to cleaner and more maintanable code

	//extensibility:new functionalities can be added like first we wrapped the loggingservice then the newmetric service after.
	//makes it convenient to scale and add new features to the system

	//testability:the abstraction and modularity makes it easier to write unit tests. mock implementations can be created for testing each
	//component in isolation. ensuring the individual functionalities work correctly and are thoroughly tested.

	//flexibility: we can easily modify specific functionalities without affecting other parts of the system, which allows for easy customization
	//and adaptation if things were to change.

	//maintainability: code complexity is reduced and systems are easier to maintain over time after a lot of changes are made. with changes
	//or updates developers can focus on specific components without worrying about unintended side effects on other parts of the codebase.

	//we create a new JSON API server.
	server := NewJSONAPIServer(*listenAddr, svc)
	//we run the server, and the run method listens for incoming HTTP requests and handles them based on the defined API endpoints.
	server.Run()
	//fmt.Println(price)

}
