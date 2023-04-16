package main

import (
	"github.com/stonear/api_gateway/db"
	"github.com/stonear/api_gateway/log"
	"github.com/stonear/api_gateway/repository"
	"github.com/stonear/api_gateway/route"
	"github.com/stonear/api_gateway/server"
)

func main() {
	// Load logger
	l := log.New()

	// Load database
	d, err := db.New()
	if err != nil {
		l.Errorln(err)
		return
	}
	d.Migrate()

	// Load route list
	r := repository.NewRouteRepository(d)
	routeList, err := r.Get()
	if err != nil {
		l.Errorln(err)
		return
	}

	// Parse routes
	router := route.New(l, routeList)

	// Run server
	s := server.New(l, "3000", router)
	s.RunGracefully()
}
