package stats

import "time"

// Average is data around the average for a single item
type Average struct {
	Name        string `json:"name"`
	Total       int    `json:"total"`
	AvgMicroSec int64  `json:"average"`
}

// AverageTracker helps with keeping track of the averages of any number of items
type AverageTracker struct {
	items map[string]*timeTracker
}

// NewAverageTracker returns a new AverageTracker instance
func NewAverageTracker() *AverageTracker {
	return &AverageTracker{
		items: make(map[string]*timeTracker),
	}
}

type timeTracker struct {
	totalTime time.Duration
	totalCount int
}

// AddCycleTime will add one instance having taken the provided duration for the named item
func (a *AverageTracker) AddCycleTime(name string, time time.Duration) {
	if _, ok := a.items[name]; !ok {
		a.items[name] = &timeTracker{
			totalTime: 0,
			totalCount: 0,
		}
	}

	a.items[name].totalCount++
	a.items[name].totalTime += time
}

// GetAverages returns a list of averages for all items currently tracked
func (a *AverageTracker) GetAverages() []Average {
	avgs := make([]Average, len(a.items))
	i := 0
	for name, tracker := range a.items {
		avgs[i] = Average{
			Name:        name,
			Total: tracker.totalCount,
			AvgMicroSec: int64(tracker.totalTime / time.Microsecond / time.Duration(tracker.totalCount) ),
		}
		i++
	}
	return avgs
}