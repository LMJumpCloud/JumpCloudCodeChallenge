package tests

import (
	"fmt"
	"github.com/MondayHopscotch/JumpCloudCodeChallenge/internal/app"
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
	})
}

func postPassword(pw string, port int) (*http.Response, error) {
	return http.Post(fmt.Sprintf("http://localhost:%v/hash", port), "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf(`password="%s"`, pw)))
}