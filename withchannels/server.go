package withchannels

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var stopHTTPServerChan chan bool

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := fmt.Fprintln(w, "<h1>Home of Channels</h1>")
	if err != nil{
		panic(err)
	}
}

func ExitHandler(w http.ResponseWriter, r *http.Request)  {
	w.WriteHeader(http.StatusOK)
	_, err := fmt.Fprintln(w, "<h1>Bye from Channels</h1>")
	if err != nil{
		panic(err)
	}
	// sends a signal to the channel, this could even be false it doesn't matter
	stopHTTPServerChan <- true
}

func StartServer() {
	stopHTTPServerChan = make(chan bool)
	r:= mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/exit", ExitHandler)

	fmt.Println("Server started at http://127.0.0.1:8000")

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {

		// always returns error. ErrServerClosed on graceful close
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// unexpected error. port in use?
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	// wait here till a signal is received
	<- stopHTTPServerChan
	if err := srv.Shutdown(context.TODO()); err != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}
	fmt.Println("Server closed - Channels")
}