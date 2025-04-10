package main

import (
	"banco/database"
	"banco/server"
)

func main() {
	database.ConnectToDatabase()
	server.RunServer()
}