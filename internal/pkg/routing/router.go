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
	paramPaths []*ParameterizedPath

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
		Handler: router,
	}
	router.srv = srv
	return router
}

// Serve starts the router as an http server
func (r *Router) Serve() {
	fmt.Println("Server starting on port:", r.port)
	r.errChan<-r.srv.ListenAndServe()
}

// Shutdown calls shutdown on the HTTP server, blocking until shutdown has finished, returning any error that occurs.
// If the error is http.ErrServerClosed, nil is returned instead as this is the expected error and indicates a clean
// shutdown of the server
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
		if IsParameterizedPath(path) {
			r.paramPaths = append(r.paramPaths, ParseParameterizedPath(path))
		}
		r.registeredPaths[path] = handler
		r.mux.HandleFunc(path, handler)
	}
}

// AvailablePaths returns all registered paths for this server
func (r *Router) AvailablePaths() []string {
	paths := make([]string, len(r.registeredPaths) + len(r.paramPaths))
	i := 0
	for path := range r.registeredPaths {
		paths[i] = path
		i++
	}
	for _, path := range r.paramPaths {
		paths[i] = path.Path
		i++
	}
	sort.Strings(paths)
	return paths
}

// ServeHTTP looks at all incoming requests and handles parsing any parameterized paths before passing the
// request to the correct underlying handler for processing
func (r *Router) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	for _, paramPath := range r.paramPaths {
		if paramPath.ParseRequest(req) {
			// Request has been updated. No more work needs to be done before serving
			break
		}
	}
	r.mux.ServeHTTP(writer, req)
}
