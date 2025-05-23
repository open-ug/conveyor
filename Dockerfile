# Building the binary of the App
FROM golang:1.23.0 AS build

# `boilerplate` should be replaced with your project name
WORKDIR /go/src/boilerplate

# Copy all the Code and stuff to compile everything
COPY . .

# Downloads all the dependencies in advance (could be left out, but it's more clear this way)
RUN go mod tidy

# Builds the application as a staticly linked one, to allow it to run on alpine
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o app .


# Moving the binary to the 'final Image' to make it smaller
FROM ubuntu:latest AS release

WORKDIR /app


# `boilerplate` should be replaced here as well
COPY --from=build /go/src/boilerplate/app .

RUN 

# Add packages
RUN apt-get update && \
  apt-get install -y ca-certificates


RUN ls


# Expose the port the app runs on
EXPOSE 3000

ENTRYPOINT [ "/app/app" ]

CMD ["api-server"]