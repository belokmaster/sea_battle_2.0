package main

import "net/http"

func newRouter() *http.ServeMux {
	mux := http.NewServeMux()

	apiMux := http.NewServeMux()
	apiMux.HandleFunc("/game", gameStatusHandler)
	apiMux.HandleFunc("/newgame", newGameHandler)
	apiMux.HandleFunc("/attack", attackHandler)
	apiMux.HandleFunc("/ability", abilityHandler)
	apiMux.HandleFunc("/save", saveGameHandler)
	apiMux.HandleFunc("/load", loadGameHandler)

	mux.Handle("/api/", http.StripPrefix("/api", apiMux))

	fileServer := http.FileServer(http.Dir("."))
	mux.Handle("/", fileServer)

	return mux
}
