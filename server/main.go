package main

import (
	api "github.com/bloblet/fenix/server/api"
	// "flag"
	// "fmt"
	// "os"
)

func main() {
	// fl := flag.NewFlagSet("Fenix", flag.ExitOnError)
	// password := fl.String("password", "", "User's password")

	// fl.Parse(os.Args[1:])

	// if *password == "" {
	// 	fmt.Print("Help for Fenix:")
	// 	fl.Usage()
	// 	os.Exit(1)
	// }

	// a := api.NewAPI(*user, *password, "")
	// a.Serve(false)

	a := api.NewGRPCApi()
	a.Serve()
}
