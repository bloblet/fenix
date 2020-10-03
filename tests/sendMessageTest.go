package tests

import (
	"testing"
	"time"

	"github.com/bloblet/fenix/client/client"
	"github.com/bloblet/fenix/server/api"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func setupTestCase(t *testing.T) (string, func(t *testing.T)) {
	// Tells the server goroutine when to stop the server
	stop := make(chan bool)
	// Pauses the main thread until the api is running

	a := api.GRPCApi{}
	a.Prepare()
	lis = bufconn.Listen(bufSize)
	go func() {
		go a.Listen(lis)
		<-stop
		a.S.Stop()
	}()

	net, err := lis.Dial()

	if err != nil {
		t.Fatal(err)
	}

	return net.LocalAddr().String(), func(t *testing.T) {
		stop <- true
	}
}

func TestSendMessage(t *testing.T) {
	addr, teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	c := client.Client{}
	c.Connect("test", addr)
	
	content := "look at this perfect messsage!"
	readMessage := make(chan bool)

	go func() {
		msg := <-c.Messages 
		readMessage <- true
		
		if msg.ID != c.Username || msg.Content != content  {
			t.Fail()
		}
	}()
	
	c.SendMessage(content)
	select {
	case <-readMessage:
	case <-time.After(2 * time.Second):
		t.Fail()
	}
}
