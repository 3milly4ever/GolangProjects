FROM golang:1.20.5 

WORKDIR /gdocker-microservice

COPY go.mod /gdocker-microservice/
COPY go.sum /gdocker-microservice/

COPY . /gdocker-microservice/

RUN go build -o /pricefetcher

EXPOSE 3000

CMD [ "/pricefetcher"]