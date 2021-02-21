package main

import (
	api "github.com/bloblet/fenix/server/api"
	"github.com/bloblet/fenix/server/utils"
)

func main() {
	utils.Log().Info("Bam")
	a := api.GRPCApi{}
	a.Serve()
}
