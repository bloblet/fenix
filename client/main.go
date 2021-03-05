package main

import (
	"bufio"
	"flag"
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
	addr := flag.String("addr", "localhost:4545", "Address of fenix to connect to")
	flag.Parse()

	fmt.Printf("Connecting to %v\n", *addr)

	c := client.Client{}
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Pick a username: ")
	username, _ := reader.ReadString('\n')

	c.Connect(sanitize(username), *addr)

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
