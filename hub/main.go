package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func serveHome(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Println(r.URL)
	http.ServeFile(w, r, "home.html")
}

func main() {
	router := httprouter.New()
	router.GET("/", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		serveHome(w, r, ps)
	})
	hub := newHub()
	go hub.run()
	router.GET("/chat/:room", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		serveWs(hub, w, r, ps)
	})
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
