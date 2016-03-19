package proxy

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/ksev/webr/build"
	"github.com/ksev/webr/config"
)

// Start start the proxy that blocks requests during compilation and displays build errors on failure
func Start() {
	revpr := httputil.NewSingleHostReverseProxy(build.App.URL)
	revpr.ErrorLog = log.New(ioutil.Discard, "webr", 0)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := build.App.CheckAndWait(); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		revpr.ServeHTTP(w, r)
	})

	fmt.Printf("Listening on %s\n", config.Current.Bind)
	http.ListenAndServe(config.Current.Bind, nil)
}
