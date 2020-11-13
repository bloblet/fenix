# Fenix v6

[![Codeac](https://static.codeac.io/badges/2-281254941.svg "Codeac.io")](https://app.codeac.io/github/bloblet/fenix)
![Tests](https://github.com/bloblet/fenix/workflows/Tests/badge.svg)

Fenix is a social media platform inspired by [Discord](https://discord.com), but offers more security. 

DMs will have an end to end encryption option.

Different levels of account security can be enabled, yet to be decided on.

If you would like to help with Fenix, or make a Fenix client, [open an issue.](https://github.com/bloblet/fenix/issues).


# Workflow

- [Golang](https://golang.org/)
- [CassandraDB](https://cassandra.apache.org/)
- [gRPC](https://grpc.io/)
- [protobufs](https://developers.google.com/protocol-buffers)
- [Revive](https://github.com/mgechev/revive)

# Setting up on Linux
```bash
# Install the latest version of Go
sudo apt install golang-1.15

# Install the latest version of protoc
sudo apt install protobuf-compiler

# Install Fenix (This may take a while, there are large files in git history that we have to delete)
cd ~/
git clone https://github.com/bloblet/fenix

# Add it to GOROOT
ln -s fenix ../../usr/local/go/src

cd fenix

# Install dependencies
go get ./...

# Add gopath/bin to PATH in .bashrc
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.bashrc
source ~/.bashrc

# Install linter
go get -u github.com/mgechev/revive
```

# Setting up on windows
[TODO]
# Build Instructions
- Ensure you have latest golang installed and on your PATH variable.

#### To build the server:
- Run `go build -o fenix-server ./server`

#### To run the server:
- Run `go run server/main.go`

#### To build the client:
- Run `go build -o fenix-client ./client`

#### To run the client:
- Run `go run client/main.go`

# Fenix 6.0.1 API Documentation

**This document is a work in progress and is subject to change!**
