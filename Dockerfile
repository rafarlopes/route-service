FROM golang:1.15 as build

WORKDIR /go/src/github.com/rafarlopes/route-service
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
RUN make setup
RUN make

FROM debian

COPY --from=build /go/src/github.com/rafarlopes/route-service/dist /app
EXPOSE 8080

ENTRYPOINT [ "/app/route-service" ]

