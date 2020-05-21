package routing

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/MondayHopscotch/JumpCloudCodeChallenge/internal/pkg/stats"
	"net/http"
	"sort"
	"time"
)

// Router holds route and server state
type Router struct {
	mux *http.ServeMux
	registeredPaths map[string]http.HandlerFunc
	paramPaths []*ParameterizedPath

	stats *stats.AverageTracker

	port int
	srv *http.Server
	errChan chan error
}

type RouterStatsResponse struct {
	StatsList []stats.Average `json:statsList`
}

// NewRouter returns a new instance of a router with no registered routes
func NewRouter(port int) *Router {
	router := &Router{
		mux: http.NewServeMux(),
		registeredPaths: make(map[string]http.HandlerFunc),
		stats: stats.NewAverageTracker(),
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

// RegisterStatsEndpoint registers a self-reporting statistics endpoint to show timing metrics on all endpoints
// registered with this router
func (r *Router) RegisterStatsEndpoint() {
	r.RegisterPaths(map[string]http.HandlerFunc{
		"/stats": r.selfStatsHandler,
	})
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
	timer := time.Now()
	req.ParseForm()
	for _, paramPath := range r.paramPaths {
		if paramPath.ParseRequest(req) {
			// Request has been updated. No more work needs to be done before serving
			break
		}
	}

	// All endpoints will return JSON
	writer.Header().Add("Content-Type", "application/json")
	r.mux.ServeHTTP(writer, req)
	r.stats.AddCycleTime(fmt.Sprintf("%v %v", req.URL.Path, req.Method), time.Since(timer))
}

func (r *Router) selfStatsHandler(writer http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	response := RouterStatsResponse{
		StatsList: r.stats.GetAverages(),
	}
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
		writer.WriteHeader(503)
		writer.Write([]byte("failed to generate averages data"))
		return
	}

	writer.WriteHeader(200)
	writer.Write(jsonBytes)
}
