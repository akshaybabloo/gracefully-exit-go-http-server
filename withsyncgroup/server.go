package withsyncgroup

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

var wg sync.WaitGroup

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := fmt.Fprintln(w, "<h1>Home of SyncGroup</h1><br><a href='/exit'>Exit</a>")
	if err != nil {
		panic(err)
	}
}

func ExitHandler(w http.ResponseWriter, r *http.Request) {
	// subtract 1 from the WaitGroup
	defer wg.Done()
	w.WriteHeader(http.StatusOK)
	_, err := fmt.Fprintln(w, "<h1>Bye from SyncGroup</h1>")
	if err != nil {
		panic(err)
	}
}

func StartServer() {

	// add an integer, as there is only one exit that we want so add 1
	wg.Add(1)

	r := mux.NewRouter()
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

	// wait till the counter becomes 0
	wg.Wait()

	if err := srv.Shutdown(context.TODO()); err != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}
	fmt.Println("Server closed - SyncGroup")
}
