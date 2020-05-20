package routing

import (
	"context"
	"fmt"
	"net/http"
	"sort"
)

// Router holds route and server state
type Router struct {
	mux *http.ServeMux
	registeredPaths map[string]http.HandlerFunc

	port int
	srv *http.Server
	errChan chan error
}

// NewRouter returns a new instance of a router with no registered routes
func NewRouter(port int) *Router {
	router := &Router{
		mux: http.NewServeMux(),
		registeredPaths: make(map[string]http.HandlerFunc),
		port: port,
		errChan: make(chan error, 0),
	}
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler: router.mux,
	}
	router.srv = srv
	return router
}

// Serve starts the router as an http server
func (r *Router) Serve() {
	go func() {
		r.errChan<-r.srv.ListenAndServe()
	}()
}

func (r *Router) Shutdown() error {
	r.srv.Shutdown(context.Background())
	shutdownErr := <-r.errChan
	if shutdownErr == http.ErrServerClosed {
		return nil
	}
	return shutdownErr
}

// RegisterPaths registers the provided paths with this router
func (r *Router) RegisterPaths(routes map[string]http.HandlerFunc) {
	for path, handler := range routes {
		r.registeredPaths[path] = handler
		r.mux.HandleFunc(path, handler)
	}
}

func (r *Router) AvailablePaths() []string {
	paths := make([]string, len(r.registeredPaths))
	i := 0
	for path, _ := range r.registeredPaths {
		paths[i] = path
		i++
	}
	sort.Strings(paths)
	return paths
}
