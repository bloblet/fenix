package tests

import (
	"fmt"
	"net"
	"strconv"
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

func TestSendMessage(t *testing.T) {
	teardownTestCase := setupTestCase(t)

	c := client.Client{}
	c.Connect("test", "localhost:4545")

	content := "look at this perfect message!"
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
	teardownTestCase(t)
}

func TestRequestMessageHistory(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	c := client.Client{}
	c.Connect("test", "localhost:4545")
	message := make([]string, 60)

	for i, _ := range message {
		message[i] = strconv.FormatInt(int64(i), 10)
		c.SendMessage(message[i])
	}

	fmt.Printf("%v", c.RequestMessageHistory(time.Now().Add(-time.Hour)))
}