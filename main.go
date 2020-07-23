package main

import databases "fenix/databases"

func main() {
	database := databases.UserDatabase{}
	print(database.UserExists("YAY"))
	// database.CreateFakeUser("YAY")
	print(database.UserExists("YAY"))
}
