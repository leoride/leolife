package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/leoride/leolife/handler"
	"github.com/leoride/leolife/util/rest"
	"log"
	"net/http"
)

var restRouter *mux.Router

func SetupRestConfig() {
	restRouter = mux.NewRouter().StrictSlash(true)

	setupPersonHandler(restRouter)
	setupAccountHandler(restRouter)

	log.Fatal(http.ListenAndServe(fmt.Sprint(":", "8080"), restRouter))
}

func setupPersonHandler(r *mux.Router) {
	personHandler := handler.NewPersonHandler(sessionFactory)

	r.Path("/api/person/{id}").Methods("GET").HandlerFunc(rest.RestAuthHandlerWrapper(personHandler.HandleGetOne))
	r.Path("/api/person/{id}").Methods("PUT").HandlerFunc(rest.RestAuthHandlerWrapper(personHandler.HandleUpdateOne))
	r.Path("/api/person/{id}").Methods("DELETE").HandlerFunc(rest.RestAuthHandlerWrapper(personHandler.HandleDeleteOne))
	r.Path("/api/person").Methods("GET").HandlerFunc(rest.RestAuthHandlerWrapper(personHandler.HandleGetAll))
	r.Path("/api/person").Methods("POST").HandlerFunc(rest.RestAuthHandlerWrapper(personHandler.HandleCreate))
}

func setupAccountHandler(r *mux.Router) {
	accountHandler := handler.NewAccountHandler(sessionFactory)

	r.Path("/api/account").Methods("POST").HandlerFunc(rest.RestHandlerWrapper(accountHandler.HandleCreate))
	r.Path("/api/account/auth").Methods("POST").HandlerFunc(rest.RestHandlerWrapper(accountHandler.HandleAuthentication))
}
