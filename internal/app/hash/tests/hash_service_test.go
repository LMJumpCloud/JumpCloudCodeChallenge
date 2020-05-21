package tests

import (
	"encoding/json"
	"fmt"
	"github.com/MondayHopscotch/JumpCloudCodeChallenge/internal/app/hash"
	"github.com/MondayHopscotch/JumpCloudCodeChallenge/internal/pkg/hashing"
	"github.com/MondayHopscotch/JumpCloudCodeChallenge/internal/pkg/routing"
	"github.com/MondayHopscotch/JumpCloudCodeChallenge/internal/pkg/test"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
)

const input = `password`

// raw hash: "b109f3bbbc244eb82441917ed06d618b9008dd09b3befd1b5e07394c706a8bb980b1d7785e5976ec049b46df5f1326af5a2ea6d103fd07c95385ffab0cacbc86"
const knownSHA512HashBase64 = "YjEwOWYzYmJiYzI0NGViODI0NDE5MTdlZDA2ZDYxOGI5MDA4ZGQwOWIzYmVmZDFiNWUwNzM5NGM3MDZhOGJiOTgwYjFkNzc4NWU1OTc2ZWMwNDliNDZkZjVmMTMyNmFmNWEyZWE2ZDEwM2ZkMDdjOTUzODVmZmFiMGNhY2JjODY="


func TestHashService(t *testing.T) {
	t.Run("hash endpoint supports PUT", func(t *testing.T) {
		port := 50123
		service := hash.New(port)
		go service.Start()

		expectedID := 1
		resp, err := postPassword(input, port)
		test.AssertNil(t, err, "HTTP error should be null")
		assertPostResponse(t, resp, expectedID)

		expectedID = 2
		resp, err = postPassword(input, port)
		test.AssertNil(t, err, "HTTP error should be null")
		assertPostResponse(t, resp, expectedID)

		service.Stop()
	})

	t.Run("able to retrieve hash", func(t *testing.T) {
		port := 50124
		service := hash.New(port)
		go service.Start()

		expectedID := 1
		resp, err := postPassword(input, port)
		test.AssertNil(t, err, "HTTP error should be null")
		assertPostResponse(t, resp, expectedID)

		resp, err = http.Get(fmt.Sprintf("http://localhost:%v/hash/%v", port, expectedID))
		test.AssertNil(t, err, "HTTP error should be null")
		test.AssertEqual(t, resp.StatusCode, 404, "immediate query should be not found")
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		test.AssertEqual(t, string(bodyBytes), "no hash for id '1' available", "body indicates error")
		resp.Body.Close()

		// sleeping in tests should generally be avoided
		time.Sleep(6 * time.Second)

		resp, err = http.Get(fmt.Sprintf("http://localhost:%v/hash/%v", port, expectedID))
		test.AssertNil(t, err, "HTTP error should be null")
		assertGetResponse(t, resp, expectedID, knownSHA512HashBase64)

		resp.Body.Close()

		service.Stop()
	})

	t.Run("test shutdown call", func(t *testing.T) {
		port := 50125
		service := hash.New(port)
		go service.Start()

		resp, err := http.Get(fmt.Sprintf("http://localhost:%v/shutdown", port))
		test.AssertNil(t, err, "HTTP error should be null")
		test.AssertEqual(t, resp.StatusCode, 200, "request accepted ok")

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		respObj := hash.SimpleMessage{}
		err = json.Unmarshal(bodyBytes, &respObj)
		test.AssertNil(t, err, "unmarshal should not error")
		expected := hash.SimpleMessage{
			Message: "server shutting down",
		}
		test.AssertEqual(t, respObj, expected, "should receive proper get response")
		test.AssertEqual(t, respObj, expected, "body indicates shutdown started")
		resp.Body.Close()

		resp, err = postPassword(input, port)
		test.AssertNotNil(t, err, "HTTP should be rejected resulting in error")
	})

	t.Run("test stats call", func(t *testing.T) {
		port := 50125
		service := hash.New(port)
		go service.Start()

		resp, err := postPassword(input, port)
		test.AssertNil(t, err, "HTTP error should be null")
		test.AssertEqual(t, resp.StatusCode, 201, "201 indicating password hash created")

		resp, err = http.Get(fmt.Sprintf("http://localhost:%v/stats", port))
		test.AssertNil(t, err, "HTTP error should be null")
		test.AssertEqual(t, resp.StatusCode, 200, "request accepted ok")

		bodyBytes, err := ioutil.ReadAll(resp.Body)

		statsResp := routing.RouterStatsResponse{}
		err = json.Unmarshal(bodyBytes, &statsResp)
		test.AssertNil(t, err, "body should be valid json")

		test.AssertEqual(t, len(statsResp.StatsList), 1, "expected number of endpoint stats")
		test.AssertEqual(t, statsResp.StatsList[0].Name, "/hash POST", "properly report POST call name")
		test.AssertEqual(t, statsResp.StatsList[0].Total, 1, "properly report POST call count")
	})
}

func assertPostResponse(t *testing.T, resp *http.Response, expectedID int) {
	test.AssertEqual(t, resp.StatusCode, 201, "201 indicating password hash created")
	bodyContents, err := ioutil.ReadAll(resp.Body)
	respObj := hashing.SubmitResponse{}
	err = json.Unmarshal(bodyContents, &respObj)
	test.AssertNil(t, err, "unmarshal should not error")
	expected := hashing.SubmitResponse{ID: int64(expectedID)}
	test.AssertEqual(t, respObj, expected,"should receive proper post response")
}

func assertGetResponse(t *testing.T, resp *http.Response, id int, hash string) {
	test.AssertEqual(t, resp.StatusCode, 200, "after wait, hash should be available")
	bodyContents, err := ioutil.ReadAll(resp.Body)
	respObj := hashing.GetResponse{}
	err = json.Unmarshal(bodyContents, &respObj)
	test.AssertNil(t, err, "unmarshal should not error")
	expected := hashing.GetResponse{
		ID:   int64(id),
		Hash: hash,
	}
	test.AssertEqual(t, respObj, expected, "should receive proper get response")
}

func postPassword(pw string, port int) (*http.Response, error) {
	return http.Post(fmt.Sprintf("http://localhost:%v/hash", port), "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf(`password=%s`, pw)))
}