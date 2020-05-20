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

func (r *Router) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	pathSplits := SplitPath(req.URL.Path)
	match := r.checkPathMatch(pathSplits)
	if match != nil {
		rigRequest(req, match)
	}
	r.mux.ServeHTTP(writer, req)
}

func (r *Router) checkPathMatch(in []string) *ParameterizedPath {
	var matchFound bool
	for _, path := range r.paramPaths {
		if path.Length == len(in) {
			matchFound = true
			for i, segment := range path.Route {
				if in[i] != segment {
					matchFound = false
				}
			}
			if matchFound {
				return path
			}
		}
	}

	return nil
}

func rigRequest(req *http.Request, pathParams *ParameterizedPath) {
	pathSplits := SplitPath(req.URL.Path)
	for i, s := range pathParams.Subs {
		req.Form.Add(s, pathSplits[i])
	}
	req.URL.Path = pathParams.Path
}
