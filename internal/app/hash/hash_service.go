package hash

import (
	"encoding/json"
	"fmt"
	"github.com/MondayHopscotch/JumpCloudCodeChallenge/internal/app/hash/endpoints"
	"github.com/MondayHopscotch/JumpCloudCodeChallenge/internal/pkg/hashing"
	"github.com/MondayHopscotch/JumpCloudCodeChallenge/internal/pkg/routing"
	"net/http"
)

// Service ties the router and the hash store together
type Service struct {
	router *routing.Router
	hashStore *hashing.InMemoryHashStore
	done chan struct{}
}

// SimpleMessage is an object with a message
type SimpleMessage struct {
	Message string `json:message`
}

// NewService returns a new instance of the hashing service
func NewService(port int) *Service {
	return &Service{
		router:    routing.NewRouter(port),
		hashStore: hashing.NewInMemoryHashStore(),
		done: make(chan struct{}, 0),
	}
}

// Start will register all endpoints and start the HTTP server
func (h *Service) Start() {
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

func (h *Service) shutdownHandler(writer http.ResponseWriter, req *http.Request) {
	resp := SimpleMessage{Message: "server shutting down"}
	bytes, err := json.Marshal(resp)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("failed to generate server shutdown message"))
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write(bytes)

	// call this in a goroutine so this request can return correctly
	go h.Stop()
}

// Stop shuts down the HTTP server gracefully and waits for all pending password hashes to finish
func (h *Service) Stop() {
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