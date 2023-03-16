package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestJitterUp(t *testing.T) {
	testCases := []struct {
		name      string
		duration  time.Duration
		jitter    float64
		numTrials int
	}{
		{
			name:      "no jitter",
			duration:  100 * time.Millisecond,
			jitter:    0,
			numTrials: 100,
		},
		{
			name:      "positive jitter",
			duration:  100 * time.Millisecond,
			jitter:    0.1,
			numTrials: 100,
		},
		{
			name:      "large jitter",
			duration:  100 * time.Millisecond,
			jitter:    5,
			numTrials: 100,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			minDuration := float64(tc.duration) * (1 - tc.jitter)
			maxDuration := float64(tc.duration) * (1 + tc.jitter)

			for i := 0; i < tc.numTrials; i++ {
				jitteredDuration := JitterUp(tc.duration, tc.jitter)
				assert.True(t, float64(jitteredDuration) >= minDuration && float64(jitteredDuration) <= maxDuration,
					"Jittered duration is out of the expected range")
			}
		})
	}
}
