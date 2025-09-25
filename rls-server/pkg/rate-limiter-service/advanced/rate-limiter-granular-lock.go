package advanced

import (
	"context"
	"sync"
	"time"
)

type UserData struct {
	timeLog []time.Time
	mutex sync.RWMutex
}

type RateLimiter struct {
	rateLimit int
	windowSize time.Duration
	users sync.Map	
}

func (rl *RateLimiter) cleanup() {	
	windowStart := time.Now().Add(-rl.windowSize)

	rl.users.Range(func (key, value interface{}) bool {
		userData := value.(*UserData)

		userData.mutex.Lock()
		windowStartIdx := searchWindowStart(userData.timeLog, windowStart)

		if windowStartIdx == len(userData.timeLog) {
			rl.users.Delete(key)
		} else {
			userData.timeLog = userData.timeLog[windowStartIdx:]
		}
		userData.mutex.Unlock()
		return true
	})
}

func (rl *RateLimiter) startCleanup(ctx context.Context) {
	interval := min(rl.windowSize/2, time.Minute)
	ticker := time.NewTicker(interval)

	go func(){
		defer ticker.Stop()
		for {
			select{
			case <-ticker.C:
				rl.cleanup()
			case <-ctx.Done():
				return
			}
		}
	}()
}

func NewRateLimiter(ctx context.Context, limit int, window time.Duration) *RateLimiter {
	rateLimiter := &RateLimiter{
		rateLimit: limit, 
		windowSize: window,
	}	
	rateLimiter.startCleanup(ctx)
	return rateLimiter
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

func (rl *RateLimiter) ShouldAllow(userId string) bool {
	value,_ := rl.users.LoadOrStore(userId, &UserData{
		timeLog: make([]time.Time, 0),
	})
	userData := value.(*UserData)
	
	userData.mutex.RLock()

	curTime := time.Now()
	windowStart := curTime.Add(-rl.windowSize)

 	windowStartIdx := searchWindowStart(userData.timeLog, windowStart)

	if len(userData.timeLog[windowStartIdx:]) >= rl.rateLimit {
		userData.mutex.RUnlock()
		return false
	}
	userData.mutex.RUnlock()

	userData.mutex.Lock()
	defer userData.mutex.Unlock()

	windowStartIdx = searchWindowStart(userData.timeLog, windowStart)
	userData.timeLog = userData.timeLog[windowStartIdx:]
	if len(userData.timeLog) >= rl.rateLimit {
		return false
	}
	userData.timeLog = append(userData.timeLog, curTime)	
	return true
}

// func main() {
// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	limiter := NewRateLimiter(ctx, 5, 10*time.Second)

// 	fmt.Println("Testing the Sliding Window Log Rate Limiter")

// 	for i := 1; i <= 7; i++ {
// 		allowed := limiter.ShouldAllow("user1")
// 		fmt.Printf("Request %d: Allowed: %v\n", i, allowed)
// 		time.Sleep(500 * time.Millisecond)
// 	}

// 	fmt.Println("\nWaiting for window to slide...")
// 	time.Sleep(10 * time.Second)

// 	fmt.Println("\nTesting 'user1' again after window has passed")
// 	allowed := limiter.ShouldAllow("user1")
// 	fmt.Printf("Request 8: Allowed: %v\n", allowed)

// }