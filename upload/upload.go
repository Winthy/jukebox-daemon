package upload

import (
	"log"
	"net/http"
)

func FileServer() {
	port := "9345"
	directory := "cache/"
	http.Handle("/", http.FileServer(http.Dir(directory)))
	log.Printf("Serving %s on HTTP port: %s\n", directory, port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
