package main

import (
	"flag"
	"net/http"
)

func main() {
	bind := flag.String("bind", "", "Bind for http listen")
	flag.Parse()

	if *bind == "" {
		panic("Bind needs to be specified")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("Hello"))
	})
	http.ListenAndServe(*bind, nil)
}
