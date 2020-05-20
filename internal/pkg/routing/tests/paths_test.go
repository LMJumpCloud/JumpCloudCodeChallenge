package tests

import (
	"github.com/MondayHopscotch/JumpCloudCodeChallenge/internal/pkg/routing"
	"github.com/MondayHopscotch/JumpCloudCodeChallenge/internal/pkg/test"
	"net/http"
	"net/url"
	"testing"
)

func TestPaths(t *testing.T) {
	t.Run("Test path split", func(t *testing.T) {
		splits := routing.SplitPath("/path/with/four/elements")
		test.AssertEqual(t, len(splits), 5, "proper split length")
		test.AssertEqual(t, splits[0], "", "root element is empty string")
		test.AssertEqual(t, splits[1], "path", "first element 'path'")
		test.AssertEqual(t, splits[2], "with", "second element 'with'")
		test.AssertEqual(t, splits[3], "four", "third element 'four'")
		test.AssertEqual(t, splits[4], "elements", "fourth element 'elements'")
	})

	t.Run("not parameterized path", func(t *testing.T) {
		test.AssertEqual(t, routing.IsParameterizedPath("/this/is/not"), false, "non-parameterized path")
	})
	t.Run("parameterized path matching", func(t *testing.T) {
		test.AssertEqual(t, routing.IsParameterizedPath("/this/is/{a}/param/path"), true, "parameterized path")
	})

	t.Run("parse parameterize path", func(t *testing.T) {
		path := routing.ParseParameterizedPath("/this/{parses}/well/{enough}")
		test.AssertEqual(t, path.Length, 5, "path length of 5")
		test.AssertEqual(t, len(path.Subs), 2, "two parameter in path")
		test.AssertEqual(t, path.Subs[2], "parses", "index `2` has `parses` parameter")
		test.AssertEqual(t, path.Subs[4], "enough", "index `4` has `enough` parameter")

		test.AssertEqual(t, path.Route[1], "this", "index `1` has `this` segment")
		test.AssertEqual(t, path.Route[3], "well", "index `3` has `well` segment")
	})

	t.Run("test parsing request", func(t *testing.T) {
		req := &http.Request{
			Method:           "",
			URL:              &url.URL{
				Path:       "/this/is/sparta",
			},
		}
		req.ParseForm()

		paramPath := routing.ParseParameterizedPath("/this/is/{something}")
		paramPath.ParseRequest(req)

		test.AssertEqual(t, req.Form.Get("something"), "sparta", "form properly populated")
	})
}
