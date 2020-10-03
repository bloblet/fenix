package main

import (
	"bufio"
	"os"
	"strings"
	"fmt"
	"github.com/pborman/ansi"
)



func main() {
	c := Client{}
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Pick a username: ")
	username, _ := reader.ReadString('\n')

	c.Connect(strings.ReplaceAll(username, "\n", ""))

	go func(){
		for true {
			msg := <-c.Messages
			fmt.Printf("<%v> %v\n", msg.GetUserID(), msg.GetContent())
		}
	}()
	
	fmt.Println("Connected to Fenix")
	for true {
		text, _ := reader.ReadString('\n')
		fmt.Printf("%v%v", ansi.CPL, ansi.DL)
		c.SendMessage(text)
	}
}
