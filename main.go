package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"log"
	"net/http"
	"time"
)

func main() {
	errorChain := alice.New(loggerHandler, recoverHandler)

	var r = mux.NewRouter()
	//r.HandleFunc("/", rootHandler)
	r.HandleFunc("/welcome", rootHandler).Name("welcome")
	r.HandleFunc("/status", statusHandler).Name("status")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./data/")))

	http.Handle("/", errorChain.Then(r))

	server := &http.Server{
		Addr: ":8080",
	}

	log.Printf("Service UP\n")

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	//render("./data/index.html", w, r)
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "UP")
}

func loggerHandler(h http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h.ServeHTTP(w, r)
		log.Printf("<< %s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}

func recoverHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v", err)
				http.Error(w, http.StatusText(500), 500)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
