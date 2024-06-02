package main

import (
	"github.com/minacio00/easyCourt/config"
	"github.com/minacio00/easyCourt/server"
)

func main() {
	config.LdVars()
	server.StartServer()
}
