package tests

import (
	"github.com/MondayHopscotch/JumpCloudCodeChallenge/internal/pkg/stats"
	"github.com/MondayHopscotch/JumpCloudCodeChallenge/internal/pkg/test"
	"testing"
	"time"
)

func TestAverages(t *testing.T) {
	t.Run("test average", func(t *testing.T) {
		avgr := stats.NewAverageTracker()
		avgr.AddCycleTime("test", 100 * time.Microsecond)
		avgr.AddCycleTime("test", 300 * time.Microsecond)

		allAverages := avgr.GetAverages()
		test.AssertEqual(t, len(allAverages), 1, "only one item being averaged")
		test.AssertEqual(t, allAverages[0].Name, "test", "correct name")
		test.AssertEqual(t, allAverages[0].Total, 2, "correct call count")
		test.AssertEqual(t, allAverages[0].AvgMicroSec, int64(200), "200 micro average")
	})
}
