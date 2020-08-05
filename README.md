# Golang Chatrooms Example

This repository contains multiple implementations of multiple realtime chatrooms in golang.

In the "independent" directory, each individual room serve's a websocket connection through an http.HandleFunc. Clients connect to, register, send messages, and unregister throught independent rooms. Rooms broadcast messages to all connected clients
