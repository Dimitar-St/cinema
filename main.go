package main

import (
	"cinema/db"
	"cinema/localserver"
)

func main() {
	println("It compiles!")
	db.InitDB()
	println("It compiles!")

	srv := localserver.Server{}
	srv.Start()

}
