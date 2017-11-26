package main

import (
	"log"
	"flag"
	"net/http"
)

func main() {
	port := flag.String("p", "8000", "port")
	dir := flag.String("d", "./statics", "dir")
	flag.Parse()

	http.Handle("/", http.FileServer(http.Dir(*dir)))
	log.Printf("Serving %s on HTTP port %s \n", *dir, *port)
	log.Fatal(http.ListenAndServe(":" + *port, nil))
}
