// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("static")))
	for _, name := range []string{"arduino", "java", "go", "scala"} {
		room := newRoom(name)
		http.HandleFunc("/chat/"+name, func(w http.ResponseWriter, r *http.Request) {
			serveWs(room, w, r)
		})
		go room.run()
	}
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	log.Printf("started chat server on port 8080")
}
