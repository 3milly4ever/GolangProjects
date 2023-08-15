FROM golang:1.20.5 
#this directory will be used as context for the following commands
WORKDIR /gdocker-microservice
#copies the go.mod and go.sum files from our local file system into the /gdocker-microservice/ directory inside our container
COPY go.mod /gdocker-microservice/
COPY go.sum /gdocker-microservice/
#copies the entire content of our current directory, where this Docker file is located, into the gdocker-microservice directory inside the container
#includes our application source code and any additional files needed for microservices.
COPY . /gdocker-microservice/
#an executable binary named pricefetcher is created and is saved in the root directory of the container 
RUN go build -o /pricefetcher
#indicates that the microservice will be listening for incoming connections on port 3000
EXPOSE 3000
#specifies the default command to run when the container starts, runs the pricefetcher executable
CMD [ "/pricefetcher"]