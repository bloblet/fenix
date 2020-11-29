package tests

import (
	"net"
	"testing"
	"time"

	// "time"

	"github.com/bloblet/fenix/client/client"
	"github.com/bloblet/fenix/server/api"
)

func setupTestCase(t *testing.T) func(t *testing.T) {
	// Tells the server goroutine when to stop the server
	stop := make(chan bool)
	// Pauses the main thread until the api is running

	a := api.GRPCApi{}
	a.Prepare()
	lis, err := net.Listen("tcp", "localhost:4545")

	if err != nil {
		t.Fatal(err)
	}

	go func() {
		go a.Listen(lis)
		<-stop
		a.S.Stop()
	}()

	return func(t *testing.T) {
		stop <- true
	}
}

func TestSendMessage(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	c := client.Client{}
	c.Connect("test", "localhost:4545")

	content := "look at this perfect messsage!"
	readMessage := make(chan bool)
	done := make(chan bool)
	go func() {
		msg := <-c.Messages
		readMessage <- true

		if msg.UserID != c.Username || msg.Content != content {
			t.Error("Unexpected message")
		}
		done <- true
	}()

	c.SendMessage(content)
	select {
	case <-readMessage:
	case <-time.After(2 * time.Second):
		t.Error("Timeout")
	}
	<-done
}
