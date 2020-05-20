package tests

import (
	"encoding/json"
	"fmt"
	"github.com/MondayHopscotch/JumpCloudCodeChallenge/internal/app"
	"github.com/MondayHopscotch/JumpCloudCodeChallenge/internal/pkg/stats"
	"github.com/MondayHopscotch/JumpCloudCodeChallenge/internal/pkg/test"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestHashService(t *testing.T) {
	t.Run("hash endpoint supports PUT", func(t *testing.T) {
		port := 50123
		service := app.NewHashService(port)
		go service.Start()

		resp, err := postPassword("password", port)
		test.AssertNil(t, err, "HTTP error should be null")
		test.AssertEqual(t, resp.StatusCode, 201, "201 indicating password hash created")

		expectedID := 1

		bodyContents, err := ioutil.ReadAll(resp.Body)
		test.AssertEqual(t, string(bodyContents), fmt.Sprintf("%v", expectedID), fmt.Sprintf("Should receive id '%v'", expectedID))

		resp, err = postPassword("password", port)
		test.AssertNil(t, err, "HTTP error should be null")
		test.AssertEqual(t, resp.StatusCode, 201, "201 indicating password hash created")

		expectedID = 2

		bodyContents, err = ioutil.ReadAll(resp.Body)
		test.AssertEqual(t, string(bodyContents), fmt.Sprintf("%v", expectedID), fmt.Sprintf("Should receive id '%v'", expectedID))

		service.Stop()
	})

	t.Run("able to retrieve hash", func(t *testing.T) {
		port := 50124
		service := app.NewHashService(port)
		go service.Start()

		resp, err := postPassword("password", port)
		test.AssertNil(t, err, "HTTP error should be null")
		test.AssertEqual(t, resp.StatusCode, 201, "201 indicating password hash created")

		expectedID := 1

		bodyContents, err := ioutil.ReadAll(resp.Body)
		test.AssertEqual(t, string(bodyContents), fmt.Sprintf("%v", expectedID), fmt.Sprintf("Should receive id '%v'", expectedID))

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
		test.AssertEqual(t, resp.StatusCode, 200, "after wait, hash should be available")
		resp.Body.Close()

		service.Stop()
	})

	t.Run("test shutdown call", func(t *testing.T) {
		port := 50125
		service := app.NewHashService(port)
		go service.Start()

		resp, err := http.Get(fmt.Sprintf("http://localhost:%v/shutdown", port))
		test.AssertNil(t, err, "HTTP error should be null")
		test.AssertEqual(t, resp.StatusCode, 200, "request accepted ok")

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		test.AssertEqual(t, string(bodyBytes), "server shutting down", "body indicates shutdown started")
		resp.Body.Close()

		resp, err = postPassword("password", port)
		test.AssertNotNil(t, err, "HTTP should be rejected resulting in error")
	})

	t.Run("test stats call", func(t *testing.T) {
		port := 50125
		service := app.NewHashService(port)
		go service.Start()

		resp, err := postPassword("password", port)
		test.AssertNil(t, err, "HTTP error should be null")
		test.AssertEqual(t, resp.StatusCode, 201, "201 indicating password hash created")

		resp, err = http.Get(fmt.Sprintf("http://localhost:%v/stats", port))
		test.AssertNil(t, err, "HTTP error should be null")
		test.AssertEqual(t, resp.StatusCode, 200, "request accepted ok")

		bodyBytes, err := ioutil.ReadAll(resp.Body)

		statsList := make([]stats.Average, 1)
		err = json.Unmarshal(bodyBytes, &statsList)
		test.AssertNil(t, err, "body should be valid json")

		test.AssertEqual(t, statsList[0].Name, "/hash POST", "properly report POST call name")
		test.AssertEqual(t, statsList[0].Total, 1, "properly report POST call count")
	})
}

func postPassword(pw string, port int) (*http.Response, error) {
	return http.Post(fmt.Sprintf("http://localhost:%v/hash", port), "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf(`password="%s"`, pw)))
}