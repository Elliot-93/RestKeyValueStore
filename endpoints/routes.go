package endpoints

import (
	"net/http"
)

type RouteDetails struct {
	Handler      http.Handler
	AuthRequired bool
}

var Routes = map[string]RouteDetails{
	PingRoute:     {Handler: PingHandler{}, AuthRequired: false},
	StoreRoute:    {Handler: StoreHandler{}, AuthRequired: true},
	ListRoute:     {Handler: ListHandler{}, AuthRequired: false},
	ShutdownRoute: {Handler: ShutdownHandler{}, AuthRequired: true},
	LoginRoute:    {Handler: LoginHandler{}, AuthRequired: false},
}
