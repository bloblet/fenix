package main

import (
	api "github.com/bloblet/fenix/server/api"
)

func main() {
	a := api.GRPCApi{}
	a.Serve()
}
