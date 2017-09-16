package main

import (
	"net/http"

	"github.com/sbiscigl/load-balancer/params"
	"github.com/sbiscigl/load-balancer/requesthandler"
	"github.com/sbiscigl/load-balancer/server"
)

func main() {
	env := params.New()
	health := server.NewServerHealthMap(env)
	balancer := requesthandler.New(health)

	health.PrintMap()

	http.Handle("/", balancer)
	http.ListenAndServe(":"+env.GetPort(), nil)
}
