package tests

import (
	"fmt"
	"github.com/bloblet/fenix/client/client"
	"strconv"
	"testing"
	"time"
)

func TestRequestMessageHistory(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	c := client.Client{}
	c.Connect("test", "localhost:4545")
	message := make([]string, 60)

	for i := range message {
		message[i] = strconv.FormatInt(int64(i), 10)
		c.SendMessage(message[i])
	}

	fmt.Printf("%v", c.RequestMessageHistory(time.Now().Add(-time.Hour)))
}
