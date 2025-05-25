package utils

import (
	"sync"
	"time"
)

type AttemptData struct {
	Count       int
	LockedUntil time.Time
}

var attemptStore = make(map[string]*AttemptData)
var mu sync.Mutex

const maxAttempts = 7
const lockDuration = 7 * time.Minute

func IsBlocked(key string) bool {
	mu.Lock()
	defer mu.Unlock()

	if data, ok := attemptStore[key]; ok {
		if data.Count >= maxAttempts && time.Now().Before(data.LockedUntil) {
			return true
		}
	}
	return false
}

func RegisterFail(key string) {
	mu.Lock()
	defer mu.Unlock()

	data, exists := attemptStore[key]
	if !exists {
		attemptStore[key] = &AttemptData{Count: 1}
		return
	}

	data.Count++
	if data.Count >= maxAttempts {
		data.LockedUntil = time.Now().Add(lockDuration)
	}
}

func ResetAttempts(key string) {
	mu.Lock()
	defer mu.Unlock()
	delete(attemptStore, key)
}
