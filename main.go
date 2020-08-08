package main

import (
	api "fenix/api"
	"flag"
	"fmt"
	"os"
)

func main() {
	fl := flag.NewFlagSet("Fenix", flag.ExitOnError)
	user := fl.String("user", "root", "User to connect to etcd with")
	password := fl.String("password", "", "User's password")

	
	fl.Parse(os.Args[1:])

	if *password == "" {
		fmt.Print("Help for Fenix:")
		fl.Usage()
		os.Exit(1)
	}
	
	a := api.NewAPI(*user, *password)
	err := make(chan api.Error)

	a.Serve(err, false)
}
