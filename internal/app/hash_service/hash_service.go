package hash_service

import (
	"fmt"
	"github.com/MondayHopscotch/JumpCloudCodeChallenge/internal/app/hash_service/endpoints"
	"github.com/MondayHopscotch/JumpCloudCodeChallenge/internal/pkg/hashing"
	"github.com/MondayHopscotch/JumpCloudCodeChallenge/internal/pkg/routing"
	"net/http"
)

// HashService ties the router and the hash store together
type HashService struct {
	router *routing.Router
	hashStore *hashing.InMemoryHashStore
	done chan struct{}
}

// New returns a new instance of HashService
func New(port int) *HashService {
	return &HashService{
		router:    routing.NewRouter(port),
		hashStore: hashing.NewInMemoryHashStore(),
		done: make(chan struct{}, 0),
	}
}

// Start will register all endpoints and start the HTTP server
func (h *HashService) Start() {
	hashEndpoint := endpoints.HashEndpointForStore(h.hashStore)
	h.router.RegisterPaths(map[string]http.HandlerFunc{
		"/hash": hashEndpoint.HandlePost,
		"/hash/{id}": hashEndpoint.HandleGet,
		"/shutdown": h.shutdownHandler,
	})
	h.router.RegisterStatsEndpoint()
	h.router.Serve()
	<-h.done
}

func (h *HashService) shutdownHandler(writer http.ResponseWriter, req *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("server shutting down"))

	// call this in a goroutine so this request can return correctly
	go h.Stop()
}

// Stop shuts down the HTTP server gracefully and waits for all pending password hashes to finish
func (h *HashService) Stop() {
	fmt.Println("Hash Service shutting down")
	err := h.router.Shutdown()
	if err != nil {
		fmt.Println("error while shutting down router: ", err)
	}
	fmt.Println("HTTP Server shutdown")

	h.hashStore.Flush()
	fmt.Println("All hash processing finished")

	h.done<-struct{}{}
}