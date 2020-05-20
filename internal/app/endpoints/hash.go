package endpoints

import (
	"fmt"
	"github.com/MondayHopscotch/JumpCloudCodeChallenge/internal/pkg/hashing"
	"net/http"
	"strconv"
)

// HashEndpoint is a wrapper around the hash endpoint and its interaction with the HashStore
type HashEndpoint struct {
	store *hashing.HashStore
}

// HashEndpointForStore returns a new instance of HashEndpoint based on the provided HashStore
func HashEndpointForStore(store *hashing.HashStore) *HashEndpoint {
	return &HashEndpoint{store: store}
}

// Handler is the HTTP server tie-in to handle requests
func (he *HashEndpoint) Handler(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		he.handlePost(writer, req)
	case http.MethodGet:
		he.handleGet(writer, req)
	default:
		writer.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (he *HashEndpoint) handlePost(writer http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	userPassword := req.Form.Get("password")
	if userPassword == "" {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("must provide 'password' field"))
		return
	}

	id := he.store.SubmitPassword(userPassword)
	writer.WriteHeader(http.StatusCreated)
	writer.Write([]byte(fmt.Sprintf("%d", id)))

}

func (he *HashEndpoint) handleGet(writer http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("unable to parse form data"))
		return
	}

	idParam := req.Form.Get("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(err.Error()))
		return
	}

	hash := he.store.GetHash(int64(id))
	if hash == "" {
		writer.WriteHeader(http.StatusNotFound)
		writer.Write([]byte(fmt.Sprintf("no hash for id '%v' available", id)))
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(hash))
}