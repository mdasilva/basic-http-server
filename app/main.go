package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var creds = make(map[string]string)

func init() {

	// load basic credentials from flag or env
	raw := flag.String("creds", os.Getenv("AUTH_CREDS"),
		"comma separated list of user:pass credentails")
	pairs := strings.Split(*raw, ",")
	for _, pair := range pairs {
		if c := strings.Split(pair, ":"); len(c) == 2 {
			creds[c[0]] = c[1]
		}
	}
}

func main() {
	rtr := mux.NewRouter()
	fs := http.FileServer(http.Dir("www"))
	rtr.PathPrefix("/").Handler(authWrapper(fs))
	loggedRouter := handlers.LoggingHandler(os.Stdout, rtr)
	fmt.Println("listening...")
	log.Fatal(http.ListenAndServe(":8080", loggedRouter))
}

func authWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if user, pass, ok := r.BasicAuth(); !ok {
				unauthorized(w)
			} else {
				if creds[user] != pass {
					unauthorized(w)
				} else {
					h.ServeHTTP(w, r)
				}
			}
		})
}

func unauthorized(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Basic realm="Login required"`)
	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
}
