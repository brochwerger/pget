
# Build minimal image based on https://github.com/mr-pascal/medium-go-docker/blob/master/Dockerfile

FROM golang:1.18 as builder
RUN mkdir /app
WORKDIR /app
COPY go.mod .

### Setting a proxy for downloading modules
ENV GOPROXY https://proxy.golang.org,direct

### Download Go application module dependencies
RUN go mod download

### Copy actual source code for building the application
COPY . .

### CGO has to be disabled cross platform builds
### Otherwise the application won't be able to start
ENV CGO_ENABLED=0

RUN GOOS=linux go build ./pget.go

FROM alpine

COPY --from=builder /app/pget /usr/bin/pget

RUN apk update && \
    apk add bash bind-tools curl

CMD tail -f /dev/null
