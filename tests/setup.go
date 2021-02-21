package tests

import (
	"fmt"
	"github.com/bloblet/fenix/server/api"
	"net"
	"testing"
)

func setupTestCase(t *testing.T) func(t *testing.T) {
	// Tells the server goroutine when to stop the server
	stop := make(chan bool)
	// Pauses the main thread until the api is running

	a := api.GRPCApi{}
	a.Prepare()
	lis, err := net.Listen("tcp", "localhost:4545")

	fmt.Println(t.Name())
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		finished := false
		go func() {
			a.Listen(lis)
			if !finished {
				t.Fail()
			}
		}()

		finished = <-stop
		a.S.Stop()
	}()

	return func(t *testing.T) {
		fmt.Printf(t.Name())
		stop <- true
	}
}
