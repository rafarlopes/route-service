FROM golang:1.15 as build

WORKDIR /app
COPY . .

RUN make setup verify sec

ENV CGO_ENABLED 0
RUN make bin

FROM scratch

COPY --from=build /app/dist /
EXPOSE 8080

ENTRYPOINT [ "/route-service" ]