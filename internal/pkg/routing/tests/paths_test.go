package tests

import (
	"github.com/MondayHopscotch/JumpCloudCodeChallenge/internal/pkg/routing"
	"github.com/MondayHopscotch/JumpCloudCodeChallenge/internal/pkg/test"
	"testing"
)

func TestPaths(t *testing.T) {
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
}
