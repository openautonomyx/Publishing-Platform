package main

import (
	"sync"
	"time"
)

// ============================================================================
// Rate Limiter Implementation (Token Bucket)
// ============================================================================

type RateLimiter struct {
	requestsPerSecond int
	burstSize         int
	clients           map[string]*ClientBucket
	mu                sync.RWMutex
	cleanupTicker     *time.Ticker
}

type ClientBucket struct {
	tokens    float64
	lastCheck time.Time
}

// ============================================================================
// Constructor
// ============================================================================

func NewRateLimiter(requestsPerSecond, burstSize int) *RateLimiter {
	rl := &RateLimiter{
		requestsPerSecond: requestsPerSecond,
		burstSize:         burstSize,
		clients:           make(map[string]*ClientBucket),
	}

	// Cleanup stale clients every 5 minutes
	rl.cleanupTicker = time.NewTicker(5 * time.Minute)
	go rl.cleanupClients()

	return rl
}

// ============================================================================
// Rate Limiting Logic
// ============================================================================

func (rl *RateLimiter) IsAllowed(clientIP string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	bucket, exists := rl.clients[clientIP]
	if !exists {
		// New client: start with burst capacity, minus this request's token
		rl.clients[clientIP] = &ClientBucket{
			tokens:    float64(rl.burstSize) - 1,
			lastCheck: time.Now(),
		}
		return true
	}

	// Add tokens based on time elapsed
	now := time.Now()
	elapsed := now.Sub(bucket.lastCheck).Seconds()
	tokensAdded := elapsed * float64(rl.requestsPerSecond)

	bucket.tokens = min(bucket.tokens+tokensAdded, float64(rl.burstSize))
	bucket.lastCheck = now

	// Check if we can allow the request
	if bucket.tokens >= 1.0 {
		bucket.tokens -= 1.0
		return true
	}

	return false
}

// ============================================================================
// Cleanup
// ============================================================================

func (rl *RateLimiter) cleanupClients() {
	for range rl.cleanupTicker.C {
		rl.mu.Lock()

		now := time.Now()
		for ip, bucket := range rl.clients {
			// Remove clients inactive for more than 1 hour
			if now.Sub(bucket.lastCheck) > time.Hour {
				delete(rl.clients, ip)
			}
		}

		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) Stop() {
	if rl.cleanupTicker != nil {
		rl.cleanupTicker.Stop()
	}
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
