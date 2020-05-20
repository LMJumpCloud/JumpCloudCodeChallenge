package tests

import (
	"github.com/MondayHopscotch/JumpCloudCodeChallenge/internal/pkg/routing"
	"github.com/MondayHopscotch/JumpCloudCodeChallenge/internal/pkg/test"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestRouting(t *testing.T) {
	t.Run("Simple route path", func(t *testing.T) {
		r := routing.NewRouter(0)
		r.RegisterPaths(map[string]http.HandlerFunc{
			"/test": func(writer http.ResponseWriter, request *http.Request) {

			},
		})

		test.AssertEqual(t, len(r.AvailablePaths()), 1, "one registered path")
	})

	t.Run("Paths sorted", func(t *testing.T) {
		r := routing.NewRouter(0)
		r.RegisterPaths(map[string]http.HandlerFunc{
			"/test/longer": func(writer http.ResponseWriter, request *http.Request) {},
			"/test": func(writer http.ResponseWriter, request *http.Request) {},
		})

		test.AssertEqual(t, len(r.AvailablePaths()), 2, "two registered path")
		test.AssertEqual(t, r.AvailablePaths()[0], "/test", "short path first")
		test.AssertEqual(t, r.AvailablePaths()[1], "/test/longer", "long path second")
	})

	t.Run("server starts", func(t *testing.T) {
		r := routing.NewRouter(8098)
		r.RegisterPaths(map[string]http.HandlerFunc{
			"/test": func(writer http.ResponseWriter, request *http.Request) {
				writer.Write([]byte("helloWorld"))
				writer.WriteHeader(200)
			},
		})

		r.Serve()
		defer func() {
			err := r.Shutdown()
			test.AssertNil(t, err, "no server close error expected")
		}()

		resp, err := http.Get("http://127.0.0.1:8098/test")
		test.AssertNil(t, err, "no error on http GET")

		test.AssertEqual(t, resp.StatusCode, 200, "successful response")

		body, err := ioutil.ReadAll(resp.Body)
		test.AssertEqual(t, string(body), "helloWorld", "hello world body")
	})
}
