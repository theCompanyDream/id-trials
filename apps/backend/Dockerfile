FROM golang:1.25.0-alpine

# Install Air (requires git and build-base)
RUN apk add --no-cache git build-base

RUN go install github.com/air-verse/air@latest

# Creates an app directory to hold your app’s source code
WORKDIR /opt/app

# Copies everything from your root directory into /opt/app
COPY . .

# Installs Go dependencies
RUN go mod download

# Specifies the executable command that runs when the container starts
ENTRYPOINT ["air"]
