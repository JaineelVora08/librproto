package router

import (
	"github.com/JaineelVora08/librproto/controller"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/fetch/{ts}", controller.GetTimeMessages).Methods("GET")
	r.HandleFunc("/submit", controller.AddNewMessage).Methods("POST")
	r.HandleFunc("/fetchall", controller.GetAllAcceptedMessages).Methods("GET")

	return r
}
