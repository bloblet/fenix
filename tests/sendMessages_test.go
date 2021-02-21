package tests

import (
	"testing"
	"time"

	"github.com/bloblet/fenix/client/client"
)

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
