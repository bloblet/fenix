package main

import (
	databases "fenix/databases"
	"fmt"
)

func main() {
	database := databases.UserDatabase{}
	fmt.Print(database.UserExists("YAY"))
	// database.CreateFakeUser("YAY")
	fmt.Print(database.UserExists("YAY"))
}
