package main

import (
	"github.com/minacio00/easyCourt/config"
	"github.com/minacio00/easyCourt/database"
	"github.com/minacio00/easyCourt/server"
)

func main() {
	config.LdVars()
	database.ConnectDb()
	server.StartServer()
}
