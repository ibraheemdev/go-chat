package main

import (
	"log"
	"net/http"
)

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func main() {
	for _, name := range []string{"arduino", "java", "go", "scala"} {
		room := newRoom(name)
		http.HandleFunc("/chat/"+name, func(w http.ResponseWriter, r *http.Request) {
			serveWs(room, w, r)
		})
		go room.run()
	}
	http.HandleFunc("/", serveHome)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
