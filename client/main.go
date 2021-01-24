package main

import (
	"bufio"
	"fmt"
	"github.com/bloblet/fenix/client/client"
	"github.com/pborman/ansi"
	"os"
	"strings"
	"time"
)

func sanitize(dirty string) string {
	return strings.ReplaceAll(strings.ReplaceAll(dirty, "\n", ""), "\r", "")
}

func main() {
	c := client.Client{}
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Pick a username: ")
	username, _ := reader.ReadString('\n')

	c.Connect(sanitize(username), "vps.bloblet.com:4000")

	go func() {
		for true {
			msg := <-c.Messages
			fmt.Printf("<%v> %v\n", msg.GetUserID(), msg.GetContent())
		}
	}()

	fmt.Println("Connected to Fenix")
	for true {
		text, _ := reader.ReadString('\n')
		fmt.Printf("%v%v", ansi.CPL, ansi.DL)

		if strings.HasPrefix(text, "/before") {
			messages := c.RequestMessageHistory(time.Now())
			for _, message := range messages {
				fmt.Printf(">>> %v <%v> %v\n", message.SentAt.AsTime().String(), message.UserID, message.Content)
			}
		} else {
			c.SendMessage(sanitize(text))
		}
	}
}
