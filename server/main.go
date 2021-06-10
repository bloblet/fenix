package main

import (
	"github.com/bloblet/fenix/server/api"
)

func main() {
	api := api.GRPCApi{}
	api.Serve()
}
