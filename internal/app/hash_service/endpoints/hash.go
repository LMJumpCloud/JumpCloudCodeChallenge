package endpoints

import (
	"encoding/json"
	"fmt"
	"github.com/MondayHopscotch/JumpCloudCodeChallenge/internal/pkg/hashing"
	"net/http"
	"strconv"
)

const passwordField = "password"
const idField = "id"

// HashEndpoint is a wrapper around the hash endpoint and its interaction with the InMemoryHashStore
type HashEndpoint struct {
	store hashing.HashStorer
}

// HashEndpointForStore returns a new instance of HashEndpoint based on the provided InMemoryHashStore
func HashEndpointForStore(store hashing.HashStorer) *HashEndpoint {
	return &HashEndpoint{store: store}
}

// HandlePost is responsible for submitting new passwords to be hashed
func (he *HashEndpoint) HandlePost(writer http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	err := req.ParseForm()
	if err != nil {
		fmt.Println(err)
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("unable to parse form data"))
		return
	}

	userPassword := req.Form.Get(passwordField)
	if userPassword == "" {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(fmt.Sprintf("must provide '%v' field", passwordField)))
		return
	}

	submitResp := he.store.SubmitPassword(userPassword)
	bytes, err := json.Marshal(submitResp)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("failed to marshal response"))
		return
	}

	writer.WriteHeader(http.StatusCreated)
	writer.Write(bytes)

}

// HandleGet is responsible for getting hashes out of the store
func (he *HashEndpoint) HandleGet(writer http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	err := req.ParseForm()
	if err != nil {
		fmt.Println(err)
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("unable to parse form data"))
		return
	}

	idParam := req.Form.Get(idField)
	id, err := strconv.Atoi(idParam)
	if err != nil {
		fmt.Println(err)
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(fmt.Sprintf("provided id '%v' is not a valid integer", idParam)))
		return
	}

	getResp := he.store.GetHash(int64(id))
	if getResp.Hash == "" {
		writer.WriteHeader(http.StatusNotFound)
		writer.Write([]byte(fmt.Sprintf("no hash for id '%v' available", id)))
		return
	}

	bytes, err := json.Marshal(getResp)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("failed to marshal response"))
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write(bytes)
}