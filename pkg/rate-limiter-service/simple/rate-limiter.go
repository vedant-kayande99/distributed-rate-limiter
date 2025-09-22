package simple

import (
	"sync"
	"time"
)

type RateLimiter struct {
	rateLimit int
	windowSize time.Duration
	users map[string][]time.Time
	mutex sync.Mutex
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		rateLimit: limit,
		windowSize: window,
		users: make(map[string][]time.Time),
	}
}

/*
	Binary Search for the window_start in the timeLog slice and return the first index of the log which is >= the window_start.
*/
func searchWindowStart (timeLog []time.Time, windowStart time.Time) int {
	if len(timeLog) == 0 {
		return 0
	}

	start := 0
	end := len(timeLog)

	for start < end {
		mid := start + (end - start)/2

		if windowStart.Equal(timeLog[mid]) || windowStart.Before(timeLog[mid]) {
			end = mid
		} else {
			start = mid + 1
		}
	}

	return start
}

func (rl *RateLimiter) ShouldAllow(userId string) bool{
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	curTime := time.Now()
	windowStart := curTime.Add(-rl.windowSize)

	userTimeLog := rl.users[userId]
	windowStartIdx := searchWindowStart(userTimeLog, windowStart)
	userTimeLog = userTimeLog[windowStartIdx:]

	if len(userTimeLog) >= rl.rateLimit {
		return false
	}

	userTimeLog = append(userTimeLog, curTime)
	rl.users[userId] = userTimeLog
	return true
}
