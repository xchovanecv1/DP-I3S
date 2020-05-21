# Base build image
FROM golang:1.11-alpine AS build
 
# Install some dependencies needed to build the project
RUN apk add bash ca-certificates git gcc g++ libc-dev
WORKDIR /go/src/go-base
 
# Force the go compiler to use modules
ENV GO111MODULE=on
 
# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .
 
#This is the ‘magic’ step that will download all the dependencies that are specified in 
# the go.mod and go.sum file.
# Because of how the layer caching system works in Docker, the  go mod download 
# command will _ only_ be re-run when the go.mod or go.sum file change 
# (or when we add another docker instruction this line)

RUN go mod download
 
# This image builds the weavaite server
FROM build AS builder
# Here we copy the rest of the source code
COPY . .
# And compile the project
RUN env GOOS=linux GOARCH=amd64 go build -o main-amd64 .

RUN env GOOS=linux GOARCH=arm GOARM=7 go build -o main-armv7 .

RUN file="$(ls)" && echo $file
#In this last stage, we start from a fresh Alpine image, to reduce the image size and not ship the Go compiler in our production artifacts.

FROM alpine:latest AS amd64

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /go/src/go-base/main-amd64 .
#COPY --from=builder /go/src/go-base/config.yaml .
#COPY --from=builder /go/src/go-base/templates ./templates
#COPY --from=builder /go/src/go-base/wait.sh ./wait.sh

RUN ["chmod", "+x", "./main-amd64"]
#RUN ["chmod", "777", "./templates"]

EXPOSE 4040

ENTRYPOINT ["./main-amd64"]

#CMD [ "serve", "--config", "config.yaml" ]

FROM arm32v7/alpine:latest AS armv7

#RUN apk add --no-cache ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /go/src/go-base/main-armv7 .
#COPY --from=builder /go/src/go-base/config.yaml .
#COPY --from=builder /go/src/go-base/templates ./templates
#COPY --from=builder /go/src/go-base/wait.sh ./wait.sh

#RUN chmod +x wait.sh
#RUN ["chmod", "777", "./templates"]

EXPOSE 4040

ENTRYPOINT ["./main-armv7"]

#CMD [ "serve", "--config", "config.yaml" ]
