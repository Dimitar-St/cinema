package main

import (
	"cinema/db"
	"cinema/localserver"
)

func main() {
	db.InitDB()

	srv := localserver.Server{}
	srv.Start()

	println("It compiles!")
}
