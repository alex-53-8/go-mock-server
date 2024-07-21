# build
FROM golang:1.22.5-alpine3.19 as build

WORKDIR /build

COPY go.mod /build/go.mod
COPY go.sum /build/go.sum
COPY main.go /build/main.go
COPY server /build/server

RUN go build -o mock-server /build

# image
FROM golang:1.22.5-alpine3.19

WORKDIR /app

COPY --from=build /build/mock-server /app/mock-server
RUN chmod 775 /app/mock-server

EXPOSE 8081

CMD [ "/app/mock-server", "-file", "model/server.yml" ]