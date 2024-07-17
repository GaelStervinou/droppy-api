# Start from golang base image
FROM golang:alpine

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git && apk add --no-cache bash && apk add build-base

# Setup folders
RUN mkdir /app
WORKDIR /app

# Copy the source from the current directory to the working Directory inside the container
COPY . .
COPY .env .

# Download all the dependencies
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

#Setup hot-reload for dev stage
RUN go install -mod=mod github.com/githubnemo/CompileDaemon
RUN go install golang.org/x/tools/gopls@latest

ENTRYPOINT CompileDaemon --build="go build -a -installsuffix cgo -o main ." --command=./main