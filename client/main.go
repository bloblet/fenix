package main

import (
	"bufio"
	"fmt"
	"github.com/bloblet/fenix/client/client"
	"github.com/pborman/ansi"
	"os"
	"strings"
)

func sanitize(dirty string) string {
	return strings.ReplaceAll(strings.ReplaceAll(dirty, "\n", ""), "\r", "")
}

func main() {
	c := client.Client{}
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Pick a username: ")
	username, _ := reader.ReadString('\n')

	c.Connect(sanitize(username), "bloblet.com:4000")

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
		c.SendMessage(sanitize(text))
	}
}
