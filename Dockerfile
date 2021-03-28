FROM golang:1.16

# RUN apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /app/qr-code-generator

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Build the Go app
RUN go build -o ./out/qr-code-generator .


# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the binary program produced by `go install`
CMD ["./out/qr-code-generator"]