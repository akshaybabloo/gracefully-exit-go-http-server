package withcontext

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type httpServerHelper struct {
	cancelFunc context.CancelFunc
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := fmt.Fprintln(w, "<h1>Home of Context</h1><br><a href='/exit'>Exit</a>")
	if err != nil {
		panic(err)
	}
}

func (helper *httpServerHelper) ExitHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := fmt.Fprintln(w, "<h1>Bye from Context</h1>")
	if err != nil {
		panic(err)
	}
	// Execute a cancel function
	helper.cancelFunc()
}

func StartServer() {
	stopHTTPServerCtx, cancel := context.WithCancel(context.Background())
	serverHelper := &httpServerHelper{cancelFunc: cancel}
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/exit", serverHelper.ExitHandler)

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

	// Wait till a cancel is executed
	<-stopHTTPServerCtx.Done()
	if err := srv.Shutdown(stopHTTPServerCtx); err != nil && err != context.Canceled {
		panic(err) // failure/timeout shutting down the server gracefully
	}
	fmt.Println("Server closed - Context")
}
