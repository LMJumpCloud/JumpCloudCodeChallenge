package endpoints

import (
	"fmt"
	"github.com/MondayHopscotch/JumpCloudCodeChallenge/internal/pkg/hashing"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

var passwordID int64 = 0
var availableHashes = map[int64]string{}

func HashEndpointHandler(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		handlePost(writer, req)
	case http.MethodGet:
		handleGet(writer, req)
	default:
		writer.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handlePost(writer http.ResponseWriter, req *http.Request) {
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

	id := atomic.AddInt64(&passwordID, 1)
	go waitAndPrepare(id, userPassword)
	writer.WriteHeader(http.StatusCreated)
	writer.Write([]byte(fmt.Sprintf("%d", passwordID)))

}

func handleGet(writer http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("unable to parse form data"))
		return
	}

	id, err := getHashIDFromReq(req)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(err.Error()))
	}

	if hash, ok := availableHashes[id]; ok {
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte(hash))
		return
	} else {
		writer.WriteHeader(http.StatusNotFound)
		writer.Write([]byte(fmt.Sprintf("no hash for id '%v' available", id)))
		return
	}
}

func getHashIDFromReq(req *http.Request) (int64, error) {
	splits := strings.Split(req.URL.Path, "/")
	for i, val := range splits {
		if val == "hash" {
			if len(splits) > i {
				id, err := strconv.Atoi(splits[i+1])
				return int64(id), err
			}
		}
	}
	return -1, fmt.Errorf("no id path parameter provided")
}

func waitAndPrepare(id int64, pw string) {
	time.Sleep(5 * time.Second)
	availableHashes[id] = hashing.GetHash(pw)
}