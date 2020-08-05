# Golang Chatrooms Example

This repository contains multiple implementations of multiple realtime chatrooms in golang. Both directories were adapted from [the gorilla websocket chat example](https://github.com/gorilla/websocket/tree/master/examples/chat) to handle JSON messages and multiple chatrooms.

In the "hub" directory, clients register their websocket connection into the hub. The hub maintains a map of rooms, and rooms maintain a map of clients. When a client broadcasts a message to the hub, the message contains the name of the room they are registered in. The hub loops through the room's clients and sends them the message. The hub also handles closing (removing from the `rooms` map) rooms that have no connected clients.

In the "independent" directory, each individual room serve's a websocket connection through an http.HandleFunc. Clients register, send messages, and unregister throught independent rooms. Rooms broadcast messages to all connected clients. Currently, deleting stale or empty rooms is not being handled. This is more difficult than in the "hub" directory, as there is nothing keeping track of all the registered rooms and their clients.

Both of the examples are valid implementation of multiple chatrooms in golang. However, one may be more suited to your particular usecase.

