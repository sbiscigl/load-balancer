package main

import (
	"log"
	"net/http"

	"github.com/sbiscigl/load-balancer/params"
	"github.com/sbiscigl/load-balancer/requesthandler"
	"github.com/sbiscigl/load-balancer/server"
)

func main() {
	env := params.New()
	health := server.NewServerHealthMap(env)
	balancer := requesthandler.New(health)

	http.Handle("/", balancer)
	err := http.ListenAndServe(":"+env.GetPort(), nil)
	if err != nil {
		log.Println("cannot start server on that port")
	}
}
