# go-qr
repository to create your own qr code in Go

## Running with source code
**pre reqs:**
- Go installed
- Git installed

1) Fork and then Clone the repository into $GOPATH/src/github.com/
2) `cd` into the cloned repository
3) Run the command: `go run main.go` OR `go build -o qr-generator` and then `./qr-generator`
4) In a browser navigate to `localhost:8080`

## Running with Docker
**pre reqs:**
- Docker account and Docker desktop client installed

1) Fork and then Clone the repository onto your machine
2) `cd` into the cloned repository
3) Run the command: `docker build -t <your-docker-id>/qr-code-generator .`
4) Run the command: `docker run -dp 8080:8080 <your-docker-id>/qr-code-generator`
5) In a browser navigate to `localhost:8080`