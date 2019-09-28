package core

import "net/http"

type Controller interface {
	Bean

	RegisterRoutes() map[string]func(writer http.ResponseWriter, request *http.Request)

	HandleRoutes(writer http.ResponseWriter, request *http.Request) (func(writer http.ResponseWriter, request *http.Request), bool)
}
