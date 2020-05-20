package app

import (
	"fmt"
	"github.com/MondayHopscotch/JumpCloudCodeChallenge/internal/app/endpoints"
	"github.com/MondayHopscotch/JumpCloudCodeChallenge/internal/pkg/routing"
	"net/http"
)

type HashService struct {
	router *routing.Router
	hashStore *HashStore
}

func New(port int) *HashService {
	return &HashService{
		router: routing.NewRouter(port),
		hashStore: NewHashStore(),
	}
}

func (h *HashService) Start() {
	h.router.RegisterPaths(map[string]http.HandlerFunc{
		"/hash": endpoints.HashEndpointHandler,
	})
	h.router.Serve()
}

func (h *HashService) Stop() {
	err := h.router.Shutdown()
	if err != nil {
		fmt.Println("ERROR: error while shutting down router: ", err)
	}
}