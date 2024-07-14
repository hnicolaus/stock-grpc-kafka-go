# Stock

## Prerequisites
Make sure you have installed all of the following prerequisites on your development machine:
* go version : [1.20](https://golang.org/dl/)
* Linux operating system

## Code structure
This service attempts to implement the Clean Architecture concept by separating the software into layers. [link](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html).

The project is mainly strucutred as follows:
- server
- handler
- usecase
- repo

The contracts and entities being used by the project are organized as follows:
- proto
- model


## Local Development
Build and start the apps:
- `go build .` to build a binary
- `./challenge` to start
or simply build and start the app directly:
- `go run .` to build a binary then start the app

### Test and Lint
golangci-lint run
gotest -v --race ./...

For manual testing the GRPC server in local environment:
- you can use any GUI client for gRPC services, some recommendations are gRPCox [ref](https://github.com/gusaul/grpcox#installation) or BloomRPC [ref](https://github.com/bloomrpc/bloomrpc)
- please use `localhost:50051` or `0.0.0.0:50051` as the target gRPC Server.
- stock.proto file is provided in the root directory of this project
