
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
