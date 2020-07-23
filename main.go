package main

import (
	api "fenix/api"
)

func main() {
	api := api.API{}
	api.Serve()
}
