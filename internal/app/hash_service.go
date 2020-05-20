package app

import (
	"fmt"
	"github.com/MondayHopscotch/JumpCloudCodeChallenge/internal/app/endpoints"
	"github.com/MondayHopscotch/JumpCloudCodeChallenge/internal/pkg/hashing"
	"github.com/MondayHopscotch/JumpCloudCodeChallenge/internal/pkg/routing"
	"net/http"
)

// HashService ties the router and the hash store together
type HashService struct {
	router *routing.Router
	hashStore *hashing.HashStore
}

// NewHashService returns a new instance of HashService
func NewHashService(port int) *HashService {
	return &HashService{
		router:    routing.NewRouter(port),
		hashStore: hashing.NewHashStore(),
	}
}

// Start will register all endpoints and start the HTTP server
func (h *HashService) Start() {
	hashEndpoint := endpoints.HashEndpointForStore(h.hashStore)
	h.router.RegisterPaths(map[string]http.HandlerFunc{
		"/hash": hashEndpoint.Handler,
		"/hash/{id}": hashEndpoint.Handler,
		"/shutdown": h.shutdownHandler,
	})
	h.router.Serve()
}

func (h *HashService) shutdownHandler(writer http.ResponseWriter, req *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("server shutting down"))

	// call this in a goroutine so this request can return correctly
	go h.Stop()
}

// Stop shuts down the HTTP server
func (h *HashService) Stop() {
	err := h.router.Shutdown()
	if err != nil {
		fmt.Println("ERROR: error while shutting down router: ", err)
	}

	h.hashStore.Flush()
}